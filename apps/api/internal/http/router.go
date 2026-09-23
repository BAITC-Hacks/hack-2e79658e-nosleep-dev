package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const maxRequestBytes = 1 << 20

type Config struct {
	WebOrigin  string
	LLMTimeout time.Duration
}

func DefaultConfig() Config {
	return Config{WebOrigin: "http://localhost:3000", LLMTimeout: 25 * time.Second}
}

func NewRouter(engine Engine, explainer Explainer, deck Deck, config Config) *gin.Engine {
	if config.WebOrigin == "" {
		config.WebOrigin = "http://localhost:3000"
	}
	if config.LLMTimeout <= 0 {
		config.LLMTimeout = 25 * time.Second
	}

	router := gin.New()
	router.HandleMethodNotAllowed = true
	router.Use(gin.Logger(), errorRecovery(), cors.New(cors.Config{
		AllowOrigins: []string{config.WebOrigin},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
		MaxAge:       12 * time.Hour,
	}))

	router.GET("/health", health)
	router.GET("/api/v1/health", health)

	v1 := router.Group("/api/v1")
	v1.GET("/catalog", func(c *gin.Context) {
		c.JSON(http.StatusOK, engine.Catalog())
	})
	v1.POST("/simulate", simulate(engine))
	v1.POST("/explain", explain(engine, explainer, config.LLMTimeout))
	v1.GET("/optimum", optimum(engine))
	v1.POST("/compare", compare(engine, explainer, config.LLMTimeout))
	v1.POST("/events/draw", draw(deck))

	router.NoRoute(func(c *gin.Context) {
		writeError(c, http.StatusNotFound, "NOT_FOUND", "маршрут не найден", nil)
	})
	router.NoMethod(func(c *gin.Context) {
		writeError(c, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "метод не поддерживается для этого маршрута", nil)
	})
	return router
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func simulate(engine Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request decisionsRequest
		if !decodeJSON(c, &request) {
			return
		}
		if request.Decisions == nil {
			writeError(c, http.StatusBadRequest, "INVALID_REQUEST", "поле decisions обязательно и должно быть массивом", nil)
			return
		}
		if !validateDecisionFields(c, *request.Decisions) {
			return
		}
		result, violations := engine.Simulate(*request.Decisions)
		if result == nil {
			writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "не удалось рассчитать сценарий", nil)
			return
		}
		if violations == nil {
			violations = []Violation{}
		}
		c.JSON(http.StatusOK, simulateResponse{
			Submittable: len(*request.Decisions) == engine.Catalog().RequiredDecisions && len(violations) == 0,
			Violations:  violations,
			Result:      result,
		})
	}
}

func explain(engine Engine, explainer Explainer, timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request decisionsRequest
		if !decodeJSON(c, &request) {
			return
		}
		if request.Decisions == nil {
			writeError(c, http.StatusBadRequest, "INVALID_REQUEST", "поле decisions обязательно и должно быть массивом", nil)
			return
		}
		if !validateDecisionFields(c, *request.Decisions) {
			return
		}
		result, violations := engine.Simulate(*request.Decisions)
		if result == nil {
			writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "не удалось рассчитать сценарий", nil)
			return
		}
		if len(violations) > 0 {
			writeViolation(c, violations)
			return
		}
		if len(*request.Decisions) != engine.Catalog().RequiredDecisions {
			writeError(c, http.StatusUnprocessableEntity, "DECISION_COUNT", fmt.Sprintf("нужно ровно %d решений, получено %d", engine.Catalog().RequiredDecisions, len(*request.Decisions)), nil)
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		answer := explainer.Explain(ctx, *request.Decisions, result)
		if answer.Source != "live" && answer.Source != "cached" {
			answer.Source = "cached"
		}
		c.JSON(http.StatusOK, explainResponse{Result: result, Explanation: answer, Source: answer.Source})
	}
}

func optimum(engine Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		reveal := false
		if raw, ok := c.GetQuery("reveal"); ok {
			switch strings.ToLower(raw) {
			case "true":
				reveal = true
			case "false":
			default:
				writeError(c, http.StatusBadRequest, "INVALID_QUERY", "параметр reveal должен быть true или false", nil)
				return
			}
		}
		c.JSON(http.StatusOK, engine.Optimum(reveal))
	}
}

func compare(engine Engine, explainer Explainer, timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request compareRequest
		if !decodeJSON(c, &request) {
			return
		}
		if request.Scenarios == nil || len(*request.Scenarios) < 2 || len(*request.Scenarios) > 3 {
			writeError(c, http.StatusBadRequest, "INVALID_REQUEST", "нужно передать от 2 до 3 сценариев", nil)
			return
		}

		response := compareResponse{Scenarios: make([]compareScenarioResult, 0, len(*request.Scenarios))}
		valid := make([]NamedResult, 0, len(*request.Scenarios))
		for i, scenario := range *request.Scenarios {
			label := strings.TrimSpace(scenario.Label)
			if label == "" || len(label) > 80 {
				writeError(c, http.StatusBadRequest, "INVALID_REQUEST", fmt.Sprintf("scenarios[%d].label должен содержать от 1 до 80 символов", i), map[string]any{"index": i})
				return
			}
			if scenario.Decisions == nil {
				writeError(c, http.StatusBadRequest, "INVALID_REQUEST", fmt.Sprintf("scenarios[%d].decisions обязательно и должно быть массивом", i), map[string]any{"index": i})
				return
			}
			if !validateDecisionFields(c, *scenario.Decisions) {
				return
			}
			result, violations := engine.Simulate(*scenario.Decisions)
			if result == nil {
				writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "не удалось рассчитать сценарий", nil)
				return
			}
			if len(violations) == 0 && len(*scenario.Decisions) != engine.Catalog().RequiredDecisions {
				violations = []Violation{{Code: "DECISION_COUNT", Message: fmt.Sprintf("нужно ровно %d решений, получено %d", engine.Catalog().RequiredDecisions, len(*scenario.Decisions))}}
			}
			if len(violations) > 0 {
				response.Scenarios = append(response.Scenarios, compareScenarioResult{Label: label, Violations: violations})
				continue
			}
			response.Scenarios = append(response.Scenarios, compareScenarioResult{Label: label, Result: result})
			valid = append(valid, NamedResult{Label: label, Decisions: *scenario.Decisions, Result: result})
		}

		if len(valid) == 0 {
			response.Comparison = Comparison{Summary: "Ни один сценарий не прошёл проверку, сравнение недоступно.", Source: "cached"}
		} else {
			ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
			response.Comparison = explainer.Compare(ctx, valid)
			cancel()
			if response.Comparison.Source != "live" && response.Comparison.Source != "cached" {
				response.Comparison.Source = "cached"
			}
		}
		c.JSON(http.StatusOK, response)
	}
}

func draw(deck Deck) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request eventRequest
		if c.Request.ContentLength != 0 {
			if !decodeJSON(c, &request) {
				return
			}
		}
		c.JSON(http.StatusOK, deck.Draw(request.Seed))
	}
}

func decodeJSON(c *gin.Context, target any) bool {
	mediaType, _, err := mime.ParseMediaType(c.GetHeader("Content-Type"))
	if err != nil || mediaType != "application/json" {
		writeError(c, http.StatusUnsupportedMediaType, "UNSUPPORTED_MEDIA_TYPE", "укажите Content-Type: application/json", nil)
		return false
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxRequestBytes)
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		writeDecodeError(c, err)
		return false
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			err = errors.New("в теле запроса должно быть одно JSON-значение")
		}
		writeDecodeError(c, err)
		return false
	}
	return true
}

func writeDecodeError(c *gin.Context, err error) {
	status := http.StatusBadRequest
	code := "INVALID_JSON"
	message := "некорректное JSON-тело запроса"
	var maxBytesError *http.MaxBytesError
	if errors.As(err, &maxBytesError) {
		status = http.StatusRequestEntityTooLarge
		code = "REQUEST_TOO_LARGE"
		message = "тело запроса превышает 1 МБ"
	} else if strings.Contains(err.Error(), "unknown field") {
		message = "в запросе есть неизвестное поле: " + err.Error()
	} else if strings.Contains(err.Error(), "cannot unmarshal") {
		message = "значение имеет неверный тип: " + err.Error()
	}
	writeError(c, status, code, message, nil)
}

func validateDecisionFields(c *gin.Context, decisions []Decision) bool {
	for i, decision := range decisions {
		if strings.TrimSpace(decision.InitiativeID) == "" {
			writeError(c, http.StatusBadRequest, "INVALID_REQUEST", fmt.Sprintf("decisions[%d].initiativeId обязательно", i), map[string]any{"index": i})
			return false
		}
		if len(decision.InitiativeID) > 32 || len(decision.DistrictID) > 64 {
			writeError(c, http.StatusBadRequest, "INVALID_REQUEST", fmt.Sprintf("decisions[%d] содержит слишком длинный идентификатор", i), map[string]any{"index": i})
			return false
		}
	}
	return true
}

func writeViolation(c *gin.Context, violations []Violation) {
	if len(violations) == 0 {
		writeError(c, http.StatusUnprocessableEntity, "INVALID_SCENARIO", "сценарий не прошёл проверку", nil)
		return
	}
	first := violations[0]
	writeError(c, http.StatusUnprocessableEntity, first.Code, first.Message, map[string]any{"violations": violations})
}

func writeError(c *gin.Context, status int, code, message string, details any) {
	errorBody := gin.H{"code": code, "message": message}
	if details != nil {
		errorBody["details"] = details
	}
	c.AbortWithStatusJSON(status, gin.H{"error": errorBody})
}

func errorRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil && !c.Writer.Written() {
				writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "внутренняя ошибка сервера", nil)
			}
		}()
		c.Next()
	}
}

type decisionsRequest struct {
	Decisions *[]Decision `json:"decisions"`
}

type compareRequest struct {
	Scenarios *[]compareScenarioRequest `json:"scenarios"`
}

type compareScenarioRequest struct {
	Label     string      `json:"label"`
	Decisions *[]Decision `json:"decisions"`
}

type eventRequest struct {
	Seed *int64 `json:"seed"`
}

type simulateResponse struct {
	Submittable bool        `json:"submittable"`
	Violations  []Violation `json:"violations"`
	*Result
}

type explainResponse struct {
	Result      *Result     `json:"result"`
	Explanation Explanation `json:"explanation"`
	Source      string      `json:"source"`
}

type compareResponse struct {
	Scenarios  []compareScenarioResult `json:"scenarios"`
	Comparison Comparison              `json:"comparison"`
}

type compareScenarioResult struct {
	Label      string      `json:"label"`
	Result     *Result     `json:"result,omitempty"`
	Violations []Violation `json:"violations,omitempty"`
}

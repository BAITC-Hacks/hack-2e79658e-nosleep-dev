package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	httpapi "akim5/api/internal/http"
	"github.com/gin-gonic/gin"
)

func main() {
	port, err := apiPort()
	if err != nil {
		log.Fatal(err)
	}
	config := httpapi.DefaultConfig()
	if origin := os.Getenv("WEB_ORIGIN"); origin != "" {
		config.WebOrigin = origin
	}
	if rawTimeout := os.Getenv("LLM_TIMEOUT_SECONDS"); rawTimeout != "" {
		seconds, err := strconv.Atoi(rawTimeout)
		if err != nil || seconds < 1 || seconds > 300 {
			log.Fatal("LLM_TIMEOUT_SECONDS must be an integer from 1 to 300")
		}
		config.LLMTimeout = time.Duration(seconds) * time.Second
	}

	// Replace this fixture service constructor when the AI/data packages land.
	services, err := httpapi.NewFixtureServices(os.Getenv("CONTRACTS_DIR"))
	if err != nil {
		log.Fatal(err)
	}

	if mode := os.Getenv("GIN_MODE"); mode != "" {
		gin.SetMode(mode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           httpapi.NewRouter(services, services, services, config),
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		log.Printf("API listening on :%s", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	}()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}

func apiPort() (string, error) {
	port := os.Getenv("API_PORT")
	if port == "" {
		return "8000", nil
	}
	value, err := strconv.Atoi(port)
	if err != nil || value < 1 || value > 65535 {
		return "", errors.New("API_PORT must be an integer from 1 to 65535")
	}
	return strconv.Itoa(value), nil
}

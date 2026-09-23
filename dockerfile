FROM golang:1.26-alpine AS build

WORKDIR /src
COPY apps/api/go.mod apps/api/go.sum ./apps/api/
WORKDIR /src/apps/api
RUN go mod download
COPY apps/api/ ./
COPY contracts/ /src/contracts/
COPY data/ /src/data/
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/api ./cmd/server

FROM alpine:3.22

WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=build /out/api ./api
COPY --from=build /src/contracts/ ./contracts/
COPY --from=build /src/data/ ./data/
EXPOSE 8000
ENV API_PORT=8000 CONTRACTS_DIR=/app/contracts
ENTRYPOINT ["/app/api"]

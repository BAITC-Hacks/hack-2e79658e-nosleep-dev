.PHONY: dev dev-api compose-up compose-down

dev:
	@test -f .env || (echo "Copy .env.example to .env first."; exit 1)
	@test -f apps/web/package.json || (echo "apps/web/package.json is not present yet."; exit 1)
	@set -eu; \
	  set -a; . ./.env; set +a; \
	  (cd apps/api && go run ./cmd/server) & api_pid=$$!; \
	  (cd apps/web && npm run dev -- --hostname 0.0.0.0) & web_pid=$$!; \
	  trap 'kill $$api_pid $$web_pid 2>/dev/null || true; wait $$api_pid $$web_pid 2>/dev/null || true' EXIT INT TERM; \
	  wait $$api_pid $$web_pid

dev-api:
	cd apps/api && go run ./cmd/server

compose-up:
	docker compose up --build

compose-down:
	docker compose down

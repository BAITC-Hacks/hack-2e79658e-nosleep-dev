DONE: Gin API, health/CORS, error envelope, fixture services for all routes, validation, tests, Makefile, Compose, Dockerfiles.
ON BRANCH / MERGED: be/api-bootstrap; commit 8b5be3f pushed to origin.
TESTED HOW: `go test ./...`, `go vet ./...`, and built binary `/health` smoke check pass.
BROKEN / RISKY: `apps/web/package.json` is not present yet; `make dev` and Compose await the frontend scaffold.
NEEDS (from whom): AI+data to confirm shared DTO ownership; HTTP mirrors the Go interface locally pending real package types.
NEXT: swap the fixture service constructor after AI/data services land; adapt the HTTP DTO boundary if their concrete types differ.

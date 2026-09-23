DONE: Resolved DATA_DIR for make dev; packaged data for Compose; passed full scenario facts to LLM; extended frontend explain timeout.
ON BRANCH / MERGED: be/startup-integration in an isolated worktree; PR pending.
TESTED HOW: Go suite, frontend lint/build, make dev, seven API paths and web 200 in DEMO_MODE=true; Compose API healthy with event data.
BROKEN / RISKY: no known runtime blocker from these checks.
NEEDS (from whom): AI+data review internal/real adapter; frontend review api-client timeout in the PR.
NEXT: push PR for AI+data and frontend review, then merge after approval.

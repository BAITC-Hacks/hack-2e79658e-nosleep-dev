DONE: OpenAI API is now the default LLM endpoint; OPENAI_API_KEY is the only secret needed and gpt-4.1-mini is the default model.
ON BRANCH / MERGED: codex/openai-api-default; PR pending.
TESTED HOW: Go packages compile with go build ./...; no live call without a user key.
BROKEN / RISKY: A missing or invalid OpenAI key returns the honestly labeled cached fallback.
NEEDS (from whom): AI+data review of its internal/llm change; one teammate approval before merge.
NEXT: open PR, insert key locally, restart make dev, confirm /explain reports source=live.

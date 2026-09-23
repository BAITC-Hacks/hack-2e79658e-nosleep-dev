DONE: OpenAI API is now the default LLM endpoint; OPENAI_API_KEY is the only secret needed and gpt-4.1-mini is the default model.
ON BRANCH / MERGED: codex/openai-api-default pushed; PR creation pending.
TESTED HOW: Go packages compile with go build ./...; no live call without a user key.
BROKEN / RISKY: Missing or invalid OpenAI key returns cached; GitHub connector returned 404 and credential-based PR creation was rejected by auto-review.
NEEDS (from whom): User/team to open PR from pushed branch; AI+data review of internal/llm and one teammate approval before merge.
NEXT: Add key locally, restart make dev, confirm /explain reports source=live.

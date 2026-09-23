DONE: scoring/optimizer/events remain green; LLM now gets concise verified facts, rejects numeric hallucinations, retries once, then returns honest cached fallback.
ON BRANCH / MERGED: ai/llm-grounding / PR pending review.
TESTED HOW: cd apps/api && go test ./...; local Ollama qwen3:4b-instruct returned source=live for /explain and /compare.
BROKEN / RISKY: local Ollama model download and root .env are machine-local; DEMO_MODE=true remains the keyless fallback.
NEEDS (from whom): teammate approval to merge; never commit credentials or machine-local .env.
NEXT: merge PR, then rehearse demo with DEMO_MODE=true and optional local live AI.

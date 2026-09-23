DONE:
ON BRANCH / MERGED:
TESTED HOW:
BROKEN / RISKY:
NEEDS (from whom):
NEXT: implement internal/scoring exactly per contracts/scoring.md, loading data/districts.json, data/initiatives.json, data/rules.json via go:embed. Write it as a table-driven test over every case in data/fixtures/golden.json (9 invalid + 3 valid) — that's your acceptance bar for Checkpoint 1. Then internal/llm (NVIDIA NIM client, structured output, DEMO_MODE + 1-retry-then-cached-fallback per contracts/api.md). See .ai/PLAN.md's task table and .ai/prompts/ai-data.md if pasting into a fresh agent.

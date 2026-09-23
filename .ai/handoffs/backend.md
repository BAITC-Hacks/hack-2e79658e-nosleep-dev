DONE: /explain facts now come from scoring; OpenAI selects 1-2 recommendations from scenario-safe options. /compare guards hypothetical language.
ON BRANCH / MERGED: codex/openai-api-default; PR #18 open against main (mergeable).
TESTED HOW: Go tests pass; 12/12 golden HTTP cases pass; /explain 6/6 live median 0.925s (was ~2.8-3.0s), /compare 3/3 live median 1.327s.
BROKEN / RISKY: cheapest fixture has no expected Score; live summaries are deterministic and LLM selects vetted recommendations.
NEEDS (from whom): AI+data review changes to owned internal/llm in PR #18; one teammate approval before merge.
NEXT: Merge PR #18 after review; API 8000 and web 3000 are running, live /explain smoke check passed.

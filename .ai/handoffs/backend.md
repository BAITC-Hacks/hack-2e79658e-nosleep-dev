DONE: /explain summary now uses two plain sentences: city quality-of-life trend and indicators below the critical threshold; OpenAI still selects vetted recommendations.
ON BRANCH / MERGED: codex/openai-api-default; PR #18 open against main (mergeable).
TESTED HOW: Summary copy and assertions reviewed; gofmt and git diff --check pass. Earlier Go tests, 12/12 fixtures, and live LLM benchmark predate this copy edit.
BROKEN / RISKY: A separate make dev restarted API 8000 before this edit; that process needs a restart to load the new summary. Cheapest fixture has no expected Score.
NEEDS (from whom): AI+data review changes to owned internal/llm in PR #18; one teammate approval before merge.
NEXT: Merge PR #18 after review; restart the existing API process when ready to view the new summary.

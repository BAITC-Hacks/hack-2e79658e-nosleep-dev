DONE: scoring/LLM are merged; the demo script and fresh-clone QA are ready; live OpenAI explanations and comparisons work with the local ignored .env.
ON BRANCH / MERGED: ai/demo-qa / PR #15 pending; OpenAI default and verification merged to main via PRs #16 and #17.
TESTED HOW: fresh clone make dev + Compose web/API 200, three browser rehearsals, cached fallback, Go/verify.py/web lint/build green; root .env gave source=live for /explain (Score 56.54307) and /compare.
BROKEN / RISKY: OpenAI key is local and must never be committed; frontend /compare timeout is 8s vs API LLM 25s. Unsaved picks reset on reload.
NEEDS (from whom): teammate review of shared STATUS/demo docs; frontend owner to check live comparison latency.
NEXT: merge PR #15 after approval, tag final main, and rehearse the judge demo with DEMO_MODE=true unless live key/network are verified immediately beforehand.

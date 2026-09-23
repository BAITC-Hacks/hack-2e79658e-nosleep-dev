DONE: scoring/LLM PR #14 merged; docs-only 60-second demo script and fresh-clone QA completed on main.
ON BRANCH / MERGED: ai/demo-qa / PR pending review.
TESTED HOW: fresh clone make dev + Compose web/API 200; three browser rehearsals, over-budget block, saved compare/reload, unreachable-LLM cached fallback; Go/verify.py/web lint/build green.
BROKEN / RISKY: hosted OpenAI has not yet been verified; frontend /compare timeout is 8s vs API LLM 25s. Unsaved picks reset on reload.
NEEDS (from whom): teammate review of shared STATUS/demo docs; OpenAI owner to confirm source=live; frontend owner to check live comparison latency.
NEXT: merge docs PR, tag final main when all code merges, perform one live-provider rehearsal.

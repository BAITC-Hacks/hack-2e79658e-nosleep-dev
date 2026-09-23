DONE: OpenAI chat-completions now returns grounded live explanations and comparisons; active local .env uses OpenAI, with the former config retained in an ignored backup.
ON BRANCH / MERGED: codex/openai-verification / OpenAI default merged to main in PR #16.
TESTED HOW: root .env + POST /explain returned source=live and Score 56.54307; /compare returned source=live with the contract examples; cd apps/api && go test ./... passed.
BROKEN / RISKY: OpenAI key is local and must never be committed; DEMO_MODE=true remains the keyless fallback. Frontend comparison timeout (8s) is shorter than the API's LLM timeout (25s).
NEEDS (from whom): teammate approval to merge this verification PR; frontend owner may want to align the compare timeout.
NEXT: merge after approval; keep the judge demo on DEMO_MODE=true unless key/network are verified immediately beforehand.

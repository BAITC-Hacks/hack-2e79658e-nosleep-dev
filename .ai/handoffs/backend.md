DONE: OpenAI key verified; schema now requires nonempty explanation lists and permits grounded vague quantities.
ON BRANCH / MERGED: codex/openai-api-default; PR creation pending.
TESTED HOW: OpenAI auth/chat HTTP 200; /explain and /compare HTTP 200 source=live; go test ./... passed.
BROKEN / RISKY: Existing port 8000 process may need restart to load the key and new code; missing/invalid key still returns cached.
NEEDS (from whom): User/team to open PR (GitHub connector 404; credential workaround rejected by auto-review); AI+data review and peer approval.
NEXT: Push fix, restart make dev, merge reviewed PR.

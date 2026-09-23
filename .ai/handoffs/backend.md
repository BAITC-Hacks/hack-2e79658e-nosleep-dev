DONE: OpenAI integration live on port 8000; all 12 golden JSON cases match expected scoring or validation results.
ON BRANCH / MERGED: codex/openai-api-default pushed; PR creation pending.
TESTED HOW: /explain 6/6 live (2 full golden cases x3), /compare 3/3 live; median 2.8-3.0s and 1.25s; verify.py matches base/example scores.
BROKEN / RISKY: LLM sometimes describes simulated actions as implemented and infers "high activity" from budget use; cheapest case has no expected Score in golden.json.
NEEDS (from whom): AI+data prompt/grounding review; user/team to open PR (GitHub connector 404, credential workaround rejected); peer approval.
NEXT: tighten hypothetical wording if requested; merge reviewed PR.

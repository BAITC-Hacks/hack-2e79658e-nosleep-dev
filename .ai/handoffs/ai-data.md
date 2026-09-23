DONE: real adapter now wires scoring, LLM, optimizer and event deck into Gin; scoring synergy shape matches the API contract.
ON BRANCH / MERGED: ai/real-services / pending push and review.
TESTED HOW: cd apps/api && go test ./... (PASS); manual DEMO_MODE endpoint checks: /simulate=56.54307, /explain source=cached, /optimum, seeded /events/draw.
BROKEN / RISKY: root data and embedded Go assets must be kept in sync when authoring data.
NEEDS (from whom): Backend review only; adapter intentionally preserves HTTP DTO ownership.
NEXT: merge PR, then validate frontend against live API.

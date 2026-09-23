DONE: Next 16 + TS + Tailwind + shadcn scaffold; editorial design; health/catalog/simulate/explain client; 5-pick flow; live SVG map + score.
ON BRANCH / MERGED: fe/live-simulator (not merged)
TESTED HOW: lint/build green; Playwright QA at 1280×1000 + 390×844; 5-pick submit/result flow passes.
BROKEN / RISKY: OSM now has 6 Astana districts, mismatching the 5-district contract; using agreed stylized SVG fallback.
NEEDS (from whom): Backend: keep /api/v1 response shapes aligned with contracts/api.md.
NEXT: commit and push; open PR; then test NEXT_PUBLIC_USE_MOCKS=false against backend.

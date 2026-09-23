DONE: Scenario naming/save/delete in localStorage; 2–3 selection; responsive /compare results with source badge, loading/error/invalid states.
ON BRANCH / MERGED: fe/scenario-compare (not merged)
TESTED HOW: lint/build green; Playwright saves 2 scenarios, reloads, compares 2 cards; desktop + 390px screenshots checked.
BROKEN / RISKY: OSM now has 6 Astana districts, mismatching the 5-district contract; using agreed stylized SVG fallback.
NEEDS (from whom): Backend: cmd/server still uses NewFixtureServices; wire merged scoring/LLM/optimizer/events so partial picks change map live.
NEXT: commit/push/PR; after merge, frontend polish only (keyboard/accessibility and demo rehearsal fixes).

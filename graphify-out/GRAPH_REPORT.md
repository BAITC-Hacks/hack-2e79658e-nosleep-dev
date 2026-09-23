# Graph Report - .  (2026-09-23)

## Corpus Check
- Corpus is ~30,651 words - fits in a single context window. You may not need a graph.

## Summary
- 394 nodes · 626 edges · 27 communities (21 shown, 6 thin omitted)
- Extraction: 98% EXTRACTED · 2% INFERRED · 0% AMBIGUOUS · INFERRED: 10 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_simulator.tsx|simulator.tsx]]
- [[_COMMUNITY_router.go|router.go]]
- [[_COMMUNITY_FixtureServices|FixtureServices]]
- [[_COMMUNITY_dependencies|dependencies]]
- [[_COMMUNITY_landing-page.tsx|landing-page.tsx]]
- [[_COMMUNITY_scoring.go|scoring.go]]
- [[_COMMUNITY_types.go|types.go]]
- [[_COMMUNITY_components.json|components.json]]
- [[_COMMUNITY_llm.go|llm.go]]
- [[_COMMUNITY_compilerOptions|compilerOptions]]
- [[_COMMUNITY_Engine|Engine]]
- [[_COMMUNITY_testRouter()|testRouter()]]
- [[_COMMUNITY_scoring_test.go|scoring_test.go]]
- [[_COMMUNITY_Deck|Deck]]
- [[_COMMUNITY_BesSheshim|BesSheshim]]
- [[_COMMUNITY_verify.py|verify.py]]
- [[_COMMUNITY_layout.tsx|layout.tsx]]
- [[_COMMUNITY_TestSeededDrawIsDeterministic()|TestSeededDrawIsDeterministic()]]
- [[_COMMUNITY_TestOptimumIsAtLeastGoldenScenario()|TestOptimumIsAtLeastGoldenScenario()]]
- [[_COMMUNITY_eslint.config.mjs|eslint.config.mjs]]
- [[_COMMUNITY_next.config.ts|next.config.ts]]
- [[_COMMUNITY_postcss.config.mjs|postcss.config.mjs]]

## God Nodes (most connected - your core abstractions)
1. `FixtureServices` - 25 edges
2. `compilerOptions` - 16 edges
3. `NewRouter()` - 14 edges
4. `NewFixtureServices()` - 13 edges
5. `writeError()` - 12 edges
6. `Engine` - 12 edges
7. `explain()` - 11 edges
8. `testRouter()` - 11 edges
9. `compare()` - 10 edges
10. `readFixture()` - 8 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `NewFixtureServices()`  [INFERRED]
  apps/api/cmd/server/main.go → apps/api/internal/http/fixtures.go
- `testRouter()` --calls--> `NewFixtureServices()`  [INFERRED]
  apps/api/internal/http/router_test.go → apps/api/internal/http/fixtures.go
- `NewFromEnv()` --calls--> `Duration`  [INFERRED]
  apps/api/internal/llm/llm.go → apps/api/internal/http/router.go
- `testRouter()` --calls--> `DefaultConfig()`  [INFERRED]
  apps/api/internal/http/router_test.go → apps/api/internal/http/router.go
- `testRouter()` --calls--> `NewRouter()`  [INFERRED]
  apps/api/internal/http/router_test.go → apps/api/internal/http/router.go

## Import Cycles
- None detected.

## Communities (27 total, 6 thin omitted)

### Community 0 - "simulator.tsx"
Cohesion: 0.07
Nodes (28): CityMap(), shapes, directionLabels, Simulator(), ApiError, catalog, configuredApiUrl, mockSimulation() (+20 more)

### Community 1 - "router.go"
Cohesion: 0.12
Nodes (38): Comparison, Context, Decision, Explanation, Result, Violation, compareScenarioRequest, compareScenarioResult (+30 more)

### Community 2 - "FixtureServices"
Cohesion: 0.12
Nodes (22): Catalog, Comparison, Context, Decision, Explanation, Initiative, OptimumResult, Result (+14 more)

### Community 3 - "dependencies"
Cohesion: 0.06
Nodes (32): dependencies, @base-ui/react, class-variance-authority, clsx, cn, @fontsource/ibm-plex-mono, @fontsource-variable/inter, gsap (+24 more)

### Community 4 - "landing-page.tsx"
Cohesion: 0.13
Nodes (19): metadata, clamp01(), DottedGrid(), getCircleRingStrength(), getPlusStrength(), getRawShapeStrength(), getSquareStrength(), getStarStrength() (+11 more)

### Community 5 - "scoring.go"
Cohesion: 0.14
Nodes (23): District, Initiative, synergy, incompatibility, Catalog, Contribution, Decision, District (+15 more)

### Community 6 - "types.go"
Cohesion: 0.09
Nodes (24): District, Initiative, Synergy, Contribution, DistrictResult, Catalog, Comparison, Contribution (+16 more)

### Community 7 - "components.json"
Cohesion: 0.09
Nodes (21): aliases, components, hooks, lib, ui, utils, iconLibrary, menuAccent (+13 more)

### Community 8 - "llm.go"
Cohesion: 0.19
Nodes (15): Context, Decision, Result, T, Client, Comparison, Explainer, Explanation (+7 more)

### Community 9 - "compilerOptions"
Cohesion: 0.10
Nodes (19): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+11 more)

### Community 10 - "Engine"
Cohesion: 0.25
Nodes (12): Catalog, Decision, Initiative, OptimumResult, Result, Violation, Engine, cloneResult() (+4 more)

### Community 11 - "testRouter()"
Cohesion: 0.44
Nodes (11): T, Handler, findRepoRoot(), postJSON(), TestExplainBadDecisionCountUsesContractErrorShape(), TestHealthAndCORS(), TestMalformedRequestTypeReturnsContractError(), testRouter() (+3 more)

### Community 14 - "scoring_test.go"
Cohesion: 0.43
Nodes (6): Decision, T, golden, close(), TestEmbeddedDataLoadsWithoutDataDir(), TestGoldenFixtures()

### Community 15 - "Deck"
Cohesion: 0.62
Nodes (5): Deck, Event, dataDir(), MustNew(), New()

### Community 16 - "BesSheshim"
Cohesion: 0.33
Nodes (6): Graphify incremental update, AI agent knowledge graph workflow, AI city budget simulator, BesSheshim, BS/5, Five decisions. One city.

### Community 17 - "verify.py"
Cohesion: 0.67
Nodes (3): clip(), Reference re-implementation of the Score formula (contracts/scoring.md), indepen, simulate()

## Knowledge Gaps
- **125 isolated node(s):** `T`, `simulationFixture`, `Initiative`, `compareScenarioRequest`, `eventRequest` (+120 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **6 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `NewFixtureServices()` connect `FixtureServices` to `router.go`, `testRouter()`?**
  _High betweenness centrality (0.037) - this node is a cross-community bridge._
- **Why does `Duration` connect `router.go` to `llm.go`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **Why does `main()` connect `router.go` to `FixtureServices`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **Are the 2 inferred relationships involving `NewRouter()` (e.g. with `testRouter()` and `main()`) actually correct?**
  _`NewRouter()` has 2 INFERRED edges - model-reasoned connections that need verification._
- **Are the 2 inferred relationships involving `NewFixtureServices()` (e.g. with `testRouter()` and `main()`) actually correct?**
  _`NewFixtureServices()` has 2 INFERRED edges - model-reasoned connections that need verification._
- **What connects `T`, `simulationFixture`, `Initiative` to the rest of the system?**
  _126 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `simulator.tsx` be split into smaller, more focused modules?**
  _Cohesion score 0.0696969696969697 - nodes in this community are weakly interconnected._
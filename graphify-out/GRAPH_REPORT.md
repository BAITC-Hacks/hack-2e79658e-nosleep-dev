# Graph Report - hack-2e79658e-nosleep-dev  (2026-09-23)

## Corpus Check
- 95 files · ~209,072 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 701 nodes · 1028 edges · 53 communities (37 shown, 16 thin omitted)
- Extraction: 97% EXTRACTED · 3% INFERRED · 0% AMBIGUOUS · INFERRED: 33 edges (avg confidence: 0.87)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `4f5db652`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- [[_COMMUNITY_Community 0|Community 0]]
- [[_COMMUNITY_Community 1|Community 1]]
- [[_COMMUNITY_Community 2|Community 2]]
- [[_COMMUNITY_Community 3|Community 3]]
- [[_COMMUNITY_Community 4|Community 4]]
- [[_COMMUNITY_Community 5|Community 5]]
- [[_COMMUNITY_Community 6|Community 6]]
- [[_COMMUNITY_Community 7|Community 7]]
- [[_COMMUNITY_Community 8|Community 8]]
- [[_COMMUNITY_Community 9|Community 9]]
- [[_COMMUNITY_Community 10|Community 10]]
- [[_COMMUNITY_Community 11|Community 11]]
- [[_COMMUNITY_Community 12|Community 12]]
- [[_COMMUNITY_Community 13|Community 13]]
- [[_COMMUNITY_Community 14|Community 14]]
- [[_COMMUNITY_Community 15|Community 15]]
- [[_COMMUNITY_Community 16|Community 16]]
- [[_COMMUNITY_Community 17|Community 17]]
- [[_COMMUNITY_Community 19|Community 19]]
- [[_COMMUNITY_Community 20|Community 20]]
- [[_COMMUNITY_Community 21|Community 21]]
- [[_COMMUNITY_Community 22|Community 22]]
- [[_COMMUNITY_Community 24|Community 24]]
- [[_COMMUNITY_Community 25|Community 25]]
- [[_COMMUNITY_Community 26|Community 26]]
- [[_COMMUNITY_Community 27|Community 27]]
- [[_COMMUNITY_Community 28|Community 28]]
- [[_COMMUNITY_Community 29|Community 29]]
- [[_COMMUNITY_Community 30|Community 30]]
- [[_COMMUNITY_Community 31|Community 31]]
- [[_COMMUNITY_Community 32|Community 32]]
- [[_COMMUNITY_Community 33|Community 33]]
- [[_COMMUNITY_Community 34|Community 34]]
- [[_COMMUNITY_Community 35|Community 35]]
- [[_COMMUNITY_Community 36|Community 36]]
- [[_COMMUNITY_Community 37|Community 37]]
- [[_COMMUNITY_Community 38|Community 38]]
- [[_COMMUNITY_Community 41|Community 41]]
- [[_COMMUNITY_Community 42|Community 42]]
- [[_COMMUNITY_Community 43|Community 43]]
- [[_COMMUNITY_Community 44|Community 44]]
- [[_COMMUNITY_Community 45|Community 45]]
- [[_COMMUNITY_Community 46|Community 46]]
- [[_COMMUNITY_Community 52|Community 52]]

## God Nodes (most connected - your core abstractions)
1. `FixtureServices` - 25 edges
2. `compilerOptions` - 16 edges
3. `NewRouter()` - 15 edges
4. `NewFixtureServices()` - 13 edges
5. `writeError()` - 12 edges
6. `Engine` - 12 edges
7. `explain()` - 11 edges
8. `testRouter()` - 11 edges
9. `NewLiveServices()` - 10 edges
10. `compare()` - 10 edges

## Surprising Connections (you probably didn't know these)
- `nextConfig` --implements--> `Next.js /api/v1/* rewrite to loopback Go API`  [INFERRED]
  apps/web/next.config.ts → README.md
- `Optional LLM environment passthrough: base URL, API key, model, timeout` --shares_data_with--> `NewFromEnv()`  [INFERRED]
  docker-compose.yml → apps/api/internal/llm/llm.go
- `Inputs` --references--> `Six-district population and indicator dataset`  [EXTRACTED]
  contracts/scoring.md → data/districts.json
- `Steps` --references--> `Six-district population and indicator dataset`  [EXTRACTED]
  contracts/scoring.md → data/districts.json
- `Warm Civic Ambience` --conceptually_related_to--> `Civic Transparent Decision Experience`  [INFERRED]
  apps/web/public/brand/hero-amber.png → docs/BRAND.md

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **** — assets_banner_ai_city_budget_simulator, assets_banner_five_decisions, assets_banner_hundred_budget_units, assets_banner_five_districts [EXTRACTED 1.00]

## Communities (53 total, 16 thin omitted)

### Community 0 - "Community 0"
Cohesion: 0.05
Nodes (43): CityMap(), shapes, DirectionFilter, MobileView, Simulator(), ApiError, catalog, configuredApiUrl (+35 more)

### Community 1 - "Community 1"
Cohesion: 0.06
Nodes (47): Compose app service, APP_BIND and APP_PORT host mapping to container port 3000, Compose restart unless stopped policy, Go API build stage, Container startup command, HTTP health check via web and API, Multi-stage API and web Docker image, Web port 3000 and internal API port 8000 (+39 more)

### Community 2 - "Community 2"
Cohesion: 0.05
Nodes (46): AI city budget simulator, Astana, Kazakhstan, BesSheshim, Five decisions, Five decisions. One city., Five districts, 100 budget units, NoSleep.dev / BAITC 2026 (+38 more)

### Community 3 - "Community 3"
Cohesion: 0.07
Nodes (37): Comparison, Context, Decision, Deck, Engine, Event, Explainer, Explanation (+29 more)

### Community 4 - "Community 4"
Cohesion: 0.11
Nodes (40): T, Comparison, Context, Decision, Deck, Engine, Explainer, Explanation (+32 more)

### Community 5 - "Community 5"
Cohesion: 0.12
Nodes (22): Catalog, Comparison, Context, Decision, Event, Explanation, Initiative, NamedResult (+14 more)

### Community 6 - "Community 6"
Cohesion: 0.10
Nodes (24): metadata, ArrowFillButton(), ArrowFillButtonProps, ArrowFillButtonStyle, clamp01(), DottedGrid(), getCircleRingStrength(), getPlusStrength() (+16 more)

### Community 7 - "Community 7"
Cohesion: 0.06
Nodes (34): dependencies, @base-ui/react, class-variance-authority, clsx, cn, @fontsource/ibm-plex-mono, @fontsource-variable/inter, gsap (+26 more)

### Community 8 - "Community 8"
Cohesion: 0.11
Nodes (27): District, Initiative, synergy, Decision, T, incompatibility, Catalog, Contribution (+19 more)

### Community 9 - "Community 9"
Cohesion: 0.13
Nodes (14): Catalog, `GET /catalog`, Inputs, Scoring formula (source of truth: `data/rules.json`), Steps, Validation rules (all must pass, in this order — first failure wins), Worked examples (verified in `scripts/verify.py`, mirrored in `data/fixtures/golden.json`), Almaty district (population share 0.2112) (+6 more)

### Community 10 - "Community 10"
Cohesion: 0.09
Nodes (24): District, Initiative, Synergy, Contribution, DistrictResult, Catalog, Comparison, Contribution (+16 more)

### Community 11 - "Community 11"
Cohesion: 0.10
Nodes (22): Brand Governance, Graphify incremental update, AI agent knowledge graph workflow, Abstract Amber Architectural Hero Background, Vertical Panel Rhythm, Warm Civic Ambience, AI city budget simulator, BES SHESHIM Logo Treatment (+14 more)

### Community 12 - "Community 12"
Cohesion: 0.17
Nodes (16): Context, Decision, Result, T, Optional LLM environment passthrough: base URL, API key, model, timeout, Sample optional LLM configuration for NVIDIA NIM or another OpenAI-compatible service, Client, Comparison (+8 more)

### Community 13 - "Community 13"
Cohesion: 0.09
Nodes (21): aliases, components, hooks, lib, ui, utils, iconLibrary, menuAccent (+13 more)

### Community 14 - "Community 14"
Cohesion: 0.10
Nodes (19): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+11 more)

### Community 15 - "Community 15"
Cohesion: 0.15
Nodes (17): AstanaThreeMap(), BOUNDARIES, Boundary, createMapTiles(), makeRegion(), MapProps, MapRuntime, ORIGIN (+9 more)

### Community 16 - "Community 16"
Cohesion: 0.44
Nodes (11): T, Handler, findRepoRoot(), postJSON(), TestExplainBadDecisionCountUsesContractErrorShape(), TestHealthAndCORS(), TestMalformedRequestTypeReturnsContractError(), testRouter() (+3 more)

### Community 17 - "Community 17"
Cohesion: 0.20
Nodes (9): Cut order if behind (from `.ai/DECISIONS.md`), Git workflow, Hard rules (never traded for speed), Mission, Ownership, Shared state, Timeline and checkpoints, WORKFLOW — 3 peers, no lead, 5 hours (+1 more)

### Community 19 - "Community 19"
Cohesion: 0.22
Nodes (8): Brief's own "Критерии проверки" → what proves it, JUDGING — how we score against the rubric, Official criteria (100 pts), README и воспроизводимость — 25 pts, Потенциал развития и оригинальность подхода — 10 pts, Соответствие задаче и работоспособность — 25 pts, Техническая реализация — 25 pts, Ценность и применимость решения — 15 pts

### Community 20 - "Community 20"
Cohesion: 0.22
Nodes (8): Constraints, Kickoff prompt — Landing page agent, Language, Non-negotiable: match DESIGN.md, ObsidianUI + micro-interactions (required, not optional polish), Reference, not a template, Sections to build, When done

### Community 21 - "Community 21"
Cohesion: 0.22
Nodes (8): Color, Design system — «Аким на 5 часов» / Astana QoL simulator, Motion (must respect `prefers-reduced-motion`), Reusable component patterns already established, Shape & spacing, Typography, Voice & content patterns, What "beautiful, matches design.md" means in practice

### Community 22 - "Community 22"
Cohesion: 0.25
Nodes (7): Build sequence (for whoever implements this next), Coordination needed (to post under `NEEDS:` in `.ai/handoffs/frontend.md`), Decisions, Frontend plan — «Аким на 5 часов», Page map, Verification checklist (for implementation time, not now), Where things stand

### Community 24 - "Community 24"
Cohesion: 0.33
Nodes (5): MVP (must ship for the demo), Non-goals (explicitly cut, not forgotten), PLAN — «Аким на 5 часов» / Astana Quality of Life simulator, Stack & architecture, Task list

### Community 25 - "Community 25"
Cohesion: 0.50
Nodes (3): Brand, Knowledge graph, Repository instructions for AI agents

### Community 26 - "Community 26"
Cohesion: 0.67
Nodes (3): clip(), Reference re-implementation of the Score formula (contracts/scoring.md), indepen, simulate()

### Community 27 - "Community 27"
Cohesion: 0.50
Nodes (3): Deploy on Vercel, Getting Started, Learn More

### Community 28 - "Community 28"
Cohesion: 0.67
Nodes (3): Graph Maintenance Workflow, Repository Knowledge Graph, Source Files as Final Authority

### Community 52 - "Community 52"
Cohesion: 0.25
Nodes (12): Catalog, Decision, Initiative, OptimumResult, Result, Violation, Engine, cloneResult() (+4 more)

## Knowledge Gaps
- **257 isolated node(s):** `T`, `simulationFixture`, `Initiative`, `Catalog`, `Violation` (+252 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **16 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `NewLiveServices()` connect `Community 3` to `Community 8`, `Community 12`, `Community 4`, `Community 52`?**
  _High betweenness centrality (0.098) - this node is a cross-community bridge._
- **Why does `NewFromEnv()` connect `Community 12` to `Community 3`, `Community 4`?**
  _High betweenness centrality (0.092) - this node is a cross-community bridge._
- **Why does `Compose app service` connect `Community 1` to `Community 3`, `Community 12`?**
  _High betweenness centrality (0.091) - this node is a cross-community bridge._
- **Are the 3 inferred relationships involving `NewRouter()` (e.g. with `TestLiveServicesServeCalculatedScenario()` and `testRouter()`) actually correct?**
  _`NewRouter()` has 3 INFERRED edges - model-reasoned connections that need verification._
- **Are the 2 inferred relationships involving `NewFixtureServices()` (e.g. with `testRouter()` and `main()`) actually correct?**
  _`NewFixtureServices()` has 2 INFERRED edges - model-reasoned connections that need verification._
- **What connects `T`, `simulationFixture`, `Initiative` to the rest of the system?**
  _260 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Community 0` be split into smaller, more focused modules?**
  _Cohesion score 0.054987212276214836 - nodes in this community are weakly interconnected._
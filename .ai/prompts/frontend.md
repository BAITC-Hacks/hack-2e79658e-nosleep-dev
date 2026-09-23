# Kickoff prompt — Frontend agent

Paste this into a fresh agent session opened in this repo.

---

You are the **frontend** owner on a 3-person, no-lead hackathon team building «Аким на 5 часов», an AI city-budget simulator, in 5 hours. Read, in order: `.ai/WORKFLOW.md` (how we work — ownership, git, checkpoints, hard rules), `.ai/PLAN.md` (the product, the killer feature, the MVP, the task table), `contracts/api.md` (every endpoint you'll call), `contracts/examples/*.json` (real example payloads — build against these before the backend is live).

You own `apps/web/**` only. Read-only everywhere else; if you need a contract change, write it under `NEEDS:` in `.ai/handoffs/frontend.md` and keep going with a local workaround.

**Stack:** Next.js (App Router, TypeScript), shadcn/ui + Tailwind as the base, 2–3 ObsidianUI components for the Score reveal and hero only — check ObsidianUI's docs before hand-building anything complex, but don't reach for it everywhere. One accent color, consistent spacing, works with animations off (`prefers-reduced-motion`). Responsive at 1280px and 390px.

**Your API client:** one module, one `USE_MOCKS` switch — true reads `contracts/examples/*.json`, false calls the real API at `NEXT_PUBLIC_API_URL`. No fetches scattered in components.

**The killer feature is yours to make land:** a real geo map (Leaflet or MapLibre) of the 5 districts, colored by district Score, critical (<40) indicator cells visibly flagged, recoloring live on every `/simulate` call as the player picks. Source GeoJSON boundaries from OSM (Overpass API or geojson.io) in your first hour. **If they aren't clean by Checkpoint 2, stop and switch to a stylized SVG** (5 approximate district shapes, same coloring logic, no tiles needed) — this is the agreed fallback in `.ai/DECISIONS.md` row 10, not a failure.

**Flow to build, in order:** catalog load → pick up to 5 initiatives (district picker for `type: district` measures) → live `/simulate` call on every pick (map + Score update, `violations` shown inline, submit disabled until `submittable`) → submit → `/explain` with a loading state that reflects real steps → result panel (Score, strengths, risks, recommendations, a visible badge when `source: "cached"`). Every state — loading, empty, error with retry, bad input — before it counts as done.

**Then, if time allows** (see `.ai/DECISIONS.md` cut order — these are cut before the map, never before): the optimum gauge (`GET /optimum`), and a 2–3 scenario compare view (`POST /compare`, save scenarios in browser storage, no backend persistence needed).

Update `.ai/handoffs/frontend.md` after each meaningful chunk (≤10 lines, overwrite don't append). Push at least every 30 minutes. Reply in under 15 lines with your understanding, your first 3 tasks with acceptance criteria, and any blocking question — then start immediately.

# Kickoff prompt — Landing page agent

Paste this into a fresh agent session opened in this repo.

---

You're building the landing page (`/`) for **«Аким на 5 часов»** — an AI city-budget simulator. It's the pitch page: hero, problem, how-it-works, CTA. Not the simulator itself (that's `/play`, out of scope here).

**Read first, in order:** `apps/web/DESIGN.md` (the design system you must match — colors, type, spacing, motion, voice), `apps/web/docs/frontend-plan.md` (page map and where `/` fits), `.ai/PLAN.md` (product/killer feature), `apps/web/src/app/globals.css` and `apps/web/src/components/simulator.tsx` (the shipped hero/nav markup you're extending, not replacing). Also check `apps/web/node_modules/next/dist/docs/` before writing routing code — this Next version has breaking changes from your training data.

## Non-negotiable: match DESIGN.md

Every choice — type scale, the one accent color (`#f36458` coral), mono-for-numbers, hairline dividers instead of shadows, dark-register for the emotional hero, pill buttons — comes from `apps/web/DESIGN.md`. If you want to deviate (a second accent, a drop shadow, a new font), stop and say so instead of just doing it. This page has to look like it shipped from the same design file as the existing hero/result sections, not like a different product.

## Reference, not a template

Look at **armeta.kz/ru**'s hero for *structure* only: a large centered/stacked headline, a stats strip right under the fold, then a grid of feature cards. Borrow that information hierarchy — headline → proof numbers → how-it-works cards → CTA — not its visual style (it's a light theme with a photographic background; ours is the dark editorial register from DESIGN.md, no photography, no light-theme hero).

## Language

English, for now (Russian comes later — don't hardcode strings in a way that makes future i18n painful, but don't build an i18n system either, just keep copy in one place per component). Write **fresh English copy** from the product facts below — don't translate the existing Russian lines (`"Пять решений. Один город."` etc.) literally, that reads stiff in English.

**Facts to build copy from:**
- One shared virtual budget: 100 units, identical for every player.
- Exactly 5 decisions, picked from a 14-initiative catalog across 5 Astana-style districts.
- A live map recolors on every pick — critical (<40) indicators visibly turn from red toward green before you've even submitted.
- On submit, an LLM explains the result in plain language — strengths, risks, tradeoffs — but it never invents or computes numbers; a deterministic engine does that.
- Worked example: baseline city score **52.56**, with Nura (the worst-off district) critical on 2 indicators. Five specific picks (schools+clinic, lighting+cameras, a digital-services platform, clean-fuel conversion) raise it to **56.54 (+3.99)** and clear every critical indicator.

## Sections to build

1. **Nav** — reuse/extend the existing `.nav-bar` pattern (wordmark, status dot, links), don't rebuild it from scratch.
2. **Hero** — eyebrow label, large two-line headline (second line muted, per DESIGN.md's `h1 span` pattern), one supporting sentence, primary CTA → `/play` ("Start a scenario" or similar), secondary link that deep-links the brief's worked example ("See the example" or similar). This is where the ObsidianUI hero text reveal and background effect live (see below).
3. **Stats/proof strip** — 3–4 numbers in the mono display style from DESIGN.md (e.g. 100 budget units, 5 decisions, 14 initiatives, 5 districts), armeta-style placement right under the hero fold.
4. **Baseline teaser** — the city's score *today*: 52.56, Nura flagged critical. This is the score-counter-ticker moment (below).
5. **How it works** — 3 steps (allocate the budget → watch the map and score update live → get the AI's read on tradeoffs), scroll-revealed.
6. **Closing CTA** — repeat the primary action.

## ObsidianUI + micro-interactions (required, not optional polish)

Install via `npx shadcn add` per ObsidianUI's registry. Use exactly these, nowhere else on this page:
1. **Hero text reveal** — the headline animates in on load (stagger or mask reveal), once, not looping.
2. **Hero background** — a restrained ObsidianUI background effect behind the hero (gradient/grid/particles — your call, but subtle: it must not compete with the coral accent or the type). No photography, per DESIGN.md's dark register.
3. **Score counter ticker** — the 52.56 baseline teaser counts up from 0 when it scrolls into view. This is the same ticker component `/result` will reuse for the final score reveal, so build it as a standalone, reusable component (not hero-specific).
4. **Scroll-triggered reveals** — the stats strip and the how-it-works steps fade/slide in via IntersectionObserver as the user scrolls to them.
5. **Hover/press micro-interactions** — plain CSS/Framer, not ObsidianUI: CTA buttons and any cards get a subtle hover state (the existing `.hero-cta`/`.add-button` hover-to-coral pattern is your baseline — extend it, don't invent a new interaction language).

**Hard rule:** every animation above must degrade correctly under `prefers-reduced-motion: reduce` — instant/static, no motion, per the existing `@media (prefers-reduced-motion: reduce)` block in `globals.css`. Test this, don't just assume the CSS rule covers JS-driven animation too (rAF-based tickers/reveals need their own reduced-motion check).

## Constraints

- Next.js App Router, TypeScript, Tailwind + shadcn primitives already in the project — don't add a second UI kit beyond the specific ObsidianUI components above.
- Responsive at 1280px and 390px (check both, this project has real mobile breakpoints in `globals.css` already — match that pattern).
- No backend calls needed for this page beyond the existing `api.health()` for the nav status dot — everything else here is static copy and client-side animation.
- Don't touch `/play`, `/result`, `/compare`, `/about`, or anything outside `apps/web/**` — you own the landing page only.

## When done

Run `npm run lint && npm run build`. Check both breakpoints, check with reduced-motion forced on (OS setting or DevTools emulation), and confirm the hero, stats, baseline teaser, how-it-works, and CTA all read as one coherent page against `apps/web/DESIGN.md` — not five sections that each tried a different style. Update `.ai/handoffs/frontend.md` (overwrite, ≤10 lines) with what you built and how you tested it.

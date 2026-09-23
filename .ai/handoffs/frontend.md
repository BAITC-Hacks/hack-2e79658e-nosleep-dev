DONE: Built the English landing page at `/`; preserved the simulator at `/play`.
DONE: Added ObsidianUI text reveal + dotted hero field, reusable score ticker, and scroll reveals.
ON BRANCH / MERGED: codex/landing-page (not merged)
TESTED HOW: `npm run lint && npm run build`; Chrome at 1280px, 390px, and forced reduced motion.
TESTED HOW: Verified static reduced-motion output and repaired the hero reveal's load-time fallback.
BROKEN / RISKY: `/play` still ignores the example query string until scenario URL parsing lands.
NEEDS: None for the landing page.
NEXT: Add scenario URL parsing when `/play` routing work begins.

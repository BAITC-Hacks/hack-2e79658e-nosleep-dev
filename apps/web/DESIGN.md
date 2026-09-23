# Design system — «Аким на 5 часов» / Astana QoL simulator

Extracted from the shipped styles in `src/app/globals.css` and the existing `simulator.tsx` markup (branch `fe/live-simulator`). This is the system to match — new pages and components should read as the same product, not a fresh skin.

## Typography

- **Sans:** Inter Variable (`@fontsource-variable/inter`) — body text, UI labels, buttons.
- **Mono:** IBM Plex Mono, 400/500 (`@fontsource/ibm-plex-mono`) — eyebrows, metadata, numbers, badges. Always uppercase, always tracked.
- **Display headlines** use the sans font at very large sizes with tight tracking, not the mono: `font-size: clamp(64px, 9vw, 128px); line-height: .88; letter-spacing: -.065em; font-weight: 400;` and the feature settings `"cv01","cv11","cv12","cv13","ss07"` (Inter's stylistic alternates — keep these on any big display text, they're what makes the numerals/letterforms feel editorial rather than default).
- **Eyebrow label** pattern (used above every section heading): `font: 12px/1.5 "IBM Plex Mono", monospace; letter-spacing: .05em; text-transform: uppercase; color: #797979` (or `#505b6c` on light backgrounds — the `.dark` variant).
- Section headings sit between the eyebrow and body copy: `clamp(42px, 5vw, 72px)`, weight 400, `letter-spacing: -.045em`.
- Body/lede copy: 18px, `line-height: 1.5`, `letter-spacing: -.18px`, color `var(--ash)` (#b9b9b9) on dark backgrounds.

## Color

Two registers, both already in use — pick per-section, don't mix within one section:

**Dark register** (nav, hero, live panel, result section):
| Token | Value | Use |
|---|---|---|
| `--canvas` | `#0b0b0b` | page/section background |
| `--canvas-soft` | `#212121` | raised surface on dark |
| `--hairline-soft` | `#353535` | borders/dividers on dark |
| `--ash` | `#b9b9b9` | secondary text on dark |
| body text | `#ffffff` | primary text on dark |

**Light register** (simulator workspace section):
| Token | Value | Use |
|---|---|---|
| background | `#ededed` | section background |
| card | `#ffffff` | panels |
| border | `#dedede` / `#ededed` | dividers |
| muted text | `#797979` | secondary text on light |
| primary text | `#0b0b0b` | headings/body on light |

**Shared accent + signal colors:**
| Token | Value | Use |
|---|---|---|
| `--brand` / `--accent` | `#f36458` (coral) | the one accent — CTAs, selected states, critical/error signal |
| success | `#37cd84` | live pulse dot, positive deltas, "ok" status |
| warning | `#e7a465` | "watch" map legend tier |
| stable | `#74d8a3` | "stable" map legend tier |
| error text on dark | `#ff9e96` on `#241312` background | violation/error banners |
| focus ring | `#0052ef` | — |

Rule: **one accent color** (coral). Everything else is neutral grayscale plus the three semantic signal colors (green/amber/red-ish) used only for status, never decoration.

## Shape & spacing

- Radius scale: `--radius-sm: 4px`, `--radius-md: 6px`, `--radius-lg: 12px`. Pills (`999px`) for buttons and status dots.
- Cards: `border: 1px solid <border-color>; border-radius: 12px;`.
- Section padding follows the container pattern: `padding-inline: max(24px, calc((100vw - 1440px) / 2))` — content maxes out at 1440px and gutters to 24px on narrower viewports.
- Dividers are always 1px hairlines, never shadows, for separating list rows and grid columns.

## Motion (must respect `prefers-reduced-motion`)

Existing rules already collapse all animation/transition durations to `.01ms` under `@media (prefers-reduced-motion: reduce)` — any new motion must route through durations/transitions that this query can zero out, not `requestAnimationFrame` loops that ignore it.

Current motion vocabulary:
- Fill/background transitions on hover: `transition: background .2s ease`.
- Map polygon recolor: `transition: fill .35s ease`.
- A `spin` keyframe for loading icons.
- `scroll-behavior: smooth` on `html`.

No existing page has entrance animation, parallax, or scroll-reveal yet — that vocabulary is new territory for the landing page (see the landing prompt), so keep it restrained and consistent with the terse, confident tone below rather than decorative.

## Voice & content patterns

- Eyebrow + big terse headline + one supporting sentence. Headlines are short, declarative, often two short sentences broken across lines (e.g. "Five decisions.<br/>One city."), with the second clause in muted gray (`h1 span { color: #797979 }`).
- Numbers are always in mono, always large, always the visual anchor of a panel (`.score-lockup strong { font-size: 74px }`, `.result-score strong { font-size: 64px }`).
- Status/metadata is always mono, uppercase, tracked, and small (10–12px).
- Buttons: pill-shaped, coral fill for primary actions, white/dark-inverted for the dark-register primary button, hover state swaps to the accent color, disabled state drops opacity/desaturates rather than hiding.

## Reusable component patterns already established

- `.nav-bar`: sticky, 64px, translucent dark blur (`background: rgba(11,11,11,.92); backdrop-filter: blur(18px)`).
- `.status-dot`: 6px circle, gray idle / green "ok" with a soft glow ring / coral "error".
- `.pulse`: same treatment as `.status-dot.ok`, used to mark "this number is live."
- Panel header row: 52px min-height, hairline bottom border, mono uppercase label left, count/status right.
- Badge (shadcn `Badge`): pill, `#212121` bg, `#b9b9b9` text, `#353535` border — used for the "cached" AI-source badge.

## What "beautiful, matches design.md" means in practice

New UI should look like it shipped from the same design file as the hero and result section already do: large confident type, one accent color used sparingly, mono for every number and label, generous negative space, hairline dividers instead of card shadows, dark-register sections for the emotional beats (hero, score reveal) and light-register for working/functional sections. Anything that introduces a second accent color, drop shadows, rounded card gradients, or a decorative font is off-system.

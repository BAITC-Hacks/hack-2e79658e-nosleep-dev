# Scoring formula (source of truth: `data/rules.json`)

This is the exact algorithm `internal/scoring` (owned by AI+data) must implement. The LLM never computes or invents any of these numbers — it only explains a result this engine already produced. All source data lives in `data/*.json`; nothing here is hardcoded in code.

## Inputs

- `districts`: 6 entries, each with `population` (share, sums to 1.0) and 10 `indicators` (0–100, higher is always better) — see `data/districts.json`.
- `initiatives`: 14 measures, each with `direction`, `type` (`district` | `city`), `cost`, `lag` (quarters), and `effects` (indicator → full-strength delta) — see `data/initiatives.json`.
- `decisions`: exactly 5 `{initiativeId, districtId?}` pairs. `districtId` is required for `type: district` and forbidden for `type: city`.

## Steps

1. **Validate** the decision set (see violation codes in `contracts/api.md`). An invalid set is never scored.
2. **Realized fraction.** Horizon `H = 8` quarters. For an initiative with lag `L`, `realized = (H - L) / H`. This scales every effect of that initiative (not just some).
3. **Accumulate deltas per district.** For a `district`-type initiative, its (realized) effects apply only to its `districtId`. For a `city`-type initiative, its (realized) effects apply to all 6 districts identically.
4. **Synergies** (from `data/rules.json.synergies`): if both members of a pair are selected, add the *fixed* bonus (not scaled by lag) to the given indicator, in the district of the `anchor` initiative (always the district-type member of the pair).
5. **New indicator value.** `I'_dk = clip(I_dk + Σ deltas_dk, 0, 100)` for every district `d` and indicator `k` (districts/indicators nobody touched keep their original value).
6. **District score.** `D_d = Σ_k w_k * I'_dk` using `indicatorWeights` from `data/rules.json` (sums to 1.0).
7. **City score.** `D_avg = Σ_d population_d * D_d`.
8. **Critical count.** `N_crit` = number of `(district, indicator)` pairs with `I'_dk < 40` (strict), counted **after** step 5.
9. **Final Score.** `Score = 0.7 * D_avg + 0.3 * min_d(D_d) - 1.0 * N_crit`.

Compute this twice per request: once with all deltas at zero (the **before** baseline, always the same numbers) and once with the real deltas (**after**), so the API and the LLM both get a `scoreBefore`/`scoreAfter`/`scoreDelta` triple to work with.

## Validation rules (all must pass, in this order — first failure wins)

1. Exactly 5 decisions (`DECISION_COUNT`).
2. Every `initiativeId` exists (`UNKNOWN_INITIATIVE`) and is not repeated (`DUPLICATE_INITIATIVE`).
3. `district`-type initiatives have a valid, known `districtId` (`DISTRICT_REQUIRED`, `UNKNOWN_DISTRICT`); `city`-type initiatives have none (`CITY_INITIATIVE_HAS_DISTRICT`).
4. Total cost ≤ 100 (`BUDGET_EXCEEDED`). Unspent budget does not carry any bonus.
5. At most 2 initiatives per `direction` (`DIRECTION_LIMIT_EXCEEDED`) — this guarantees at least 3 directions are touched.
6. No incompatible pair selected (`INCOMPATIBLE`): `M1×M3` always conflict; `M4×M7` and `M5×M13` conflict only when they target the same district.

## Worked examples (verified in `scripts/verify.py`, mirrored in `data/fixtures/golden.json`)

| Scenario | Decisions | Cost | D_avg | min D (district) | N_crit | Score |
|---|---|---|---|---|---|---|
| Base (no decisions) | — | 0 | 56.8729 | 49.18 (nura) | 2 (S1, S2 in Nura) | **52.5650** |
| Brief's example | M7→nura, M8→nura, M10→nura, M12, M5→saryarka | 95 | 57.9948 | 52.9625 (nura) | 0 | **56.4851** (M10+M12 synergy fires) |
| Cheapest valid | M9→nura, M11→nura, M10→nura, M12, M4→nura | 61 | — | — | — | valid, not scored above |

The six-district model keeps the worked example near its original target: `56.4851 − 52.5650 = 3.92`. District indicator values are scenario inputs, not measured conditions.

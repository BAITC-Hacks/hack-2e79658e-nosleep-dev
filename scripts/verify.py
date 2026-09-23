"""Reference re-implementation of the Score formula (contracts/scoring.md),
independent of the Go engine, used to cross-check data/fixtures/golden.json.
Not the production implementation — see internal/scoring (owned by AI+data)."""

import json
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
raw_districts = json.loads((ROOT / "data/districts.json").read_text())
districts = {row["id"]: {"pop": row["population"], **row["indicators"]} for row in raw_districts}
W = json.loads((ROOT / "data/rules.json").read_text())["indicatorWeights"]
KEYS = list(W.keys())

raw_initiatives = json.loads((ROOT / "data/initiatives.json").read_text())
initiatives = {row["id"]: {"direction": row["direction"], "type": row["type"], "cost": row["cost"], "lag": row["lag"], "eff": row["effects"]} for row in raw_initiatives}
raw_synergies = json.loads((ROOT / "data/rules.json").read_text())["synergies"]
SYN = [(row["pair"][0], row["pair"][1], row["indicator"], row["bonus"], row["anchor"]) for row in raw_synergies]

def clip(v):
    return max(0.0, min(100.0, v))

def simulate(decisions):
    # decisions: list of (initiative_id, district_id_or_None)
    deltas = {d: {k: 0.0 for k in KEYS} for d in districts}
    picked = set()
    dist_of = {}
    for iid, did in decisions:
        picked.add(iid)
        if did:
            dist_of[iid] = did
        init = initiatives[iid]
        realized = (8 - init["lag"]) / 8
        targets = list(districts.keys()) if init["type"] == "city" else [did]
        for k, eff in init["eff"].items():
            for t in targets:
                deltas[t][k] += eff * realized
    for a, b, k, bonus, anchor in SYN:
        if a in picked and b in picked:
            deltas[dist_of[anchor]][k] += bonus

    d_scores = {}
    n_crit = 0
    for d, vals in districts.items():
        s = 0.0
        for k in KEYS:
            after = clip(vals[k] + deltas[d][k])
            s += W[k] * after
            if after < 40:
                n_crit += 1
        d_scores[d] = s
    d_avg = sum(districts[d]["pop"] * d_scores[d] for d in districts)
    min_d = min(d_scores.values())
    min_id = min(d_scores, key=d_scores.get)
    score = 0.7 * d_avg + 0.3 * min_d - 1.0 * n_crit
    return d_scores, d_avg, min_d, min_id, n_crit, score

# Base (no decisions)
d_scores, d_avg, min_d, min_id, n_crit, score = simulate([])
print("BASE  D_avg=%.4f minD=%.4f(%s) Ncrit=%d Score=%.4f" % (d_avg, min_d, min_id, n_crit, score))

# Brief's example set: M7 Nura, M8 Nura, M10 Nura, M12 city, M5 Saryarka
example = [("M7", "nura"), ("M8", "nura"), ("M10", "nura"), ("M12", None), ("M5", "saryarka")]
cost = sum(initiatives[i]["cost"] for i, _ in example)
d_scores, d_avg, min_d, min_id, n_crit, score = simulate(example)
print("EX    cost=%d D_avg=%.4f minD=%.4f(%s) Ncrit=%d Score=%.4f" % (cost, d_avg, min_d, min_id, n_crit, score))

# Cheapest valid set: M9 + M11 + M10 + M12 + M4
cheap = [("M9", "nura"), ("M11", "nura"), ("M10", "nura"), ("M12", None), ("M4", "nura")]
cost = sum(initiatives[i]["cost"] for i, _ in cheap)
print("CHEAP cost=%d" % cost)

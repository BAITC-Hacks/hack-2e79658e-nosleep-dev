"""Reference re-implementation of the Score formula (contracts/scoring.md),
independent of the Go engine, used to cross-check data/fixtures/golden.json.
Not the production implementation — see internal/scoring (owned by AI+data)."""

districts = {
    "yesil":    dict(pop=0.27, T1=45, T2=62, E1=68, E2=72, S1=48, S2=55, B1=78, B2=60, C1=75, C2=70),
    "almaty":   dict(pop=0.24, T1=40, T2=75, E1=50, E2=55, S1=60, S2=65, B1=62, B2=52, C1=50, C2=60),
    "saryarka": dict(pop=0.20, T1=50, T2=70, E1=42, E2=40, S1=62, S2=68, B1=58, B2=55, C1=45, C2=55),
    "baikonur": dict(pop=0.13, T1=52, T2=68, E1=55, E2=50, S1=58, S2=60, B1=52, B2=58, C1=55, C2=58),
    "nura":     dict(pop=0.16, T1=55, T2=40, E1=45, E2=65, S1=38, S2=35, B1=55, B2=50, C1=60, C2=50),
}

W = dict(T1=.10, T2=.10, E1=.09, E2=.11, S1=.11, S2=.11, B1=.09, B2=.09, C1=.10, C2=.10)
KEYS = list(W.keys())

initiatives = {
    "M1": dict(direction="transport", type="district", cost=18, lag=2, eff=dict(T1=6, T2=9)),
    "M2": dict(direction="transport", type="city", cost=22, lag=2, eff=dict(T1=4, B2=3)),
    "M3": dict(direction="transport", type="district", cost=30, lag=4, eff=dict(T1=16, T2=20, E2=4)),
    "M4": dict(direction="ecology", type="district", cost=15, lag=2, eff=dict(E1=12, E2=3, B1=2)),
    "M5": dict(direction="ecology", type="district", cost=25, lag=3, eff=dict(E2=14, C1=4)),
    "M6": dict(direction="ecology", type="city", cost=20, lag=4, eff=dict(E1=5, E2=3)),
    "M7": dict(direction="social", type="district", cost=24, lag=3, eff=dict(S1=16)),
    "M8": dict(direction="social", type="district", cost=20, lag=3, eff=dict(S2=14)),
    "M9": dict(direction="social", type="district", cost=10, lag=1, eff=dict(S1=3, S2=3, B1=3)),
    "M10": dict(direction="safety", type="district", cost=12, lag=1, eff=dict(B1=12, B2=2)),
    "M11": dict(direction="safety", type="district", cost=10, lag=1, eff=dict(B2=12, T1=-2)),
    "M12": dict(direction="services", type="city", cost=14, lag=1, eff=dict(C2=5)),
    "M13": dict(direction="services", type="district", cost=28, lag=4, eff=dict(C1=18, E2=2)),
    "M14": dict(direction="services", type="city", cost=16, lag=1, eff=dict(C1=5, C2=2)),
}

SYN = [("M1", "M2", "T1", 2, "M1"), ("M10", "M12", "B1", 2, "M10"), ("M5", "M6", "E2", 2, "M5")]

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

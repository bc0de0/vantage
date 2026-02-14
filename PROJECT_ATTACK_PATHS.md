# Project Attack Path Expansion

## Example graph state

- Seed evidence node: `ev-1` (reachable service fingerprint).
- Action classes:
  - `AC-1` (impact 1.0, risk 0.2) produces `hypothesis`.
  - `AC-2` (impact 1.1, risk 0.2) produces `technique` (objective).
  - `AC-RISKY` (impact 2.0, risk 0.9) produces `technique`.

## Paths found

1. `AC-1 -> AC-2`
2. `AC-2`
3. `AC-RISKY` (pruned when risk threshold is low)

## Score computation

Path scoring follows:

`score = Σ(ImpactWeight) - Σ(RiskWeight) - depthPenalty + confidenceWeight * Σ(hypothesisConfidence)`

Example (`AC-1 -> AC-2`):

- Impact sum = `2.1`
- Risk sum = `0.4`
- Depth penalty = `0.2` (depth 2, penalty 0.1)
- Confidence term = `0.2` (assuming confidence 0.5 and weight 0.2)
- Final score = `2.1 - 0.4 - 0.2 + 0.2 = 1.7`

## Why some paths were pruned

- **Risk pruning:** any path with cumulative risk greater than configured threshold is discarded.
- **Depth pruning:** any path with depth greater than `MaxDepth` is discarded.
- **Feasibility pruning:** if action preconditions fail in the virtual graph snapshot, expansion stops for that branch.
- **Phase pruning:** actions whose phase does not match current (or next) operation phase are skipped.

# Reasoning-Centric Restructure (Phase 1)

This repository now includes a **foundational reasoning-core scaffold** aligned with the cognition-centric architecture:

- `core/reasoning`: deterministic reasoning cycle, hypotheses, confidence scoring
- `core/knowledge`: in-memory operational graph with typed nodes/edges
- `core/opsec`: OPSEC-aware candidate filtering
- `planning`: weighted tactic ranking for attack-path options
- `memory`: session/campaign/global pattern stores
- `core/state/operation_phase.go` + `core/state_machine.go`: phase-driven operation lifecycle

## Design boundaries enforced

- Reasoning layer operates on structured domain objects only.
- Planner and OPSEC evaluator are injected interfaces.
- No controller/UI coupling in reasoning packages.
- Confidence is evidence-backed and explicit.

## Next implementation increments

1. Add telemetry ingest adapters into `intelligence/*`.
2. Build attack-graph path expansion over `core/knowledge.Graph`.
3. Add simulation package for defender-response replay.
4. Wire `api/controllers` to trigger one reasoning cycle and return schema-first output.

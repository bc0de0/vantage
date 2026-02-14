# Ontology Normalization Summary (Phase 2 + Phase 3)

## Canonical directory selected
- Chosen canonical source: `action-classes-normalized/`.
- Rationale: it is already wired into runtime loading (`core/reasoning/engine.go`) and now fully aligned to canonical IDs.

## Final directory tree

```text
action-classes-normalized/
├── _schema.yaml
├── AC-01-Passive-Observation.yaml
├── AC-02-Active-Surface-Discovery.yaml
├── AC-03-Reachability-Validation.yaml
├── AC-04-Service-Identification.yaml
├── AC-05-Protocol-Metadata.yaml
├── AC-06-Version-Enumeration.yaml
├── AC-07-Auth-Surface-Analysis.yaml
├── AC-08-Credential-Validation.yaml
├── AC-09-Access-Establishment.yaml
├── AC-10-Privilege-Assessment.yaml
├── AC-11-Lateral-Reachability.yaml
├── AC-12-Execution-Capability.yaml
├── AC-13-Data-Exposure.yaml
├── AC-14-Impact-Feasibility.yaml
├── AC-15-External-Execution.yaml
└── README.md
```

## Totals
- Total canonical action classes: **15**
- ID range present: **AC-01 → AC-15** (unique, filename-aligned)

## Removed files/directories
- Removed redundant legacy directory: `action-classes/` (including all nested duplicate and alias files).

## Drift resolved summary
- Corrected 8 normalized files where `id` did not match canonical `AC-XX` filename prefix.
- Eliminated alias-ID namespace drift (`AC-VERSION-ENUM`, `AC-ACCESS`, etc.).
- Eliminated duplicate semantic definitions spread across two directories.
- Preserved schema metadata by moving schema document to `action-classes-normalized/_schema.yaml`.
- Updated technique exclusions to include `AC-01`, ensuring all canonical IDs are explicitly represented in technique resolution references.

## Validation run
- `go build ./...` passed.
- `go test ./...` passed.
- Additional structural check passed:
  - 15 unique IDs in canonical directory.
  - Every filename `AC-XX` prefix matches YAML `id`.
  - Every action class ID referenced in `techniques/t1595/technique.go` exists.

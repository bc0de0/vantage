# Ontology Audit (Phase 1)

## Scope
- `action-classes/` (legacy, nested by domain)
- `action-classes-normalized/` (flat canonical candidate)

## Canonical candidates
Preferred canonical IDs inferred from code and validator contracts:
`AC-01` through `AC-15`.

Primary supporting references:
- `ai/validate.go` hard-codes valid ActionClass IDs `AC-01..AC-15`.
- `techniques/t1595/technique.go` resolves and excludes canonical numeric IDs.
- `core/reasoning/engine.go` loads from `action-classes-normalized`.

## Drift detected

### 1) Duplicate ActionClass IDs across directories
Duplicate numeric IDs were present in both trees for:
- `AC-01`, `AC-02`, `AC-03`, `AC-07`, `AC-09`, `AC-13`, `AC-14`.

### 2) Semantic duplicates under alias IDs
Multiple files described the same semantics with non-canonical IDs:
- `AC-VERSION-ENUM` vs canonical `AC-06`
- `AC-PROTOCOL-META` vs canonical `AC-05`
- `AC-AUTH-SURFACE` vs canonical `AC-07`
- `AC-ACCESS` vs canonical `AC-09`
- `AC-LATERAL-REACH` vs canonical `AC-11`
- `AC-EXEC-VALIDATE` vs canonical `AC-12`
- `AC-DATA-VERIFY` vs canonical `AC-13`
- `AC-EXTERNAL` vs canonical `AC-15`
- plus additional legacy aliases (`AC-PASSIVE-SURFACE`, `AC-SERVICE-ID`, `AC-CRED-VALIDATION`, `AC-PRIV-ASSESS`, `AC-CONTROL-WEAKNESS`).

### 3) Filename/ID mismatch drift
Several files used `AC-XX-...yaml` names but had non-matching `id:` values (e.g. `AC-04-Service-Identification.yaml` containing `id: AC-VERSION-ENUM`).

### 4) Case/style drift
Legacy tree used lowercase/slugs (`AC-07-auth-surface.yaml`) while normalized tree used title case (`AC-07-Auth-Surface-Analysis.yaml`).

## Orphan analysis
- Legacy alias-ID files were not referenced by technique resolution paths that require canonical numeric IDs.
- Canonical `AC-01..AC-15` set is the only set valid for AI validation and reasoning guardrails.

## Schema inconsistencies

### Inferred schema
Action class documents should contain:
- `id`, `name`, `description`, `intent_domains`, `roe_category`, `exposure_cost`, `operator_involvement`, `preconditions`, `notes`.

### Findings
- Action class YAMLs generally matched this key set.
- `action-classes/_schema.yaml` is a schema meta-document (different shape: `version`, `fields`, `doctrine`) and not an action class instance.
- Critical inconsistency was semantic (`id` namespace drift), not field-shape drift.

## Files safe to remove
- Entire legacy `action-classes/` tree after canonical migration.
- Alias-ID duplicates and semantic duplicates in legacy tree.

## Files requiring rename or ID correction
Before normalization, these normalized filenames had incorrect IDs:
- `AC-04-Service-Identification.yaml`
- `AC-05-Protocol-Metadata.yaml`
- `AC-06-Version-Enumeration.yaml`
- `AC-08-Credential-Validation.yaml`
- `AC-10-Privilege-Assessment.yaml`
- `AC-11-Lateral-Reachability.yaml`
- `AC-12-Execution-Capability.yaml`
- `AC-15-External-Execution.yaml`

Normalization corrected these by migrating canonical content so filename and `id` align.

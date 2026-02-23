# VANTAGE Full Security Triage

_Date:_ 2026-02-23  
_Scope:_ Entire repository (`/workspace/vantage`)  
_Method:_ Static code review + boundary analysis + available local checks

---

## 1) Executive Summary

This repository is primarily a **policy-driven execution orchestration CLI** and does not expose a long-running network service, database layer, or built-in authentication subsystem. That reduces some classic web-service attack surface, but several security-relevant weaknesses still exist at trust boundaries:

- **Integrity model weakness**: evidence integrity is implemented as a plain SHA-256 hash without a secret/private key, making signatures forgeable by any actor that can write artifacts.
- **Mutable trust-object risk**: the validated intent contract is stored as a pointer and not defensively copied, enabling post-validation mutation by in-process callers.
- **External AI boundary hardening gaps**: no HTTP status enforcement and no response body size limit create reliability/DoS and ambiguous-failure risk.
- **Configuration parsing fragility**: action class “YAML” files are parsed with line-splitting logic rather than a YAML parser, which can produce policy drift or misinterpretation.

---

## 2) Methodology

### Repository-wide checks performed

- Enumerated security-relevant code paths (ROE enforcement, intent validation, execution boundary, evidence integrity, AI/network calls, filesystem parsing).
- Reviewed trust transitions:
  - CLI input → runtime contract
  - contract → ROE gate
  - executor → optional external AI service
  - filesystem YAML-like data → action-class model
  - evidence generation → integrity verification
- Ran unit/integration tests where available.
- Attempted dependency vulnerability scan tool (`govulncheck`) but it is not installed in this environment.

### Commands run

- `rg --files -g '*.go'`
- `rg -n "exec\.Command|os\.Open|os\.Create|http\.Client|InsecureSkipVerify|token|secret|WriteFile|ReadFile|json\.Unmarshal"`
- `go test ./...`
- `govulncheck ./...` (failed: command unavailable)

---

## 3) Architecture + Trust Boundary Notes

- **No built-in auth/authz user model**: authorization is represented as contract + ROE checks inside process, not principal/session-based auth.
- **Primary authority boundary**: `core/executor/engine.go` + `core/roe/enforcer.go`.
- **External trust boundary**: `ai/callTogether()` performs outbound HTTPS calls using env-provided API key.
- **Filesystem trust boundary**: action class metadata is loaded from local YAML files.
- **No native DB boundary** in this codebase; no SQL driver usage found.

---

## 4) Findings

## F-01: Evidence integrity is forgeable (hash-only, no secret/private key)

- **Severity:** High
- **Category:** Encryption / Integrity / Non-repudiation
- **Where:** `core/evidence/signature.go`
- **Why it matters:**
  - `Sign()` computes `sha256(canonicalPayload)` and stores hex digest.
  - There is **no keyed MAC** and **no asymmetric signing key**.
  - Any attacker with artifact write access can modify fields and recompute integrity, defeating tamper-evidence and non-repudiation claims.
- **Exploit sketch:** Modify artifact `Success`/`Output`, recompute SHA-256, overwrite `Integrity`; `Verify()` still passes.
- **Recommended fix:**
  - Replace with `HMAC-SHA256` (minimum) using secret from protected key management, or preferably Ed25519/ECDSA signatures with private key isolation.
  - Include key ID and algorithm metadata in artifact.

## F-02: Intent contract can be mutated after validation (TOCTOU by shared pointer)

- **Severity:** High
- **Category:** Authorization / Trust boundary control
- **Where:** `core/executor/engine.go` (`New` + stored `contract *intent.Contract`)
- **Why it matters:**
  - Engine validates once in `New()`, then stores original pointer.
  - Caller holding the same pointer can mutate `AllowedTechniques`, `Targets`, or time bounds after validation.
  - This violates the stated immutability security model and creates in-process policy bypass risk.
- **Recommended fix:**
  - Deep-copy contract into engine at construction.
  - Optionally re-validate critical fields at execution-time, or seal contract in immutable type.

## F-03: Outbound AI call lacks HTTP status validation

- **Severity:** Medium
- **Category:** Validation / External service trust boundary
- **Where:** `ai/client.go`
- **Why it matters:**
  - Response JSON is decoded regardless of HTTP status.
  - Upstream 4xx/5xx or intermediary error pages can be treated as parse failures without clear handling semantics.
  - Security telemetry and incident triage become ambiguous; control-plane reliability weakens.
- **Recommended fix:**
  - Enforce `2xx` status code before decode.
  - Capture bounded error body for diagnostics; avoid logging secrets.

## F-04: No response size limit on AI JSON decode (memory pressure / DoS risk)

- **Severity:** Medium
- **Category:** Input validation / Availability
- **Where:** `ai/client.go`
- **Why it matters:**
  - `json.NewDecoder(resp.Body).Decode(&parsed)` reads unbounded stream.
  - Malicious/compromised upstream or proxy can return excessively large body.
- **Recommended fix:**
  - Wrap body in `io.LimitReader` with strict max bytes.
  - Add decoder disallow unknown fields where schema strictness is required.

## F-05: Action-class “YAML” parser is custom line parser (misparse/policy drift risk)

- **Severity:** Medium
- **Category:** Validation / Configuration integrity
- **Where:** `core/reasoning/action_class.go`
- **Why it matters:**
  - Parser uses line `SplitN(':', 2)` and simple inline list parsing.
  - Real YAML constructs (multiline strings, nested structures, comments/colons in values) can be silently ignored or misread.
  - Could lead to incorrect phase/precondition assignment and governance drift.
- **Recommended fix:**
  - Use `gopkg.in/yaml.v3` with strict struct decoding.
  - Validate required fields and reject unknown/ambiguous structures.

## F-06: AI advisory failures are fully swallowed (security visibility blind spot)

- **Severity:** Low
- **Category:** Monitoring / Detection / Governance
- **Where:** `core/executor/engine.go`
- **Why it matters:**
  - `_, _ = ai.Advise(...)` discards all errors.
  - While non-authoritative behavior is intentional, complete suppression removes observability for tampering/outage/contract drift.
- **Recommended fix:**
  - Keep non-blocking behavior but record bounded metrics/log events for failure classes.

## F-07: Plaintext evidence output field may store sensitive data without policy guard

- **Severity:** Low
- **Category:** Secrets management / Data handling
- **Where:** `core/evidence/artifact.go`
- **Why it matters:**
  - `Output` is free-form raw content and not redacted/encrypted.
  - If future techniques ingest credentials/tokens, sensitive material may persist in artifacts.
- **Recommended fix:**
  - Add data-classification policy, secret scrubbing, and optional at-rest encryption strategy for storage destinations.

---

## 5) Domain-specific Coverage Requested

### Authentication

- No first-class user authentication implementation exists in repository (CLI local process model).
- Security implication: identity/actor controls must be provided by surrounding platform (CI/CD, host IAM, runtime operator controls).

### Authorization

- Strong intent + ROE checks exist, but mutable contract pointer introduces bypass risk (F-02).

### Encryption / Integrity

- Integrity exists but is non-keyed hash (F-01), insufficient for adversarial tamper resistance.
- No transport customization observed that disables TLS verification.

### Secrets management

- API key sourced from environment variable (`TOGETHER_API_KEY`), which is common but requires runtime hygiene.
- No explicit secret redaction path in evidence outputs (F-07).

### Validation

- Intent and ROE validations are strict.
- Weaknesses exist in external response handling (F-03/F-04) and config parsing robustness (F-05).

### Trust-boundary shifts (cmd → service → db/filesystem/OS)

- `cmd` → core: constrained by validation + ROE; no direct OS command execution discovered.
- core → external service: AI HTTPS call boundary has missing status/body hardening.
- core → filesystem: action class load path relies on fragile parser.
- db boundary: not present in repository.

---

## 6) Prioritized Remediation Plan

1. **Immediate (High):** Replace evidence signature with keyed/asymmetric signing (F-01).
2. **Immediate (High):** Deep-copy + freeze intent contract inside executor (F-02).
3. **Near-term (Medium):** Add HTTP status enforcement + body limits for AI client (F-03/F-04).
4. **Near-term (Medium):** Migrate to strict YAML parser and schema validation (F-05).
5. **Ongoing (Low):** Add non-blocking advisory telemetry + output redaction policies (F-06/F-07).

---

## 7) Validation Notes

- `go test ./...` currently fails due to a pre-existing stress-test threshold breach (`TestPlanCampaignStress1000Techniques`).
- `govulncheck` could not be run because tool is not installed in this execution environment.


module vantage

go 1.25

// --------------------------------------------------------------------
// VANTAGE — DEPENDENCY DOCTRINE (IMMUTABLE)
//
// 1. MINIMALISM
//    Every dependency must justify its presence in a court-admissible
//    evidence chain.
//
// 2. STABILITY
//    Dependencies must be mature, actively maintained, and widely
//    adopted in tier-1 infrastructure projects (e.g., Kubernetes,
//    Docker, Terraform, Vault).
//
// 3. SECURITY
//    All dependencies must pass `govulncheck` with zero HIGH or CRITICAL
//    vulnerabilities at time of inclusion.
//
// 4. AUDITABILITY
//    No dependency may introduce opaque behavior, background activity,
//    or uncontrolled side effects.
// --------------------------------------------------------------------

// --------------------------------------------------------------------
// CORE DEPENDENCIES (NON-NEGOTIABLE)
// --------------------------------------------------------------------

// Cobra — Structured CLI framework.
// Justification:
// - Deterministic command trees
// - Explicit subcommand execution
// - Industry-standard (Kubernetes, Docker)
require github.com/spf13/cobra v1.8.1

// Pflag — POSIX/GNU-style flags.
// Required by Cobra; no standalone usage.
require github.com/spf13/pflag v1.0.5 // indirect

// Mousetrap — Required by Cobra on Windows.
// Included transitively; no direct usage.
require github.com/inconshreveable/mousetrap v1.1.0 // indirect

// YAML v3 — Declarative intent & ROE parsing.
// Used strictly for configuration, never execution.
require gopkg.in/yaml.v3 v3.0.1

// --------------------------------------------------------------------
// CRYPTOGRAPHY (EVIDENCE INTEGRITY)
// --------------------------------------------------------------------

// x/crypto — Supplemental cryptographic primitives.
// Used for hashing/signing evidence artifacts.
// First-party Go project, security-reviewed.
require golang.org/x/crypto v0.19.0

// UUID — Deterministic artifact identifiers.
//
// Justification:
// - Evidence artifacts require globally unique, collision-resistant IDs
// - Used only for identifiers, NOT security primitives
// - Industry standard (Kubernetes, Docker, Terraform)
//
// Policy note:
// - This is NOT used for cryptography
// - This is NOT used for randomness beyond uniqueness
require github.com/google/uuid v1.6.0


// --------------------------------------------------------------------
// EXPLICITLY DEFERRED DEPENDENCIES
// --------------------------------------------------------------------
//
// The following are intentionally NOT included:
//
// - Viper (configuration sprawl, nondeterminism)
// - Logrus/Zap (hooks increase attack surface)
// - UUID libraries (custom deterministic IDs suffice)
// - fsnotify (hot reload is a security risk)
// - validator/mapstructure (custom validation is explicit)
//
// Inclusion of any deferred dependency REQUIRES:
// 1. Demonstrated operational need
// 2. Security review
// 3. ROE impact assessment
// --------------------------------------------------------------------

// --------------------------------------------------------------------
// OPERATIONAL REQUIREMENTS
// --------------------------------------------------------------------
//
// REQUIRED BEFORE EVERY BUILD:
//   go mod verify
//
// REQUIRED BEFORE EVERY COMMIT:
//   govulncheck ./...
//
// This file is considered part of the security boundary.
// Unauthorized modification is a policy violation.
// --------------------------------------------------------------------

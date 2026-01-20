package evidence

import (
	"errors"
	"time"
)

// -----------------------------------------------------------------------------
// Evidence Artifact — FACTUAL RECORD
//
// An Artifact represents a single, immutable fact produced by execution.
//
// Evidence is NOT:
// - a log
// - a report
// - an interpretation
//
// Evidence IS:
// - a timestamped fact
// - bound to intent, technique, and target
// - cryptographically verifiable
//
// DESIGN PRINCIPLES:
//
// 1. FACT OVER OPINION
//    Artifacts record what happened, not what it means.
//
// 2. IMMUTABLE AFTER SIGNING
//    Any mutation invalidates integrity.
//
// 3. MINIMUM NECESSARY DATA
//    Only store what is required to prove execution.
// -----------------------------------------------------------------------------

// Artifact represents a single unit of evidence.
type Artifact struct {

	// ArtifactID uniquely identifies this evidence unit.
	// Generated once and never reused.
	ArtifactID string

	// CampaignID binds evidence to a declared intent.
	CampaignID string

	// TechniqueID identifies the technique executed.
	TechniqueID string

	// Target identifies the execution target.
	Target string

	// ExecutedAt records when execution completed (UTC).
	ExecutedAt time.Time

	// Success indicates whether the technique completed successfully.
	Success bool

	// Output contains raw, uninterpreted execution output.
	//
	// This may include:
	// - command output
	// - banners
	// - protocol responses
	//
	// It MUST NOT contain:
	// - analysis
	// - summaries
	// - inferred impact
	Output string

	// ExposureScore captures exposure at time of execution.
	ExposureScore uint64

	// Integrity contains the cryptographic signature
	// over all other fields.
	Integrity string
}

// Validate performs structural validation prior to signing.
func (a *Artifact) Validate() error {

	if a == nil {
		return errors.New("nil artifact")
	}

	if a.ArtifactID == "" {
		return errors.New("artifact missing artifact_id")
	}

	if a.CampaignID == "" {
		return errors.New("artifact missing campaign_id")
	}

	if a.TechniqueID == "" {
		return errors.New("artifact missing technique_id")
	}

	if a.Target == "" {
		return errors.New("artifact missing target")
	}

	if a.ExecutedAt.IsZero() {
		return errors.New("artifact missing execution timestamp")
	}

	return nil
}

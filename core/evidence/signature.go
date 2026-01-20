package evidence

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"
)

// -----------------------------------------------------------------------------
// Evidence Integrity — TAMPER PROTECTION
//
// Integrity is implemented via deterministic hashing.
// This provides:
//
// - Tamper detection
// - Chain-of-custody confidence
// - Audit defensibility
//
// We DO NOT:
// - encrypt evidence
// - hide evidence
// - obfuscate evidence
//
// Evidence is meant to be READ — not altered.
// -----------------------------------------------------------------------------

// Sign calculates and applies an integrity signature
// over the artifact's immutable fields.
func (a *Artifact) Sign() error {

	if err := a.Validate(); err != nil {
		return err
	}

	if a.Integrity != "" {
		return errors.New("artifact already signed")
	}

	payload, err := canonicalPayload(a)
	if err != nil {
		return err
	}

	hash := sha256.Sum256(payload)
	a.Integrity = hex.EncodeToString(hash[:])

	return nil
}

// Verify checks whether the artifact's integrity
// matches its contents.
func (a *Artifact) Verify() (bool, error) {

	if a.Integrity == "" {
		return false, errors.New("artifact not signed")
	}

	payload, err := canonicalPayload(a)
	if err != nil {
		return false, err
	}

	hash := sha256.Sum256(payload)
	expected := hex.EncodeToString(hash[:])

	return a.Integrity == expected, nil
}

// canonicalPayload produces a deterministic byte representation
// of the artifact excluding the Integrity field.
func canonicalPayload(a *Artifact) ([]byte, error) {

	type signedView struct {
		ArtifactID    string
		CampaignID    string
		TechniqueID   string
		Target        string
		ExecutedAt    time.Time
		Success       bool
		Output        string
		ExposureScore uint64
	}

	view := signedView{
		ArtifactID:    a.ArtifactID,
		CampaignID:    a.CampaignID,
		TechniqueID:   a.TechniqueID,
		Target:        a.Target,
		ExecutedAt:    a.ExecutedAt,
		Success:       a.Success,
		Output:        a.Output,
		ExposureScore: a.ExposureScore,
	}

	return json.Marshal(view)
}

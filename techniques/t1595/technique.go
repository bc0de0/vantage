package t1595

import (
	"errors"

	"vantage/techniques"
)

// ============================================================================
// TECHNIQUE T1595 — ACTIVE SCANNING (DECISION ONLY)
//
// This file:
// - Performs NO scanning
// - Executes NO tools
// - Has NO side effects
//
// It exists ONLY to declare which Action Classes may be considered.
// ============================================================================

type Technique1595 struct{}

// ID returns the canonical technique identifier.
func (t *Technique1595) ID() string {
	return "T1595"
}

// Description returns a non-operational summary.
func (t *Technique1595) Description() string {
	return "Active scanning of target surface to identify reachable assets and services"
}

// Resolve determines admissible Action Classes.
//
// PURE FUNCTION GUARANTEES:
// - No I/O
// - No global state
// - Deterministic output
func (t *Technique1595) Resolve(
	input techniques.ResolveInput,
) (*techniques.Resolution, error) {

	// -----------------------------------------------------------------
	// Structural sanity check
	// -----------------------------------------------------------------
	if input.TechniqueID != t.ID() {
		return nil, errors.New("techniques/t1595: technique ID mismatch")
	}

	res := &techniques.Resolution{
		TechniqueID: t.ID(),
	}

	allowed := []string{}
	excluded := []techniques.ExcludedActionClass{}

	// -----------------------------------------------------------------
	// Intent domain gate
	// -----------------------------------------------------------------
	intentOK := false
	for _, d := range input.AllowedIntentDomains {
		if d == "discovery" || d == "enumeration" {
			intentOK = true
			break
		}
	}

	if !intentOK {
		res.Rationale = "Intent does not permit discovery or enumeration"
		return res, nil
	}

	// -----------------------------------------------------------------
	// ROE category gate
	// -----------------------------------------------------------------
	roeOK := false
	for _, c := range input.AllowedROECategories {
		if c == "active_non_invasive" {
			roeOK = true
			break
		}
	}

	if !roeOK {
		res.Rationale = "ROE does not permit active non-invasive actions"
		return res, nil
	}

	// -----------------------------------------------------------------
	// Baseline allowed classes
	// -----------------------------------------------------------------
	allowed = append(allowed, "AC-03") // Reachability validation

	switch input.ExposureBudget {
	case "medium", "high":
		allowed = append(allowed, "AC-02")
	case "low":
		excluded = append(excluded, techniques.ExcludedActionClass{
			ActionClassID: "AC-02",
			Reason:        "Exposure budget insufficient",
		})
	}

	if input.ExposureBudget == "high" {
		allowed = append(allowed, "AC-04", "AC-05")
	} else {
		excluded = append(excluded,
			techniques.ExcludedActionClass{
				ActionClassID: "AC-04",
				Reason:        "Exposure budget insufficient",
			},
			techniques.ExcludedActionClass{
				ActionClassID: "AC-05",
				Reason:        "Exposure budget insufficient",
			},
		)
	}

	// -----------------------------------------------------------------
	// Hard exclusions (documented for audit)
	// -----------------------------------------------------------------
	for _, id := range []string{
		"AC-06", "AC-07", "AC-08", "AC-09", "AC-10",
		"AC-11", "AC-12", "AC-13", "AC-14", "AC-15",
	} {
		excluded = append(excluded, techniques.ExcludedActionClass{
			ActionClassID: id,
			Reason:        "Out of scope for T1595 by doctrine",
		})
	}

	res.AllowedActionClasses = allowed
	res.ExcludedActionClasses = excluded
	res.Rationale = "Resolved admissible action classes for active scanning"

	return res, nil
}

// -----------------------------------------------------------------------------
// INIT-TIME REGISTRATION
// -----------------------------------------------------------------------------
func init() {
	techniques.Register(&Technique1595{})
}

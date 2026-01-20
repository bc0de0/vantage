package intent

import (
	"errors"
	"fmt"
	"time"
)

// -----------------------------------------------------------------------------
// Intent Contract — DECLARED OPERATOR INTENT (SECURITY CRITICAL)
//
// This file defines the canonical structure used to declare WHAT an operator
// intends to do, BEFORE any execution occurs.
//
// Intent is:
// - Declarative (not procedural)
// - Explicit (nothing implied)
// - Binding (used for ROE + audit)
//
// No execution may occur without a validated Intent Contract.
//
// DESIGN PRINCIPLES:
//
// 1. DECLARATION OVER INFERENCE
//    Operators must explicitly state intent.
//    The system must never "guess".
//
// 2. MINIMUM NECESSARY POWER
//    Intent must constrain scope, time, and techniques.
//
// 3. AUDITABILITY
//    An intent contract must be reviewable by non-engineers.
//
// 4. IMMUTABILITY AT RUNTIME
//    Once validated, intent must not change during execution.
// -----------------------------------------------------------------------------

// Contract represents a validated declaration of operator intent.
//
// A Contract is loaded BEFORE execution and used to:
// - Gate ROE
// - Scope execution
// - Anchor audit trails
//
// IMPORTANT:
// - Fields are intentionally explicit.
// - Optionality is avoided where possible.
// - Defaults are NOT assumed.
type Contract struct {

	// CampaignID uniquely identifies the engagement or operation.
	//
	// This value:
	// - Appears in evidence
	// - Appears in reports
	// - Anchors audit trails
	//
	// It MUST be:
	// - Non-empty
	// - Stable for the lifetime of the campaign
	CampaignID string

	// Objective describes the high-level goal of the campaign.
	//
	// This is NOT used for execution logic.
	// It exists for:
	// - Human review
	// - Legal context
	// - Reporting
	//
	// Examples:
	// - "Validate exposure of development network services"
	// - "Confirm effectiveness of credential hygiene controls"
	Objective string

	// AllowedTechniques is the explicit list of techniques
	// the operator is permitted to execute.
	//
	// IMPORTANT:
	// - This is an allow list, not a preference list.
	// - Techniques NOT listed here are forbidden.
	//
	// This list is intersected with ROE at runtime.
	AllowedTechniques []string

	// Targets defines the explicit scope of execution.
	//
	// Each entry represents a single, executor-validated target.
	// Interpretation is technique-specific.
	//
	// Examples:
	// - "10.10.0.5"
	// - "dev-db.internal"
	//
	// Wildcards and ranges are intentionally NOT supported in v0.x.
	Targets []string

	// NotBefore defines the earliest time execution is permitted.
	//
	// Evaluated in UTC.
	NotBefore time.Time

	// NotAfter defines the latest time execution is permitted.
	//
	// Evaluated in UTC.
	NotAfter time.Time
}

// Validate performs strict validation of the intent contract.
//
// This function MUST be called immediately after loading a contract
// and BEFORE any execution logic.
//
// VALIDATION RULES:
// - Fail fast
// - Fail closed
// - No defaults
//
// A contract that fails validation MUST NOT be used.
func (c *Contract) Validate() error {

	// -----------------------------------------------------------------
	// 1. Structural Validation
	// -----------------------------------------------------------------

	if c == nil {
		return errors.New("intent contract is nil")
	}

	if c.CampaignID == "" {
		return errors.New("intent contract missing campaign_id")
	}

	if c.Objective == "" {
		return errors.New("intent contract missing objective")
	}

	// -----------------------------------------------------------------
	// 2. Technique Scope Validation
	// -----------------------------------------------------------------

	if len(c.AllowedTechniques) == 0 {
		return errors.New("intent contract defines no allowed techniques")
	}

	for _, t := range c.AllowedTechniques {
		if t == "" {
			return errors.New("intent contract contains empty technique ID")
		}
	}

	// -----------------------------------------------------------------
	// 3. Target Scope Validation
	// -----------------------------------------------------------------

	if len(c.Targets) == 0 {
		return errors.New("intent contract defines no targets")
	}

	for _, target := range c.Targets {
		if target == "" {
			return errors.New("intent contract contains empty target")
		}
	}

	// -----------------------------------------------------------------
	// 4. Temporal Validation
	// -----------------------------------------------------------------

	if c.NotBefore.IsZero() {
		return errors.New("intent contract missing not_before timestamp")
	}

	if c.NotAfter.IsZero() {
		return errors.New("intent contract missing not_after timestamp")
	}

	if !c.NotAfter.After(c.NotBefore) {
		return errors.New("intent contract has invalid time window (not_after <= not_before)")
	}

	now := time.Now().UTC()
	if now.Before(c.NotBefore) || now.After(c.NotAfter) {
		return fmt.Errorf(
			"intent contract not currently valid (valid UTC %s – %s)",
			c.NotBefore.Format(time.RFC3339),
			c.NotAfter.Format(time.RFC3339),
		)
	}

	// -----------------------------------------------------------------
	// 5. PASS — Contract is valid
	// -----------------------------------------------------------------

	return nil
}

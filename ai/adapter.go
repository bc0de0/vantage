package ai

import (
	"encoding/json"
	"errors"
	"os"
)

// ============================================================================
// AI ADVISORY ADAPTER — AUTHORITATIVE SAFE ENTRY POINT
//
// FREEZE NOTICE (v0.x):
// - The AI ↔ Action Class contract is FROZEN.
// - AI may ONLY reference canonical Action Class IDs (AC-01 → AC-15).
// - Any deviation is treated as advisory failure and MUST be ignored
//   by the caller.
//
// SECURITY GUARANTEES:
// - No authority
// - No side effects
// - No execution capability
// - Safe to disable entirely
// - Safe to ignore on failure
//
// AI FAILURE MUST NEVER BLOCK CORE LOGIC.
// ============================================================================

// -----------------------------------------------------------------------------
// Advise — SINGLE SAFE ENTRY POINT
// -----------------------------------------------------------------------------
//
// Advise performs AI-assisted Action Class classification.
//
// DOCTRINAL RULES:
// - Advisory ONLY
// - No execution recommendation
// - No new Action Classes
// - No authority
//
// FAILURE SEMANTICS:
// - Any error returned here MUST be ignored by the caller.
// - Core logic MUST proceed without AI.
func Advise(input AdvisoryInput) (*AdvisoryOutput, error) {

	// -----------------------------------------------------------------
	// 1. HARD MODE GATE — PREVENT SILENT AUTHORITY DRIFT
	// -----------------------------------------------------------------

	if os.Getenv("VANTAGE_AI_MODE") != "advisory_only" {
		return nil, errors.New("AI mode is not advisory_only")
	}

	// -----------------------------------------------------------------
	// 2. STRUCTURED INPUT SERIALIZATION
	// -----------------------------------------------------------------

	payload, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return nil, err
	}

	// -----------------------------------------------------------------
	// 3. PROMPT CONSTRUCTION (DOCTRINAL)
	// -----------------------------------------------------------------

	fullPrompt := advisoryPrompt +
		"\n\nINPUT:\n" +
		string(payload)

	// -----------------------------------------------------------------
	// 4. MODEL INVOCATION (NON-AUTHORITATIVE)
	// -----------------------------------------------------------------

	raw, err := callTogether(fullPrompt)
	if err != nil {
		return nil, err
	}

	// -----------------------------------------------------------------
	// 5. OUTPUT PARSING
	// -----------------------------------------------------------------

	var output AdvisoryOutput
	if err := json.Unmarshal([]byte(raw), &output); err != nil {
		return nil, err
	}

	// -----------------------------------------------------------------
	// 6. MANDATORY DISCLAIMER ENFORCEMENT
	// -----------------------------------------------------------------

	if output.Disclaimer != "Advisory output only. Human decision required." {
		return nil, errors.New("AI output missing mandatory disclaimer")
	}

	// -----------------------------------------------------------------
	// 7. ACTION CLASS CONTRACT VALIDATION (FREEZE GUARD)
	// -----------------------------------------------------------------

	// Validate suggested action classes
	for _, s := range output.Suggested {
		if err := validateActionClassIDs([]string{s.ID}); err != nil {
			return nil, err
		}
	}

	// Validate excluded action classes
	for _, e := range output.Excluded {
		if err := validateActionClassIDs([]string{e.ID}); err != nil {
			return nil, err
		}
	}

	// Validate novelty closest matches
	for _, n := range output.Novelty {
		if err := validateActionClassIDs(n.ClosestMatches); err != nil {
			return nil, err
		}
	}

	// -----------------------------------------------------------------
	// 8. SUCCESSFUL ADVISORY OUTPUT
	// -----------------------------------------------------------------

	return &output, nil
}

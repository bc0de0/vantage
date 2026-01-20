package ai

import (
	"encoding/json"
	"errors"
	"os"
)

// -----------------------------------------------------------------------------
// AI ADVISORY ADAPTER — SINGLE SAFE ENTRY POINT
//
// This is the ONLY function VANTAGE core may call.
//
// GUARANTEES:
// - No side effects
// - No authority
// - Safe to ignore
// - Safe to disable
//
// AI FAILURE MUST NEVER BLOCK CORE LOGIC.
// -----------------------------------------------------------------------------

// Advise performs AI-assisted Action Class classification.
//
// If AI is unavailable or violates constraints, this function
// returns an error and MUST be ignored by the caller.
func Advise(input AdvisoryInput) (*AdvisoryOutput, error) {

	// Hard mode check — prevents silent drift
	if os.Getenv("VANTAGE_AI_MODE") != "advisory_only" {
		return nil, errors.New("AI mode is not advisory_only")
	}

	// Serialize structured input
	payload, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return nil, err
	}

	// Construct full prompt
	fullPrompt := advisoryPrompt +
		"\n\nINPUT:\n" +
		string(payload)

	// Call Together AI
	raw, err := callTogether(fullPrompt)
	if err != nil {
		return nil, err
	}

	// Parse AI output
	var output AdvisoryOutput
	if err := json.Unmarshal([]byte(raw), &output); err != nil {
		return nil, err
	}

	// Mandatory disclaimer enforcement
	if output.Disclaimer != "Advisory output only. Human decision required." {
		return nil, errors.New("AI output missing mandatory disclaimer")
	}

	return &output, nil
}

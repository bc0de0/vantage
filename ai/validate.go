package ai

import "fmt"

// ============================================================================
// AI ACTION CLASS VALIDATOR — FREEZE GUARD
//
// This file enforces the frozen Action Class contract.
// Any deviation causes AI output to be rejected.
//
// This validator has NO authority.
// It exists solely to detect drift.
// ============================================================================

// CanonicalActionClasses is the frozen set of allowed Action Class IDs.
var CanonicalActionClasses = map[string]struct{}{
	"AC-01": {}, "AC-02": {}, "AC-03": {}, "AC-04": {}, "AC-05": {},
	"AC-06": {}, "AC-07": {}, "AC-08": {}, "AC-09": {}, "AC-10": {},
	"AC-11": {}, "AC-12": {}, "AC-13": {}, "AC-14": {}, "AC-15": {},
}

// validateActionClassIDs ensures AI only references canonical Action Classes.
func validateActionClassIDs(ids []string) error {
	for _, id := range ids {
		if _, ok := CanonicalActionClasses[id]; !ok {
			return fmt.Errorf("AI referenced non-canonical action class: %s", id)
		}
	}
	return nil
}

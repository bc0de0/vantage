package state

import "fmt"

// OperationPhase represents the currently active offensive lifecycle phase.
type OperationPhase string

const (
	PhaseRecon           OperationPhase = "RECON"
	PhaseInitialAccess   OperationPhase = "INITIAL_ACCESS"
	PhasePersistence     OperationPhase = "PERSISTENCE"
	PhasePrivEsc         OperationPhase = "PRIV_ESC"
	PhaseLateralMovement OperationPhase = "LATERAL_MOVEMENT"
	PhaseC2              OperationPhase = "C2"
	PhaseObjective       OperationPhase = "OBJECTIVE"
	PhaseExfil           OperationPhase = "EXFIL"
)

var phaseOrder = []OperationPhase{
	PhaseRecon,
	PhaseInitialAccess,
	PhasePersistence,
	PhasePrivEsc,
	PhaseLateralMovement,
	PhaseC2,
	PhaseObjective,
	PhaseExfil,
}

// Validate checks whether a phase is a known lifecycle value.
func (p OperationPhase) Validate() error {
	for _, candidate := range phaseOrder {
		if candidate == p {
			return nil
		}
	}

	return fmt.Errorf("invalid operation phase: %s", p)
}

// Next returns the next legal phase in the lifecycle.
// The second return value is false when already at final phase.
func (p OperationPhase) Next() (OperationPhase, bool) {
	for i, candidate := range phaseOrder {
		if candidate != p {
			continue
		}
		if i+1 >= len(phaseOrder) {
			return p, false
		}
		return phaseOrder[i+1], true
	}

	return p, false
}

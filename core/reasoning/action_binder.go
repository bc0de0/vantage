package reasoning

import (
	"fmt"
	"sync"
	"time"

	"vantage/core/state"
)

// Evidence aliases normalized execution evidence used by the reasoning engine.
type Evidence = EvidenceEvent

// ActionBinder binds action classes to reasoning hypotheses and graph mutation.
type ActionBinder interface {
	BindActionClasses([]ActionClass)
	MatchAndGenerate(*Graph, *state.State) ([]Hypothesis, error)
	ApplyAction(*Graph, ActionClass, Evidence) error
}

// DefaultActionBinder provides deterministic action-class-driven hypothesis and graph updates.
type DefaultActionBinder struct {
	mu      sync.RWMutex
	classes map[string]ActionClass
}

// NewDefaultActionBinder creates an empty action binder.
func NewDefaultActionBinder() *DefaultActionBinder {
	return &DefaultActionBinder{classes: make(map[string]ActionClass)}
}

// BindActionClasses replaces the loaded action class registry.
func (b *DefaultActionBinder) BindActionClasses(classes []ActionClass) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.classes = make(map[string]ActionClass, len(classes))
	for _, class := range classes {
		if class.ID == "" {
			continue
		}
		b.classes[class.ID] = class
	}
}

// MatchAndGenerate creates deterministic hypotheses when action-class preconditions match graph state.
func (b *DefaultActionBinder) MatchAndGenerate(graph *Graph, st *state.State) ([]Hypothesis, error) {
	if graph == nil || st == nil {
		return nil, nil
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	out := make([]Hypothesis, 0, len(b.classes))
	currentPhase := phaseForState(st)
	for _, ac := range b.classes {
		if ac.Phase != currentPhase {
			continue
		}
		if !MatchPatterns(graph, ac.Preconditions) {
			continue
		}
		out = append(out, Hypothesis{
			ID:            fmt.Sprintf("hyp-ac-%s", ac.ID),
			ActionClassID: ac.ID,
			Statement:     fmt.Sprintf("Action class %s is feasible in %s", ac.Name, ac.Phase),
			Confidence:    0.5 + ac.ConfidenceBoost,
			DerivedFrom:   []string{},
		})
	}

	return out, nil
}

func phaseForState(st *state.State) state.OperationPhase {
	if st == nil {
		return state.PhaseRecon
	}
	switch {
	case st.Executions() >= 5:
		return state.PhaseObjective
	case st.Executions() >= 3:
		return state.PhaseLateralMovement
	default:
		return state.PhaseRecon
	}
}

// ApplyAction mutates the graph using action-class production semantics.
func (b *DefaultActionBinder) ApplyAction(graph *Graph, ac ActionClass, evidence Evidence) error {
	if graph == nil {
		return fmt.Errorf("graph is nil")
	}

	now := time.Now().UTC()
	evidenceNodeID := fmt.Sprintf("ev-%d-%s", now.UnixNano(), evidence.TechniqueID)
	graph.UpsertNode(&Node{ID: evidenceNodeID, Type: NodeTypeEvidence, Label: fmt.Sprintf("%s@%s", evidence.TechniqueID, evidence.Target)})

	producedIDs := make([]string, 0, len(ac.ProducesNodes))
	for idx, nodeType := range ac.ProducesNodes {
		nodeID := fmt.Sprintf("ac-%s-%d-%d", ac.ID, idx, now.UnixNano())
		graph.UpsertNode(&Node{ID: nodeID, Type: nodeType, Label: fmt.Sprintf("%s produced %s", ac.ID, nodeType)})
		producedIDs = append(producedIDs, nodeID)
	}

	if len(producedIDs) == 0 {
		return nil
	}
	for _, edgeType := range ac.ProducesEdges {
		if err := graph.AddEdge(&Edge{From: evidenceNodeID, To: producedIDs[0], Type: edgeType, Weight: 1.0}); err != nil {
			return err
		}
	}

	return nil
}

// ActionClass resolves a loaded action class by identifier.
func (b *DefaultActionBinder) ActionClass(id string) (ActionClass, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	ac, ok := b.classes[id]
	return ac, ok
}

var _ ActionBinder = (*DefaultActionBinder)(nil)

package reasoning

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"vantage/core/state"
	"vantage/techniques"
)

// Decision is the structured planning output for executor consumption.
type Decision struct {
	Selected  RankedAction
	Ranked    []RankedAction
	CreatedAt time.Time
}

// Engine orchestrates the full reasoning lifecycle over an in-memory graph.
type Engine struct {
	mu       sync.RWMutex
	graph    *Graph
	registry *effectRegistry
	planner  *Planner
	expander HypothesisExpander
	state    *state.State
}

// NewEngine constructs a reasoning engine with default technique effects.
func NewEngine(expander HypothesisExpander) *Engine {
	registry := newEffectRegistry()
	for _, id := range techniques.List() {
		registry.RegisterTechniqueEffect(TechniqueEffect{
			TechniqueID: id,
			Impact:      0.6,
			Risk:        0.4,
			Stealth:     0.5,
			Produces:    []string{"generic_evidence"},
		})
	}
	planner := NewPlanner(registry, DefaultTechniqueScoreWeights())
	return &Engine{
		graph:    NewGraph(),
		registry: registry,
		planner:  planner,
		expander: expander,
	}
}

// Graph returns the underlying operational graph.
func (e *Engine) Graph() *Graph {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.graph
}

// RegisterTechniqueEffect registers or updates effect metadata for a technique.
func (e *Engine) RegisterTechniqueEffect(effect TechniqueEffect) {
	e.registry.RegisterTechniqueEffect(effect)
}

// EffectForTechnique returns effect metadata for a technique.
func (e *Engine) EffectForTechnique(techniqueID string) (TechniqueEffect, bool) {
	return e.registry.EffectForTechnique(techniqueID)
}

// KnownTechniques returns all techniques known to the reasoning registry.
func (e *Engine) KnownTechniques() []string {
	return e.registry.KnownTechniques()
}

// IngestEvidence updates graph state from executor evidence.
func (e *Engine) IngestEvidence(event EvidenceEvent) error {
	if event.TechniqueID == "" || event.Target == "" {
		return fmt.Errorf("evidence event missing technique or target")
	}
	nodeID := fmt.Sprintf("ev-%d-%s", time.Now().UTC().UnixNano(), event.TechniqueID)
	e.graph.UpsertNode(&Node{
		ID:    nodeID,
		Type:  NodeTypeEvidence,
		Label: fmt.Sprintf("%s@%s", event.TechniqueID, event.Target),
		Metadata: map[string]string{
			"success": fmt.Sprintf("%t", event.Success),
			"target":  event.Target,
		},
	})
	return nil
}

// PlanNextAction runs hypothesis generation, scoring, and action selection.
func (e *Engine) PlanNextAction(query PlannerQuery) (*Decision, error) {
	hypotheses := GenerateHypotheses(e.graph)
	if e.expander != nil {
		aiHypotheses, err := e.expander.Expand(e.graph, e.state)
		if err == nil {
			hypotheses = append(hypotheses, aiHypotheses...)
		}
	}
	for _, h := range hypotheses {
		e.graph.UpsertNode(&Node{ID: h.ID, Type: NodeTypeHypothesis, Label: h.Statement, Metadata: map[string]string{"confidence": fmt.Sprintf("%.2f", h.Confidence)}})
		for _, support := range h.SupportingNodeIDs {
			_ = e.graph.AddEdge(&Edge{From: support, To: h.ID, Type: EdgeTypeSupports, Weight: h.Confidence})
		}
	}

	ranked := e.planner.RankedActions(query)
	if len(ranked) == 0 {
		return nil, fmt.Errorf("no ranked actions available")
	}
	decision := &Decision{Selected: ranked[0], Ranked: ranked, CreatedAt: time.Now().UTC()}

	selectedNodeID := fmt.Sprintf("tech-%s", decision.Selected.TechniqueID)
	e.graph.UpsertNode(&Node{ID: selectedNodeID, Type: NodeTypeTechnique, Label: decision.Selected.TechniqueID})
	for _, h := range hypotheses {
		_ = e.graph.AddEdge(&Edge{From: h.ID, To: selectedNodeID, Type: EdgeTypeEnables, Weight: h.Confidence})
	}

	return decision, nil
}

// DOT returns Graphviz DOT output for the current reasoning graph.
func (e *Engine) DOT() string {
	return e.graph.ToDOT()
}

type effectRegistry struct {
	mu      sync.RWMutex
	effects map[string]TechniqueEffect
}

func newEffectRegistry() *effectRegistry {
	return &effectRegistry{effects: make(map[string]TechniqueEffect)}
}

func (r *effectRegistry) RegisterTechniqueEffect(effect TechniqueEffect) {
	if effect.TechniqueID == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.effects[effect.TechniqueID] = effect
}

func (r *effectRegistry) EffectForTechnique(techniqueID string) (TechniqueEffect, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	effect, ok := r.effects[techniqueID]
	return effect, ok
}

func (r *effectRegistry) KnownTechniques() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := make([]string, 0, len(r.effects))
	for id := range r.effects {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

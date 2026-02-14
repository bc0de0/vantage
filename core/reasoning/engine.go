package reasoning

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"vantage/core/evidence"
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
	mu               sync.RWMutex
	graph            *Graph
	registry         *effectRegistry
	planner          *Planner
	expander         HypothesisExpander
	actionBinder     ActionBinder
	state            *state.State
	cycle            CycleConfig
	attackPathConfig AttackPathConfig
}

// TechniqueExecutor executes a selected technique against a target.
type TechniqueExecutor interface {
	Run(ctx context.Context, techniqueID string, target string) (*evidence.Artifact, error)
}

// CycleConfig holds execution wiring for a full reasoning cycle.
type CycleConfig struct {
	Target            string
	AllowedTechniques []string
	Executor          TechniqueExecutor
	Timeout           time.Duration
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
	binder := NewDefaultActionBinder()
	if classes, err := LoadActionClassesFromDir("action-classes-normalized"); err == nil {
		binder.BindActionClasses(classes)
	}
	return &Engine{
		graph:            NewGraph(),
		registry:         registry,
		planner:          planner,
		expander:         expander,
		actionBinder:     binder,
		attackPathConfig: DefaultAttackPathConfig(),
	}
}

// Graph returns the underlying operational graph.
func (e *Engine) Graph() *Graph {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.graph
}

// BindActionClasses replaces the action class set driving deterministic reasoning.
func (e *Engine) BindActionClasses(classes []ActionClass) {
	if e.actionBinder != nil {
		e.actionBinder.BindActionClasses(classes)
	}
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

// GenerateHypotheses creates deterministic hypotheses from evidence and action-class matching.
// Action classes act as graph rules: when phase and preconditions match, the engine emits
// deterministic hypotheses anchored to the matching action class IDs.
func (e *Engine) GenerateHypotheses() []Hypothesis {
	hypotheses := GenerateHypotheses(e.graph)
	if e.actionBinder != nil {
		matched, err := e.actionBinder.MatchAndGenerate(e.graph, e.state)
		if err == nil {
			hypotheses = append(hypotheses, matched...)
		}
	}
	return hypotheses
}

// PlanNextAction runs hypothesis generation, scoring, and action selection.
func (e *Engine) PlanNextAction(query PlannerQuery) (*Decision, error) {
	hypotheses := e.GenerateHypotheses()
	if e.expander != nil {
		aiHypotheses, err := e.expander.Expand(e.graph, e.state)
		if err == nil {
			hypotheses = append(hypotheses, aiHypotheses...)
		}
	}
	for _, h := range hypotheses {
		e.graph.UpsertNode(&Node{ID: h.ID, Type: NodeTypeHypothesis, Label: h.Statement, Metadata: map[string]string{"confidence": fmt.Sprintf("%.2f", h.Confidence), "action_class": h.ActionClassID}})
		for _, support := range h.SupportingNodeIDs {
			_ = e.graph.AddEdge(&Edge{From: support, To: h.ID, Type: EdgeTypeSupports, Weight: h.Confidence})
		}
	}

	ranked := e.planner.RankedActions(query)
	if e.state != nil {
		phase := phaseForState(e.state)
		if phase == state.PhaseLateralMovement || phase == state.PhaseObjective || phase == state.PhaseC2 {
			if paths, err := e.ExpandAttackPaths(e.state); err == nil && len(paths) > 0 {
				enrichRankedActionsWithPaths(ranked, paths)
			}
		}
	}
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

// ConfigureCycle configures runtime dependencies for RunCycle.
func (e *Engine) ConfigureCycle(cfg CycleConfig) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.cycle = cfg
}

// RunCycle executes one deterministic reasoning + execution cycle.
func (e *Engine) RunCycle(state *state.State) (*Decision, error) {
	e.mu.Lock()
	e.state = state
	cfg := e.cycle
	e.mu.Unlock()

	if cfg.Target == "" {
		return nil, fmt.Errorf("run cycle target is required")
	}
	if cfg.Executor == nil {
		return nil, fmt.Errorf("run cycle executor is required")
	}

	decision, err := e.PlanNextAction(PlannerQuery{
		Target:            cfg.Target,
		AllowedTechniques: cfg.AllowedTechniques,
		TopN:              1,
	})
	if err != nil {
		return nil, err
	}

	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	artifact, execErr := cfg.Executor.Run(ctx, decision.Selected.TechniqueID, cfg.Target)
	if artifact != nil {
		event := EvidenceEvent{TechniqueID: artifact.TechniqueID, Target: artifact.Target, Success: artifact.Success, Output: artifact.Output, Artifact: artifact}
		applied := false
		if binder, ok := e.actionBinder.(*DefaultActionBinder); ok && decision.Selected.ActionClassID != "" {
			if ac, found := binder.ActionClass(decision.Selected.ActionClassID); found {
				if err := e.actionBinder.ApplyAction(e.graph, ac, event); err == nil {
					applied = true
				}
			}
		}
		if !applied {
			_ = e.IngestEvidence(event)
		}
	}

	if execErr != nil {
		return decision, execErr
	}
	return decision, nil
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

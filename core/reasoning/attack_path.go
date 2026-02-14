package reasoning

import (
	"fmt"
	"sort"

	"vantage/core/state"
)

// AttackPath represents an ordered, scored sequence of hypotheses that model a feasible attack chain.
type AttackPath struct {
	Steps     []Hypothesis
	Score     float64
	Risk      float64
	Objective NodeType
	Valid     bool
}

// AttackPathConfig controls search depth, pruning, scoring, and objective detection.
type AttackPathConfig struct {
	MaxDepth           int
	RiskThreshold      float64
	DepthPenalty       float64
	ConfidenceWeight   float64
	StartNodeTypes     []NodeType
	ObjectiveNodeTypes []NodeType
	ROEPolicy          func(ac ActionClass, graph *Graph, st *state.State) bool
}

// DefaultAttackPathConfig returns conservative attack-path search defaults.
func DefaultAttackPathConfig() AttackPathConfig {
	return AttackPathConfig{
		MaxDepth:           4,
		RiskThreshold:      2.0,
		DepthPenalty:       0.1,
		ConfidenceWeight:   0.25,
		StartNodeTypes:     []NodeType{NodeTypeEvidence, NodeTypeHypothesis, NodeTypeTechnique},
		ObjectiveNodeTypes: []NodeType{NodeTypeAttackPath, NodeTypeTechnique},
		ROEPolicy:          func(ActionClass, *Graph, *state.State) bool { return true },
	}
}

// ExpandAttackPaths computes feasible, scored attack paths from the current graph using virtual graph simulation.
func (e *Engine) ExpandAttackPaths(st *state.State) ([]AttackPath, error) {
	if e == nil || e.graph == nil {
		return nil, fmt.Errorf("engine or graph is nil")
	}

	cfg := e.attackPathConfig
	if cfg.MaxDepth <= 0 {
		cfg = DefaultAttackPathConfig()
	}
	if cfg.ROEPolicy == nil {
		cfg.ROEPolicy = func(ActionClass, *Graph, *state.State) bool { return true }
	}

	classes := e.boundActionClasses()
	if len(classes) == 0 {
		return nil, nil
	}

	seeded := e.seedNodeCount(cfg.StartNodeTypes)
	if seeded == 0 {
		return nil, nil
	}

	currentPhase := phaseForState(st)
	paths := make([]AttackPath, 0)
	seen := make(map[string]struct{})

	for _, root := range classes {
		if !phaseAllowed(currentPhase, root.Phase) {
			continue
		}
		gCopy := cloneGraph(e.graph)
		explorePathTree(gCopy, st, cfg, classes, []ActionClass{root}, &paths, seen)
	}

	sort.Slice(paths, func(i, j int) bool {
		if paths[i].Score == paths[j].Score {
			return len(paths[i].Steps) < len(paths[j].Steps)
		}
		return paths[i].Score > paths[j].Score
	})
	return paths, nil
}

// ConfigureAttackPathExpansion sets the engine attack-path search configuration.
func (e *Engine) ConfigureAttackPathExpansion(cfg AttackPathConfig) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.attackPathConfig = cfg
}

func explorePathTree(g *Graph, st *state.State, cfg AttackPathConfig, classes []ActionClass, stack []ActionClass, out *[]AttackPath, seen map[string]struct{}) {
	if len(stack) == 0 || len(stack) > cfg.MaxDepth {
		return
	}

	latest := stack[len(stack)-1]
	if !MatchPatterns(g, latest.Preconditions) || !cfg.ROEPolicy(latest, g, st) {
		return
	}
	if err := simulateAction(g, latest); err != nil {
		return
	}

	hyp := hypothesisForAction(latest, len(stack))
	steps := buildHypotheses(stack)
	steps[len(steps)-1] = hyp
	risk := cumulativeRisk(stack)
	if cfg.RiskThreshold > 0 && risk > cfg.RiskThreshold {
		return
	}

	if objective, ok := findObjective(cfg.ObjectiveNodeTypes, latest.ProducesNodes); ok {
		path := scorePath(steps, stack, classes, objective, cfg)
		key := pathKey(path)
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			*out = append(*out, path)
		}
		return
	}

	advanced := false
	currentPhase := phaseForState(st)
	for _, next := range classes {
		if !phaseAllowed(currentPhase, next.Phase) {
			continue
		}
		gNext := cloneGraph(g)
		nextStack := append(append([]ActionClass(nil), stack...), next)
		advanced = true
		explorePathTree(gNext, st, cfg, classes, nextStack, out, seen)
	}

	if !advanced {
		path := scorePath(steps, stack, classes, "", cfg)
		path.Valid = MatchPatterns(g, latest.Preconditions)
		if path.Valid {
			key := pathKey(path)
			if _, exists := seen[key]; !exists {
				seen[key] = struct{}{}
				*out = append(*out, path)
			}
		}
	}
}

func findObjective(objectiveNodeTypes []NodeType, produced []NodeType) (NodeType, bool) {
	for _, objective := range objectiveNodeTypes {
		for _, p := range produced {
			if objective == p {
				return objective, true
			}
		}
	}
	return "", false
}

func cumulativeRisk(classes []ActionClass) float64 {
	total := 0.0
	for _, ac := range classes {
		total += ac.RiskWeight
	}
	return total
}

func buildHypotheses(classes []ActionClass) []Hypothesis {
	out := make([]Hypothesis, 0, len(classes))
	for idx, ac := range classes {
		out = append(out, hypothesisForAction(ac, idx+1))
	}
	return out
}

func hypothesisForAction(ac ActionClass, idx int) Hypothesis {
	return Hypothesis{
		ID:            fmt.Sprintf("path-hyp-%s-%d", ac.ID, idx),
		ActionClassID: ac.ID,
		Statement:     fmt.Sprintf("Action class %s is feasible", ac.Name),
		Confidence:    0.5 + ac.ConfidenceBoost,
	}
}

func simulateAction(g *Graph, ac ActionClass) error {
	if g == nil {
		return fmt.Errorf("graph is nil")
	}
	base := fmt.Sprintf("sim-%s-%d", ac.ID, len(g.NodesByType(NodeTypeEvidence))+len(g.NodesByType(NodeTypeHypothesis))+len(g.NodesByType(NodeTypeTechnique))+len(g.NodesByType(NodeTypeAttackPath)))
	for idx, nodeType := range ac.ProducesNodes {
		g.UpsertNode(&Node{ID: fmt.Sprintf("%s-node-%d", base, idx), Type: nodeType, Label: fmt.Sprintf("simulated %s", ac.ID)})
	}
	nodes := g.NodesByType(NodeTypeEvidence)
	if len(nodes) == 0 {
		return nil
	}
	for idx, edgeType := range ac.ProducesEdges {
		toID := fmt.Sprintf("%s-edge-node-%d", base, idx)
		g.UpsertNode(&Node{ID: toID, Type: NodeTypeHypothesis, Label: "simulated edge target"})
		_ = g.AddEdge(&Edge{From: nodes[0].ID, To: toID, Type: edgeType, Weight: 1})
	}
	return nil
}

func cloneGraph(g *Graph) *Graph {
	if g == nil {
		return NewGraph()
	}
	copy := NewGraph()
	g.mu.RLock()
	defer g.mu.RUnlock()
	for id, node := range g.nodes {
		metadata := make(map[string]string, len(node.Metadata))
		for k, v := range node.Metadata {
			metadata[k] = v
		}
		copy.nodes[id] = &Node{ID: node.ID, Type: node.Type, Label: node.Label, CreatedAt: node.CreatedAt, Metadata: metadata}
	}
	for _, edge := range g.edges {
		copy.edges = append(copy.edges, &Edge{From: edge.From, To: edge.To, Type: edge.Type, Weight: edge.Weight, CreatedAt: edge.CreatedAt})
	}
	return copy
}

func phaseAllowed(current, candidate state.OperationPhase) bool {
	if current == candidate {
		return true
	}
	next, ok := current.Next()
	return ok && next == candidate
}

func pathKey(path AttackPath) string {
	ids := make([]string, 0, len(path.Steps))
	for _, step := range path.Steps {
		ids = append(ids, step.ActionClassID)
	}
	return fmt.Sprintf("%v|%s", ids, path.Objective)
}

func (e *Engine) boundActionClasses() []ActionClass {
	binder, ok := e.actionBinder.(*DefaultActionBinder)
	if !ok {
		return nil
	}
	return binder.Classes()
}

func (e *Engine) seedNodeCount(types []NodeType) int {
	total := 0
	for _, t := range types {
		total += len(e.graph.NodesByType(t))
	}
	return total
}

func enrichRankedActionsWithPaths(ranked []RankedAction, paths []AttackPath) {
	if len(ranked) == 0 || len(paths) == 0 {
		return
	}
	bestByAction := make(map[string]float64)
	for _, path := range paths {
		if len(path.Steps) == 0 {
			continue
		}
		actionID := path.Steps[0].ActionClassID
		if path.Score > bestByAction[actionID] {
			bestByAction[actionID] = path.Score
		}
	}
	for i := range ranked {
		if bonus, ok := bestByAction[ranked[i].ActionClassID]; ok {
			ranked[i].Score += bonus * 0.1
			ranked[i].Reason = fmt.Sprintf("%s path_bonus=%.2f", ranked[i].Reason, bonus*0.1)
		}
	}
	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].Score == ranked[j].Score {
			return ranked[i].TechniqueID < ranked[j].TechniqueID
		}
		return ranked[i].Score > ranked[j].Score
	})
}

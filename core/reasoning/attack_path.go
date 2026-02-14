package reasoning

import (
	"fmt"
	"sort"
	"strings"

	"vantage/core/state"
)

// AttackPath represents an ordered, scored sequence of hypotheses that model a feasible attack chain.
type AttackPath struct {
	Steps                   []Hypothesis
	Score                   float64
	Risk                    float64
	Objective               NodeType
	ObjectiveProximityScore float64
	Valid                   bool
}

// AttackPathConfig controls search depth, pruning, scoring, and objective detection.
type AttackPathConfig struct {
	MaxDepth           int
	BeamWidth          int
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
		BeamWidth:          25,
		RiskThreshold:      2.0,
		DepthPenalty:       0.1,
		ConfidenceWeight:   0.25,
		StartNodeTypes:     []NodeType{NodeTypeEvidence, NodeTypeHypothesis, NodeTypeTechnique},
		ObjectiveNodeTypes: []NodeType{NodeTypeAttackPath, NodeTypeTechnique},
		ROEPolicy:          func(ActionClass, *Graph, *state.State) bool { return true },
	}
}

// CampaignProjectionState captures per-candidate virtual graph and phase progress during campaign projection.
type CampaignProjectionState struct {
	Graph         *graphSnapshot
	PhaseProgress []state.OperationPhase
}

type attackCandidate struct {
	graph *graphSnapshot
	stack []ActionClass
	score float64
	key   string
}

type graphSnapshot struct {
	nodeCounts map[NodeType]int
	edgeCounts map[EdgeType]int
}

type actionClassIndex struct {
	byRequiredNode map[NodeType][]ActionClass
	withoutReq     []ActionClass
}

func snapshotFromGraph(g *Graph) *graphSnapshot {
	s := &graphSnapshot{nodeCounts: map[NodeType]int{}, edgeCounts: map[EdgeType]int{}}
	if g == nil {
		return s
	}
	g.mu.RLock()
	defer g.mu.RUnlock()
	for _, n := range g.nodes {
		s.nodeCounts[n.Type]++
	}
	for _, e := range g.edges {
		s.edgeCounts[e.Type]++
	}
	return s
}

func (s *graphSnapshot) clone() *graphSnapshot {
	if s == nil {
		return &graphSnapshot{nodeCounts: map[NodeType]int{}, edgeCounts: map[EdgeType]int{}}
	}
	nodes := make(map[NodeType]int, len(s.nodeCounts))
	for k, v := range s.nodeCounts {
		nodes[k] = v
	}
	edges := make(map[EdgeType]int, len(s.edgeCounts))
	for k, v := range s.edgeCounts {
		edges[k] = v
	}
	return &graphSnapshot{nodeCounts: nodes, edgeCounts: edges}
}

func (s *graphSnapshot) hasNodeType(t NodeType) bool { return s != nil && s.nodeCounts[t] > 0 }
func (s *graphSnapshot) hasEdgeType(t EdgeType) bool { return s != nil && s.edgeCounts[t] > 0 }

func (s *graphSnapshot) applyAction(ac ActionClass) {
	if s == nil {
		return
	}
	for _, n := range ac.ProducesNodes {
		s.nodeCounts[n]++
	}
	for _, e := range ac.ProducesEdges {
		s.edgeCounts[e]++
	}
}

func (s *graphSnapshot) hash() string {
	if s == nil {
		return ""
	}
	nodeKeys := make([]string, 0, len(s.nodeCounts))
	for n, c := range s.nodeCounts {
		nodeKeys = append(nodeKeys, fmt.Sprintf("n:%s:%d", n, c))
	}
	sort.Strings(nodeKeys)
	edgeKeys := make([]string, 0, len(s.edgeCounts))
	for e, c := range s.edgeCounts {
		edgeKeys = append(edgeKeys, fmt.Sprintf("e:%s:%d", e, c))
	}
	sort.Strings(edgeKeys)
	return strings.Join(append(nodeKeys, edgeKeys...), "|")
}

func matchSnapshotPatterns(snapshot *graphSnapshot, patterns []GraphPattern) bool {
	if snapshot == nil {
		return false
	}
	for _, pattern := range patterns {
		for _, nodeType := range pattern.RequiredNodeTypes {
			if !snapshot.hasNodeType(nodeType) {
				return false
			}
		}
		for _, edgeType := range pattern.RequiredEdges {
			if !snapshot.hasEdgeType(edgeType) {
				return false
			}
		}
	}
	return true
}

func buildActionClassIndex(classes []ActionClass) actionClassIndex {
	idx := actionClassIndex{byRequiredNode: map[NodeType][]ActionClass{}}
	for _, ac := range classes {
		nodes := requiredNodes(ac.Preconditions)
		if len(nodes) == 0 {
			idx.withoutReq = append(idx.withoutReq, ac)
			continue
		}
		for _, n := range nodes {
			idx.byRequiredNode[n] = append(idx.byRequiredNode[n], ac)
		}
	}
	return idx
}

func requiredNodes(patterns []GraphPattern) []NodeType {
	set := map[NodeType]struct{}{}
	for _, p := range patterns {
		for _, n := range p.RequiredNodeTypes {
			set[n] = struct{}{}
		}
	}
	out := make([]NodeType, 0, len(set))
	for n := range set {
		out = append(out, n)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func (idx actionClassIndex) eligible(snapshot *graphSnapshot) []ActionClass {
	seen := map[string]ActionClass{}
	for _, ac := range idx.withoutReq {
		seen[ac.ID] = ac
	}
	for n := range snapshot.nodeCounts {
		if snapshot.nodeCounts[n] == 0 {
			continue
		}
		for _, ac := range idx.byRequiredNode[n] {
			seen[ac.ID] = ac
		}
	}
	out := make([]ActionClass, 0, len(seen))
	for _, ac := range seen {
		out = append(out, ac)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func projectCampaignState(base CampaignProjectionState, ac ActionClass) (CampaignProjectionState, error) {
	next := CampaignProjectionState{Graph: base.Graph.clone(), PhaseProgress: append(append([]state.OperationPhase(nil), base.PhaseProgress...), ac.Phase)}
	if next.Graph == nil {
		next.Graph = &graphSnapshot{nodeCounts: map[NodeType]int{}, edgeCounts: map[EdgeType]int{}}
	}
	if !matchSnapshotPatterns(next.Graph, ac.Preconditions) {
		return CampaignProjectionState{}, fmt.Errorf("preconditions do not match")
	}
	next.Graph.applyAction(ac)
	return next, nil
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
	if cfg.BeamWidth <= 0 {
		cfg.BeamWidth = DefaultAttackPathConfig().BeamWidth
	}
	if cfg.ROEPolicy == nil {
		cfg.ROEPolicy = func(ActionClass, *Graph, *state.State) bool { return true }
	}

	classes := e.boundActionClasses()
	if len(classes) == 0 {
		return nil, nil
	}

	idx := buildActionClassIndex(classes)
	unlockCache := map[string]float64{}

	seeded := e.seedNodeCount(cfg.StartNodeTypes)
	if seeded == 0 {
		return nil, nil
	}

	currentPhase := phaseForState(st)
	baseSnapshot := snapshotFromGraph(e.graph)
	beam := make([]attackCandidate, 0, len(classes))
	paths := make([]AttackPath, 0)
	seen := make(map[string]struct{})

	for _, root := range idx.eligible(baseSnapshot) {
		if !phaseAllowed(currentPhase, root.Phase) || !cfg.ROEPolicy(root, e.graph, st) || !matchSnapshotPatterns(baseSnapshot, root.Preconditions) {
			continue
		}
		stack := []ActionClass{root}
		scored := scorePathWithCache(buildHypotheses(stack), stack, classes, "", cfg, unlockCache, baseSnapshot.hash())
		beam = append(beam, attackCandidate{graph: baseSnapshot.clone(), stack: stack, score: scored.Score, key: actionStackKey(stack)})
	}
	beam = pruneAttackBeam(beam, cfg.BeamWidth)

	for depth := 1; depth <= cfg.MaxDepth && len(beam) > 0; depth++ {
		nextBeam := make([]attackCandidate, 0, len(beam)*len(classes))
		for _, cand := range beam {
			gCopy := cand.graph.clone()
			latest := cand.stack[len(cand.stack)-1]
			if !matchSnapshotPatterns(gCopy, latest.Preconditions) || !cfg.ROEPolicy(latest, e.graph, st) {
				continue
			}
			gCopy.applyAction(latest)

			risk := cumulativeRisk(cand.stack)
			riskLimit := cfg.RiskThreshold
			if riskLimit > 0 && riskLimit < 2.0 {
				riskLimit *= 0.9
			}
			if riskLimit > 0 && risk > riskLimit {
				continue
			}
			objective, reached := findObjective(cfg.ObjectiveNodeTypes, latest.ProducesNodes)
			path := scorePathWithCache(buildHypotheses(cand.stack), cand.stack, classes, objective, cfg, unlockCache, gCopy.hash())
			if reached {
				key := pathKey(path)
				if _, exists := seen[key]; !exists {
					seen[key] = struct{}{}
					paths = append(paths, path)
				}
			}

			if depth == cfg.MaxDepth {
				continue
			}
			for _, next := range idx.eligible(gCopy) {
				if !phaseAllowed(currentPhase, next.Phase) || actionInStack(cand.stack, next.ID) {
					continue
				}
				if !matchSnapshotPatterns(gCopy, next.Preconditions) {
					continue
				}
				nextStack := append(append([]ActionClass(nil), cand.stack...), next)
				nextScored := scorePathWithCache(buildHypotheses(nextStack), nextStack, classes, "", cfg, unlockCache, gCopy.hash())
				nextBeam = append(nextBeam, attackCandidate{graph: gCopy, stack: nextStack, score: nextScored.Score, key: actionStackKey(nextStack)})
			}
		}
		beam = pruneAttackBeam(nextBeam, cfg.BeamWidth)
	}

	sort.Slice(paths, func(i, j int) bool {
		if paths[i].Score == paths[j].Score {
			return len(paths[i].Steps) < len(paths[j].Steps)
		}
		return paths[i].Score > paths[j].Score
	})
	return paths, nil
}

func pruneAttackBeam(beam []attackCandidate, width int) []attackCandidate {
	sort.Slice(beam, func(i, j int) bool {
		if beam[i].score == beam[j].score {
			return beam[i].key < beam[j].key
		}
		return beam[i].score > beam[j].score
	})
	if len(beam) > width {
		return beam[:width]
	}
	return beam
}

func actionStackKey(stack []ActionClass) string {
	ids := make([]string, 0, len(stack))
	for _, step := range stack {
		ids = append(ids, step.ID)
	}
	return fmt.Sprintf("%v", ids)
}

// ConfigureAttackPathExpansion sets the engine attack-path search configuration.
func (e *Engine) ConfigureAttackPathExpansion(cfg AttackPathConfig) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.attackPathConfig = cfg
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

func actionInStack(stack []ActionClass, id string) bool {
	for _, step := range stack {
		if step.ID == id {
			return true
		}
	}
	return false
}

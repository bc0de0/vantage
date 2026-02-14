package reasoning

import (
	"fmt"
	"sort"
)

const (
	// ConfidenceWeight emphasizes confidence-backed paths.
	ConfidenceWeight = 0.8
	// FeasibilityWeight rewards paths with preconditions that are actually satisfiable.
	FeasibilityWeight = 1.1
	// UnlockFactor rewards paths that unlock additional action classes.
	UnlockFactor = 0.15
	// RiskThreshold defines when risk transitions to extreme-risk behavior.
	RiskThreshold = 2.0
	// SmallRiskFactor scales low-risk linear penalty.
	SmallRiskFactor = 0.4
	// DepthFactor penalizes deeper paths to prioritize practical chains.
	DepthFactor = 0.25
	// ObjectiveProximityFactor boosts chains that end by producing the requested objective.
	ObjectiveProximityFactor = 1.35
)

type TechniqueScoreWeights struct {
	ImpactWeight  float64
	RiskWeight    float64
	StealthWeight float64
}

func DefaultTechniqueScoreWeights() TechniqueScoreWeights {
	return TechniqueScoreWeights{ImpactWeight: 0.5, RiskWeight: 0.2, StealthWeight: 0.3}
}

func ScoreTechnique(effect TechniqueEffect, weights TechniqueScoreWeights) float64 {
	if weights == (TechniqueScoreWeights{}) {
		weights = DefaultTechniqueScoreWeights()
	}
	return (effect.Impact * weights.ImpactWeight) + ((1-effect.Risk)*weights.RiskWeight) + (effect.Stealth * weights.StealthWeight)
}

func scorePath(steps []Hypothesis, pathClasses []ActionClass, allClasses []ActionClass, objective NodeType, cfg AttackPathConfig) AttackPath {
	return scorePathWithCache(steps, pathClasses, allClasses, objective, cfg, nil, "")
}

func scorePathWithCache(steps []Hypothesis, pathClasses []ActionClass, allClasses []ActionClass, objective NodeType, cfg AttackPathConfig, unlockCache map[string]float64, graphHash string) AttackPath {
	totalConfidence := 0.0
	risk := 0.0
	for i := range pathClasses {
		risk += pathClasses[i].RiskWeight
		if i < len(steps) {
			totalConfidence += steps[i].Confidence
		}
	}

	averageConfidence := 0.0
	if len(steps) > 0 {
		averageConfidence = totalConfidence / float64(len(steps))
	}
	feasibilityScore := averageFeasibility(pathClasses)
	unlockBonus := unlockedActionCount(pathClasses, allClasses, unlockCache, graphHash) * UnlockFactor
	score := (averageConfidence * ConfidenceWeight) + (feasibilityScore * FeasibilityWeight) + unlockBonus - riskPenalty(risk) - (float64(len(steps)) * DepthFactor)

	proximity := objectiveProximity(pathClasses, objective, cfg)
	score += proximity
	if objective != "" {
		score *= ObjectiveProximityFactor
	}

	return AttackPath{Steps: steps, Score: score, Risk: risk, Objective: objective, ObjectiveProximityScore: proximity, Valid: true}
}

func objectiveProximity(pathClasses []ActionClass, objective NodeType, cfg AttackPathConfig) float64 {
	if len(pathClasses) == 0 {
		return 0
	}
	if objective != "" && producesNode(pathClasses[len(pathClasses)-1].ProducesNodes, objective) {
		return 1
	}
	if objective == "" && len(cfg.ObjectiveNodeTypes) > 0 {
		for _, o := range cfg.ObjectiveNodeTypes {
			if producesNode(pathClasses[len(pathClasses)-1].ProducesNodes, o) {
				return 0.9
			}
		}
	}
	return 1 / float64(len(pathClasses)+1)
}

func averageFeasibility(path []ActionClass) float64 { /* unchanged */
	if len(path) == 0 {
		return 0
	}
	nodeTypes := map[NodeType]struct{}{NodeTypeEvidence: {}}
	edgeTypes := map[EdgeType]struct{}{}
	totalRatio := 0.0
	for _, ac := range path {
		matched, total := matchedPreconditions(ac.Preconditions, nodeTypes, edgeTypes)
		ratio := 1.0
		if total > 0 {
			ratio = float64(matched) / float64(total)
		}
		totalRatio += ratio
		for _, n := range ac.ProducesNodes {
			nodeTypes[n] = struct{}{}
		}
		for _, e := range ac.ProducesEdges {
			edgeTypes[e] = struct{}{}
		}
	}
	return totalRatio / float64(len(path))
}

func unlockedActionCount(path []ActionClass, universe []ActionClass, cache map[string]float64, graphHash string) float64 {
	if len(path) == 0 || len(universe) == 0 {
		return 0
	}
	beforeNodes, beforeEdges := availabilityBeforeLast(path)
	afterNodes, afterEdges := availabilityAfterPath(path)
	selected := make(map[string]struct{}, len(path))
	for _, ac := range path {
		selected[ac.ID] = struct{}{}
	}
	cacheKey := ""
	if cache != nil {
		cacheKey = fmt.Sprintf("%s|%s|%s", graphHash, availabilityHash(beforeNodes, beforeEdges), availabilityHash(afterNodes, afterEdges))
		if v, ok := cache[cacheKey]; ok {
			return v
		}
	}
	unlocked := 0
	for _, candidate := range universe {
		if _, used := selected[candidate.ID]; used {
			continue
		}
		wasEligible := preconditionsEligible(candidate.Preconditions, beforeNodes, beforeEdges)
		isEligible := preconditionsEligible(candidate.Preconditions, afterNodes, afterEdges)
		if !wasEligible && isEligible {
			unlocked++
		}
	}
	if cache != nil {
		cache[cacheKey] = float64(unlocked)
	}
	return float64(unlocked)
}

func availabilityHash(nodes map[NodeType]struct{}, edges map[EdgeType]struct{}) string {
	nodeKeys := make([]string, 0, len(nodes))
	for n := range nodes {
		nodeKeys = append(nodeKeys, string(n))
	}
	sort.Strings(nodeKeys)
	edgeKeys := make([]string, 0, len(edges))
	for e := range edges {
		edgeKeys = append(edgeKeys, string(e))
	}
	sort.Strings(edgeKeys)
	return fmt.Sprintf("n=%v|e=%v", nodeKeys, edgeKeys)
}

func riskPenalty(risk float64) float64 {
	if risk > RiskThreshold {
		return risk * risk
	}
	return risk * SmallRiskFactor
}

func availabilityBeforeLast(path []ActionClass) (map[NodeType]struct{}, map[EdgeType]struct{}) {
	if len(path) <= 1 {
		return map[NodeType]struct{}{NodeTypeEvidence: {}}, map[EdgeType]struct{}{}
	}
	return availabilityAfterPath(path[:len(path)-1])
}

func availabilityAfterPath(path []ActionClass) (map[NodeType]struct{}, map[EdgeType]struct{}) {
	nodeTypes := map[NodeType]struct{}{NodeTypeEvidence: {}}
	edgeTypes := map[EdgeType]struct{}{}
	for _, ac := range path {
		for _, n := range ac.ProducesNodes {
			nodeTypes[n] = struct{}{}
		}
		for _, e := range ac.ProducesEdges {
			edgeTypes[e] = struct{}{}
		}
	}
	return nodeTypes, edgeTypes
}

func matchedPreconditions(patterns []GraphPattern, nodeTypes map[NodeType]struct{}, edgeTypes map[EdgeType]struct{}) (matched int, total int) {
	total = len(patterns)
	for _, pattern := range patterns {
		if patternSatisfied(pattern, nodeTypes, edgeTypes) {
			matched++
		}
	}
	return matched, total
}

func preconditionsEligible(patterns []GraphPattern, nodeTypes map[NodeType]struct{}, edgeTypes map[EdgeType]struct{}) bool {
	for _, pattern := range patterns {
		if !patternSatisfied(pattern, nodeTypes, edgeTypes) {
			return false
		}
	}
	return true
}

func patternSatisfied(pattern GraphPattern, nodeTypes map[NodeType]struct{}, edgeTypes map[EdgeType]struct{}) bool {
	for _, requiredNode := range pattern.RequiredNodeTypes {
		if _, ok := nodeTypes[requiredNode]; !ok {
			return false
		}
	}
	for _, requiredEdge := range pattern.RequiredEdges {
		if _, ok := edgeTypes[requiredEdge]; !ok {
			return false
		}
	}
	return true
}

package reasoning

import (
	"fmt"
	"sort"

	"vantage/techniques"
	"vantage/techniqueset"
)

// Planner ranks action candidates using technique effects.
type Planner struct {
	registry TechniqueEffectRegistry
	weights  TechniqueScoreWeights
}

// NewPlanner creates a planner with a technique effect registry.
func NewPlanner(registry TechniqueEffectRegistry, weights TechniqueScoreWeights) *Planner {
	return &Planner{registry: registry, weights: weights}
}

// RankedActions returns sorted candidates for a target.
func (p *Planner) RankedActions(query PlannerQuery) []RankedAction {
	if p == nil || p.registry == nil {
		return nil
	}
	techniquesToScore := query.AllowedTechniques
	if len(techniquesToScore) == 0 {
		techniquesToScore = p.registry.KnownTechniques()
	}

	out := make([]RankedAction, 0, len(techniquesToScore))
	for _, id := range techniquesToScore {
		effect, ok := p.registry.EffectForTechnique(id)
		if !ok {
			continue
		}
		score := ScoreTechnique(effect, p.weights)
		out = append(out, RankedAction{TechniqueID: id, ActionClassID: effect.ActionClassID, Target: query.Target, Score: score, Impact: effect.Impact, Risk: effect.Risk, Stealth: effect.Stealth, Reason: fmt.Sprintf("impact=%.2f risk=%.2f stealth=%.2f", effect.Impact, effect.Risk, effect.Stealth)})
	}

	sortRanked(out)
	if query.TopN > 0 && len(out) > query.TopN {
		out = out[:query.TopN]
	}
	return out
}

// RankedActionsForHypotheses scores techniques that map to each hypothesis action class.
func (p *Planner) RankedActionsForHypotheses(graph *Graph, target string, hypotheses []Hypothesis, classLookup func(string) (ActionClass, bool), topN int) []RankedAction {
	if graph == nil || classLookup == nil {
		return nil
	}
	candidates := map[string]RankedAction{}
	for _, h := range hypotheses {
		if h.ActionClassID == "" {
			continue
		}
		ac, ok := classLookup(h.ActionClassID)
		if !ok {
			continue
		}
		snapshot := techniqueGraphSnapshot(graph)
		for _, tech := range techniqueset.ForActionClass(h.ActionClassID) {
			if !tech.Evaluate(snapshot) {
				continue
			}
			score := (ac.ImpactWeight * tech.ImpactModifier()) + ((1 - ac.RiskWeight) * (1 - tech.RiskModifier()))
			ra := RankedAction{TechniqueID: tech.ID(), ActionClassID: h.ActionClassID, Target: target, Score: score, Impact: tech.ImpactModifier(), Risk: tech.RiskModifier(), Stealth: 1 - tech.RiskModifier(), Reason: fmt.Sprintf("ac_impact=%.2f ac_risk=%.2f tech_impact=%.2f tech_risk=%.2f", ac.ImpactWeight, ac.RiskWeight, tech.ImpactModifier(), tech.RiskModifier())}
			if existing, exists := candidates[ra.TechniqueID]; !exists || ra.Score > existing.Score {
				candidates[ra.TechniqueID] = ra
			}
		}
	}
	out := make([]RankedAction, 0, len(candidates))
	for _, c := range candidates {
		out = append(out, c)
	}
	sortRanked(out)
	if topN > 0 && len(out) > topN {
		out = out[:topN]
	}
	return out
}

func sortRanked(out []RankedAction) {
	sort.Slice(out, func(i, j int) bool {
		if out[i].Score == out[j].Score {
			return out[i].TechniqueID < out[j].TechniqueID
		}
		return out[i].Score > out[j].Score
	})
}

func techniqueGraphSnapshot(graph *Graph) *techniques.Graph {
	if graph == nil {
		return nil
	}
	return &techniques.Graph{
		EvidenceNodes:   len(graph.NodesByType(NodeTypeEvidence)),
		HypothesisNodes: len(graph.NodesByType(NodeTypeHypothesis)),
		TechniqueNodes:  len(graph.NodesByType(NodeTypeTechnique)),
		HasSupportsEdge: graph.HasEdgeType(EdgeTypeSupports),
		HasEnablesEdge:  graph.HasEdgeType(EdgeTypeEnables),
	}
}

package reasoning

import (
	"fmt"
	"sort"

	"vantage/core/state"
)

// AttackStep is a campaign-oriented projection of an action-class step.
type AttackStep struct {
	ActionClassID string
	Statement     string
	Confidence    float64
	Phase         state.OperationPhase
}

// Campaign is a strategic sequence of attack steps toward an objective.
type Campaign struct {
	Steps      []AttackStep
	Score      float64
	Risk       float64
	Objective  NodeType
	Confidence float64
}

// CampaignOptions controls campaign search bounds and pruning behavior.
type CampaignOptions struct {
	MaxDepth                int
	RiskTolerance           float64
	ConfidenceThreshold     float64
	BeamWidth               int
	TopN                    int
	ObjectiveBiasWeight     float64
	ObjectiveProximityScore float64
}

// DefaultCampaignOptions returns conservative deterministic planning defaults.
func DefaultCampaignOptions() CampaignOptions {
	return CampaignOptions{MaxDepth: 5, RiskTolerance: 2.0, ConfidenceThreshold: 0.55, BeamWidth: 25, TopN: 10, ObjectiveBiasWeight: 0.35}
}

type campaignCandidate struct {
	graph            *graphSnapshot
	actions          []ActionClass
	steps            []AttackStep
	score            float64
	risk             float64
	confidence       float64
	objectiveReached bool
	phaseProgress    []state.OperationPhase
	feasibility      float64
}

// PlanCampaign computes prioritized strategic campaigns for a requested objective node type.
func (e *Engine) PlanCampaign(objective NodeType, opts CampaignOptions) ([]Campaign, error) {
	if e == nil {
		return nil, fmt.Errorf("engine is nil")
	}
	if objective == "" {
		return nil, fmt.Errorf("objective is required")
	}

	cfg := normalizeCampaignOptions(opts)
	classes := e.boundActionClasses()
	if len(classes) == 0 {
		return nil, nil
	}
	sort.Slice(classes, func(i, j int) bool { return classes[i].ID < classes[j].ID })
	if e.graph == nil {
		return nil, fmt.Errorf("start graph is nil")
	}

	index := buildActionClassIndex(classes)
	unlockCache := map[string]float64{}
	beam := []campaignCandidate{{graph: snapshotFromGraph(e.graph)}}
	currentPhase := phaseForState(e.state)
	seen := map[string]struct{}{}
	campaigns := make([]Campaign, 0)

	for depth := 1; depth <= cfg.MaxDepth; depth++ {
		beam = pruneCampaignBeam(beam, cfg.BeamWidth)
		nextBeam := make([]campaignCandidate, 0, len(beam)*len(classes))
		for _, candidate := range beam {
			for _, action := range index.eligible(candidate.graph) {
				if !campaignPhaseAllowed(currentPhase, candidate.phaseProgress, action.Phase) || !matchSnapshotPatterns(candidate.graph, action.Preconditions) {
					continue
				}
				projected, ok := projectCampaignCandidate(candidate, action, classes, objective, cfg, unlockCache)
				if !ok {
					continue
				}
				nextBeam = append(nextBeam, projected)
				if projected.objectiveReached {
					campaign := Campaign{Steps: append([]AttackStep(nil), projected.steps...), Score: projected.score, Risk: projected.risk, Objective: objective, Confidence: projected.confidence}
					key := campaignKey(campaign)
					if _, exists := seen[key]; !exists {
						seen[key] = struct{}{}
						campaigns = append(campaigns, campaign)
					}
				}
			}
		}
		nextBeam = pruneCampaignBeam(nextBeam, cfg.BeamWidth)
		if len(nextBeam) == 0 {
			break
		}
		beam = nextBeam
	}

	sort.Slice(campaigns, func(i, j int) bool {
		if campaigns[i].Score == campaigns[j].Score {
			return campaignKey(campaigns[i]) < campaignKey(campaigns[j])
		}
		return campaigns[i].Score > campaigns[j].Score
	})
	if len(campaigns) > cfg.TopN {
		campaigns = campaigns[:cfg.TopN]
	}
	return campaigns, nil
}

func pruneCampaignBeam(beam []campaignCandidate, width int) []campaignCandidate {
	sort.Slice(beam, func(i, j int) bool {
		if beam[i].score == beam[j].score {
			return candidatePathKey(beam[i]) < candidatePathKey(beam[j])
		}
		return beam[i].score > beam[j].score
	})
	if len(beam) > width {
		return beam[:width]
	}
	return beam
}

func normalizeCampaignOptions(opts CampaignOptions) CampaignOptions {
	cfg := opts
	defaults := DefaultCampaignOptions()
	if cfg.MaxDepth <= 0 {
		cfg.MaxDepth = defaults.MaxDepth
	}
	if cfg.BeamWidth <= 0 {
		cfg.BeamWidth = defaults.BeamWidth
	}
	if cfg.RiskTolerance <= 0 {
		cfg.RiskTolerance = defaults.RiskTolerance
	}
	if cfg.ConfidenceThreshold <= 0 {
		cfg.ConfidenceThreshold = defaults.ConfidenceThreshold
	}
	if cfg.TopN <= 0 {
		cfg.TopN = defaults.TopN
	}
	if cfg.ObjectiveBiasWeight <= 0 {
		cfg.ObjectiveBiasWeight = defaults.ObjectiveBiasWeight
	}
	return cfg
}

func projectCampaignCandidate(candidate campaignCandidate, action ActionClass, classes []ActionClass, objective NodeType, cfg CampaignOptions, unlockCache map[string]float64) (campaignCandidate, bool) {
	proj, err := projectCampaignState(CampaignProjectionState{Graph: candidate.graph, PhaseProgress: candidate.phaseProgress}, action)
	if err != nil {
		return campaignCandidate{}, false
	}
	actions := append(append([]ActionClass(nil), candidate.actions...), action)
	risk := cumulativeRisk(actions)
	if risk > cfg.RiskTolerance {
		return campaignCandidate{}, false
	}

	steps := append(append([]AttackStep(nil), candidate.steps...), attackStepForAction(action, len(actions)))
	confidence := averageCampaignConfidence(steps)
	if confidence < cfg.ConfidenceThreshold {
		return campaignCandidate{}, false
	}
	feasibility := averageFeasibility(actions)
	if len(candidate.steps) > 0 && feasibility+1e-9 < candidate.feasibility {
		return campaignCandidate{}, false
	}

	reached := producesNode(action.ProducesNodes, objective)
	distance := objectiveDistance(actions, objective)
	proximity := objectiveProximityScore(distance, action, objective)
	hypSteps := hypothesesFromAttackSteps(steps)
	scored := scorePathWithCache(hypSteps, actions, classes, nodeTypeIf(reached, objective), DefaultAttackPathConfig(), unlockCache, proj.Graph.hash())
	scored.Score += proximity * cfg.ObjectiveBiasWeight

	return campaignCandidate{graph: proj.Graph, actions: actions, steps: steps, score: scored.Score, risk: risk, confidence: confidence, objectiveReached: reached, phaseProgress: proj.PhaseProgress, feasibility: feasibility}, true
}

func objectiveDistance(actions []ActionClass, objective NodeType) int {
	if len(actions) == 0 {
		return 0
	}
	for i := len(actions) - 1; i >= 0; i-- {
		if producesNode(actions[i].ProducesNodes, objective) {
			return len(actions) - i - 1
		}
	}
	return len(actions)
}

func objectiveProximityScore(distance int, action ActionClass, objective NodeType) float64 {
	if producesNode(action.ProducesNodes, objective) {
		return 1.0
	}
	supporting := 0.0
	for _, p := range action.Preconditions {
		for _, req := range p.RequiredNodeTypes {
			if req == objective {
				supporting += 0.5
			}
		}
	}
	return (1 / float64(distance+1)) + supporting
}

func attackStepForAction(ac ActionClass, idx int) AttackStep {
	return AttackStep{ActionClassID: ac.ID, Statement: fmt.Sprintf("Action class %s is feasible", ac.Name), Confidence: 0.5 + ac.ConfidenceBoost, Phase: ac.Phase}
}

func hypothesesFromAttackSteps(steps []AttackStep) []Hypothesis {
	out := make([]Hypothesis, 0, len(steps))
	for i, step := range steps {
		out = append(out, Hypothesis{ID: fmt.Sprintf("campaign-hyp-%s-%d", step.ActionClassID, i+1), ActionClassID: step.ActionClassID, Statement: step.Statement, Confidence: step.Confidence})
	}
	return out
}

func averageCampaignConfidence(steps []AttackStep) float64 {
	if len(steps) == 0 {
		return 0
	}
	total := 0.0
	for _, step := range steps {
		total += step.Confidence
	}
	return total / float64(len(steps))
}

func producesNode(nodes []NodeType, objective NodeType) bool {
	for _, node := range nodes {
		if node == objective {
			return true
		}
	}
	return false
}

func campaignPhaseAllowed(root state.OperationPhase, progress []state.OperationPhase, candidate state.OperationPhase) bool {
	if len(progress) == 0 {
		return phaseAllowed(root, candidate)
	}
	return phaseAllowed(progress[len(progress)-1], candidate)
}

func nodeTypeIf(ok bool, objective NodeType) NodeType {
	if ok {
		return objective
	}
	return ""
}

func campaignKey(c Campaign) string {
	ids := make([]string, 0, len(c.Steps))
	for _, step := range c.Steps {
		ids = append(ids, step.ActionClassID)
	}
	return fmt.Sprintf("%v|%s", ids, c.Objective)
}

func candidatePathKey(c campaignCandidate) string {
	ids := make([]string, 0, len(c.steps))
	for _, step := range c.steps {
		ids = append(ids, step.ActionClassID)
	}
	return fmt.Sprintf("%v", ids)
}

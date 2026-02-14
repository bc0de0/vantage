package reasoning

import (
	"math"
	"strings"
)

const (
	// DepthFactor scales logarithmic depth penalty so deeper paths are discouraged
	// without over-penalizing meaningful multi-hop chains.
	DepthFactor = 0.35
	// ConfidenceFactor weights average hypothesis confidence to reward paths backed
	// by stronger evidence while keeping confidence subordinate to impact.
	ConfidenceFactor = 0.75
	// SynergyBonusValue rewards paths that coordinate across three or more phases,
	// capturing cross-phase attack-chain compounding effects.
	SynergyBonusValue = 0.6
)

// TechniqueScoreWeights configures weighted multi-factor scoring.
type TechniqueScoreWeights struct {
	ImpactWeight  float64
	RiskWeight    float64
	StealthWeight float64
}

// DefaultTechniqueScoreWeights returns conservative defaults.
func DefaultTechniqueScoreWeights() TechniqueScoreWeights {
	return TechniqueScoreWeights{
		ImpactWeight:  0.5,
		RiskWeight:    0.2,
		StealthWeight: 0.3,
	}
}

// ScoreTechnique computes a single normalized score from factors.
func ScoreTechnique(effect TechniqueEffect, weights TechniqueScoreWeights) float64 {
	if weights == (TechniqueScoreWeights{}) {
		weights = DefaultTechniqueScoreWeights()
	}
	return (effect.Impact * weights.ImpactWeight) +
		((1 - effect.Risk) * weights.RiskWeight) +
		(effect.Stealth * weights.StealthWeight)
}

func scorePath(steps []Hypothesis, classes []ActionClass, objective NodeType, _ AttackPathConfig) AttackPath {
	impactSum := 0.0
	risk := 0.0
	riskPenalty := 0.0
	totalConfidence := 0.0
	phaseSet := make(map[OperationPhase]struct{})

	for i, ac := range classes {
		impactSum += ac.ImpactWeight
		risk += ac.RiskWeight
		riskPenalty += ac.RiskWeight * ac.RiskWeight
		phaseSet[ac.Phase] = struct{}{}
		if i < len(steps) {
			totalConfidence += steps[i].Confidence
		}
	}

	avgConfidence := 0.0
	if len(steps) > 0 {
		avgConfidence = totalConfidence / float64(len(steps))
	}

	objectiveMultiplier := 1.0
	if len(classes) > 0 && hasHighImpactNode(classes[len(classes)-1].ProducesNodes) {
		objectiveMultiplier = 1.4
	}

	synergyBonus := 0.0
	if len(phaseSet) >= 3 {
		synergyBonus = SynergyBonusValue
	}

	depthPenalty := math.Log(float64(len(steps))+1) * DepthFactor
	score := (impactSum * objectiveMultiplier) +
		(avgConfidence * ConfidenceFactor) +
		synergyBonus -
		riskPenalty -
		depthPenalty

	return AttackPath{Steps: steps, Score: score, Risk: risk, Objective: objective, Valid: true}
}

func hasHighImpactNode(nodes []NodeType) bool {
	for _, nodeType := range nodes {
		normalized := strings.ToLower(string(nodeType))
		switch normalized {
		case "dataexposure", "data_exposure", "externalexecution", "external_execution":
			return true
		}
	}
	return false
}

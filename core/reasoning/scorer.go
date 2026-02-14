package reasoning

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

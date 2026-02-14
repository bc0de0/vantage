package reasoning

// ConfidenceModel computes bounded confidence scores from evidence metadata.
type ConfidenceModel struct{}

func NewConfidenceModel() *ConfidenceModel {
	return &ConfidenceModel{}
}

func (m *ConfidenceModel) Score(sourceReliability float64, inferenceDepth int, evidenceCount int) float64 {
	if sourceReliability < 0 {
		sourceReliability = 0
	}
	if sourceReliability > 1 {
		sourceReliability = 1
	}

	depthPenalty := 1.0 / float64(inferenceDepth+1)
	evidenceBoost := 0.2 * float64(evidenceCount)
	if evidenceBoost > 0.6 {
		evidenceBoost = 0.6
	}

	score := (sourceReliability * 0.7) + (depthPenalty * 0.2) + evidenceBoost
	if score > 1 {
		return 1
	}
	return score
}

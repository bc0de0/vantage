package reasoning

import (
	"fmt"
	"strings"
)

// Hypothesis is an inferred operational proposition.
type Hypothesis struct {
	ID                string
	ActionClassID     string
	Statement         string
	SupportingNodeIDs []string
	DerivedFrom       []string
	Confidence        float64
}

// GenerateHypotheses derives hypotheses from current graph evidence.
// GenerateHypotheses derives baseline deterministic hypotheses from current graph evidence.
func GenerateHypotheses(graph *Graph) []Hypothesis {
	if graph == nil {
		return nil
	}
	evidenceNodes := graph.NodesByType(NodeTypeEvidence)
	out := make([]Hypothesis, 0, len(evidenceNodes))
	for idx, n := range evidenceNodes {
		statement := fmt.Sprintf("evidence from %s may enable follow-on actions", n.Label)
		confidence := 0.5
		if strings.EqualFold(n.Metadata["success"], "true") {
			confidence = 0.8
		}
		out = append(out, Hypothesis{
			ID:                fmt.Sprintf("hyp-%d-%s", idx+1, n.ID),
			Statement:         statement,
			SupportingNodeIDs: []string{n.ID},
			DerivedFrom:       []string{n.ID},
			Confidence:        confidence,
		})
	}
	return out
}

package reasoning

// SeedScenario describes synthetic environment profile.
type SeedScenario string

const (
	SeedScenarioMinimal SeedScenario = "minimal"
	SeedScenarioRich    SeedScenario = "rich"
)

// SeedSyntheticEnvironment generates deterministic graph state for simulation.
func SeedSyntheticEnvironment(g *Graph, scenario SeedScenario) {
	if g == nil {
		return
	}
	seed := []Node{{ID: "env-public-web", Type: NodeTypeEvidence, Label: "public web app exposure"}}
	if scenario == SeedScenarioRich {
		seed = append(seed,
			Node{ID: "env-hybrid-cloud", Type: NodeTypeEvidence, Label: "hybrid cloud infrastructure"},
			Node{ID: "env-segment-a", Type: NodeTypeHypothesis, Label: "internal segmented network"},
			Node{ID: "env-cred-reuse", Type: NodeTypeEvidence, Label: "credential reuse across services"},
			Node{ID: "env-priv-boundary", Type: NodeTypeHypothesis, Label: "misconfigured privilege boundary"},
		)
	}
	for i := range seed {
		n := seed[i]
		g.UpsertNode(&n)
	}
	if scenario == SeedScenarioRich {
		_ = g.AddEdge(&Edge{From: "env-public-web", To: "env-segment-a", Type: EdgeTypeSupports, Weight: 1})
		_ = g.AddEdge(&Edge{From: "env-cred-reuse", To: "env-priv-boundary", Type: EdgeTypeSupports, Weight: 1})
	}
}

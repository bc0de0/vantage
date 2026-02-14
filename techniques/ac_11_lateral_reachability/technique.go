package ac_11_lateral_reachability

import (
	"context"

	"vantage/techniques/model"
)

type evalFn func(*model.Graph) bool

type profileTechnique struct {
	id      string
	name    string
	summary string
	risk    float64
	impact  float64
	eval    evalFn
	classID string
}

func (t profileTechnique) ID() string                       { return t.id }
func (t profileTechnique) Name() string                     { return t.name }
func (t profileTechnique) ActionClassID() string            { return t.classID }
func (t profileTechnique) Evaluate(graph *model.Graph) bool { return t.eval(graph) }
func (t profileTechnique) Execute(_ context.Context, _ *model.Graph) (model.Evidence, error) {
	return model.Evidence{TechniqueID: t.ID(), Summary: t.summary, Success: true}, nil
}
func (t profileTechnique) RiskModifier() float64   { return t.risk }
func (t profileTechnique) ImpactModifier() float64 { return t.impact }

var (
	evalMinimal    = func(g *model.Graph) bool { return g != nil && g.EvidenceNodes == 0 }
	evalObserved   = func(g *model.Graph) bool { return g != nil && g.EvidenceNodes >= 1 }
	evalPivot      = func(g *model.Graph) bool { return g != nil && g.EvidenceNodes >= 1 && g.HypothesisNodes >= 1 }
	evalGraphPivot = func(g *model.Graph) bool {
		return g != nil && g.TechniqueNodes >= 1 && (g.HasSupportsEdge || g.HasEnablesEdge)
	}
	evalRare = func(g *model.Graph) bool {
		return g != nil && g.EvidenceNodes >= 2 && g.HypothesisNodes >= 2 && g.TechniqueNodes >= 1 && g.HasSupportsEdge && g.HasEnablesEdge
	}
)

// LateralReachabilityAnalyzer favors passive or low-touch checks to reduce operational risk while building baselines.
// Risk profile: low confidence drift and low detection risk due to minimal interaction.
// Confidence rationale: high when little graph state exists because observations are easy to validate.
// Expected behavior: seeds follow-on discovery classes with foundational evidence.
type LateralReachabilityAnalyzer struct{}

func (t LateralReachabilityAnalyzer) impl() profileTechnique {
	return profileTechnique{id: "AC11LateralReachabilityAnalyzer", name: "LateralReachabilityAnalyzer", classID: "AC-11", summary: "assess lateral movement reachability from established access", risk: 0.18, impact: 0.28, eval: evalObserved}
}
func (t LateralReachabilityAnalyzer) ID() string                   { return t.impl().ID() }
func (t LateralReachabilityAnalyzer) Name() string                 { return t.impl().Name() }
func (t LateralReachabilityAnalyzer) ActionClassID() string        { return t.impl().ActionClassID() }
func (t LateralReachabilityAnalyzer) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t LateralReachabilityAnalyzer) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t LateralReachabilityAnalyzer) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t LateralReachabilityAnalyzer) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// NeighborHostCensus captures high-confidence low-impact context to enrich precision targeting later.
// Risk profile: low risk and low footprint, suitable for steady-state reconnaissance.
// Confidence rationale: requires at least one evidence node to avoid blind guesses.
// Expected behavior: unlocks service and protocol follow-ups with improved context quality.
type NeighborHostCensus struct{}

func (t NeighborHostCensus) impl() profileTechnique {
	return profileTechnique{id: "AC11NeighborHostCensus", name: "NeighborHostCensus", classID: "AC-11", summary: "build structured low-noise context for later chaining", risk: 0.24, impact: 0.34, eval: evalObserved}
}
func (t NeighborHostCensus) ID() string                   { return t.impl().ID() }
func (t NeighborHostCensus) Name() string                 { return t.impl().Name() }
func (t NeighborHostCensus) ActionClassID() string        { return t.impl().ActionClassID() }
func (t NeighborHostCensus) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t NeighborHostCensus) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t NeighborHostCensus) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t NeighborHostCensus) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// SegmentRouteValidation is pivot-heavy and links multiple graph hints to expose chained opportunities.
// Risk profile: medium risk due to directional probing that can trigger controls.
// Confidence rationale: requires evidence and hypothesis alignment to confirm a viable pivot.
// Expected behavior: unlocks credential, access, or lateral classes by revealing adjacency.
type SegmentRouteValidation struct{}

func (t SegmentRouteValidation) impl() profileTechnique {
	return profileTechnique{id: "AC11SegmentRouteValidation", name: "SegmentRouteValidation", classID: "AC-11", summary: "connect mid-stage observations into pivotable pathways", risk: 0.52, impact: 0.58, eval: evalPivot}
}
func (t SegmentRouteValidation) ID() string                   { return t.impl().ID() }
func (t SegmentRouteValidation) Name() string                 { return t.impl().Name() }
func (t SegmentRouteValidation) ActionClassID() string        { return t.impl().ActionClassID() }
func (t SegmentRouteValidation) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t SegmentRouteValidation) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t SegmentRouteValidation) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t SegmentRouteValidation) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// CredentialRelayPivot is a second pivot behavior focused on graph-link validation before escalation.
// Risk profile: medium-to-high because it exercises cross-node relationships.
// Confidence rationale: needs technique evidence plus supporting or enabling edges.
// Expected behavior: produces evidence that makes execution or impact classes evaluable.
type CredentialRelayPivot struct{}

func (t CredentialRelayPivot) impl() profileTechnique {
	return profileTechnique{id: "AC11CredentialRelayPivot", name: "CredentialRelayPivot", classID: "AC-11", summary: "stress pivot assumptions across linked graph paths", risk: 0.62, impact: 0.72, eval: evalGraphPivot}
}
func (t CredentialRelayPivot) ID() string                   { return t.impl().ID() }
func (t CredentialRelayPivot) Name() string                 { return t.impl().Name() }
func (t CredentialRelayPivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t CredentialRelayPivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t CredentialRelayPivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t CredentialRelayPivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t CredentialRelayPivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// SharedServicePivot models low-confidence high-impact behavior intended for rare but decisive opportunities.
// Risk profile: high operational and detection risk with potentially outsized downstream impact.
// Confidence rationale: only evaluates true for rare graph combinations across nodes and edges.
// Expected behavior: when triggered, it unlocks objective-oriented classes with strong score effects.
type SharedServicePivot struct{}

func (t SharedServicePivot) impl() profileTechnique {
	return profileTechnique{id: "AC11SharedServicePivot", name: "SharedServicePivot", classID: "AC-11", summary: "attempt rare high-consequence maneuver when graph strongly supports it", risk: 0.84, impact: 0.92, eval: evalRare}
}
func (t SharedServicePivot) ID() string                   { return t.impl().ID() }
func (t SharedServicePivot) Name() string                 { return t.impl().Name() }
func (t SharedServicePivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t SharedServicePivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t SharedServicePivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t SharedServicePivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t SharedServicePivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// All returns the diversified technique set for AC-11.
func All() []model.Technique {
	return []model.Technique{
		LateralReachabilityAnalyzer{},
		NeighborHostCensus{},
		SegmentRouteValidation{},
		CredentialRelayPivot{},
		SharedServicePivot{},
	}
}

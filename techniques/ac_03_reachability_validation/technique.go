package ac_03_reachability_validation

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

// ReachabilityValidator favors passive or low-touch checks to reduce operational risk while building baselines.
// Risk profile: low confidence drift and low detection risk due to minimal interaction.
// Confidence rationale: high when little graph state exists because observations are easy to validate.
// Expected behavior: seeds follow-on discovery classes with foundational evidence.
type ReachabilityValidator struct{}

func (t ReachabilityValidator) impl() profileTechnique {
	return profileTechnique{id: "AC03ReachabilityValidator", name: "ReachabilityValidator", classID: "AC-03", summary: "validate host reachability using non-invasive checks", risk: 0.18, impact: 0.28, eval: evalObserved}
}
func (t ReachabilityValidator) ID() string                   { return t.impl().ID() }
func (t ReachabilityValidator) Name() string                 { return t.impl().Name() }
func (t ReachabilityValidator) ActionClassID() string        { return t.impl().ActionClassID() }
func (t ReachabilityValidator) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t ReachabilityValidator) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t ReachabilityValidator) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t ReachabilityValidator) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// ControlPathHeartbeat captures high-confidence low-impact context to enrich precision targeting later.
// Risk profile: low risk and low footprint, suitable for steady-state reconnaissance.
// Confidence rationale: requires at least one evidence node to avoid blind guesses.
// Expected behavior: unlocks service and protocol follow-ups with improved context quality.
type ControlPathHeartbeat struct{}

func (t ControlPathHeartbeat) impl() profileTechnique {
	return profileTechnique{id: "AC03ControlPathHeartbeat", name: "ControlPathHeartbeat", classID: "AC-03", summary: "build structured low-noise context for later chaining", risk: 0.24, impact: 0.34, eval: evalObserved}
}
func (t ControlPathHeartbeat) ID() string                   { return t.impl().ID() }
func (t ControlPathHeartbeat) Name() string                 { return t.impl().Name() }
func (t ControlPathHeartbeat) ActionClassID() string        { return t.impl().ActionClassID() }
func (t ControlPathHeartbeat) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t ControlPathHeartbeat) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t ControlPathHeartbeat) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t ControlPathHeartbeat) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// LatencyEnvelopeCheck is pivot-heavy and links multiple graph hints to expose chained opportunities.
// Risk profile: medium risk due to directional probing that can trigger controls.
// Confidence rationale: requires evidence and hypothesis alignment to confirm a viable pivot.
// Expected behavior: unlocks credential, access, or lateral classes by revealing adjacency.
type LatencyEnvelopeCheck struct{}

func (t LatencyEnvelopeCheck) impl() profileTechnique {
	return profileTechnique{id: "AC03LatencyEnvelopeCheck", name: "LatencyEnvelopeCheck", classID: "AC-03", summary: "connect mid-stage observations into pivotable pathways", risk: 0.52, impact: 0.58, eval: evalPivot}
}
func (t LatencyEnvelopeCheck) ID() string                   { return t.impl().ID() }
func (t LatencyEnvelopeCheck) Name() string                 { return t.impl().Name() }
func (t LatencyEnvelopeCheck) ActionClassID() string        { return t.impl().ActionClassID() }
func (t LatencyEnvelopeCheck) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t LatencyEnvelopeCheck) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t LatencyEnvelopeCheck) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t LatencyEnvelopeCheck) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// TransitTrustPivot is a second pivot behavior focused on graph-link validation before escalation.
// Risk profile: medium-to-high because it exercises cross-node relationships.
// Confidence rationale: needs technique evidence plus supporting or enabling edges.
// Expected behavior: produces evidence that makes execution or impact classes evaluable.
type TransitTrustPivot struct{}

func (t TransitTrustPivot) impl() profileTechnique {
	return profileTechnique{id: "AC03TransitTrustPivot", name: "TransitTrustPivot", classID: "AC-03", summary: "stress pivot assumptions across linked graph paths", risk: 0.62, impact: 0.72, eval: evalGraphPivot}
}
func (t TransitTrustPivot) ID() string                   { return t.impl().ID() }
func (t TransitTrustPivot) Name() string                 { return t.impl().Name() }
func (t TransitTrustPivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t TransitTrustPivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t TransitTrustPivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t TransitTrustPivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t TransitTrustPivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// DualStackRoutePivot models low-confidence high-impact behavior intended for rare but decisive opportunities.
// Risk profile: high operational and detection risk with potentially outsized downstream impact.
// Confidence rationale: only evaluates true for rare graph combinations across nodes and edges.
// Expected behavior: when triggered, it unlocks objective-oriented classes with strong score effects.
type DualStackRoutePivot struct{}

func (t DualStackRoutePivot) impl() profileTechnique {
	return profileTechnique{id: "AC03DualStackRoutePivot", name: "DualStackRoutePivot", classID: "AC-03", summary: "attempt rare high-consequence maneuver when graph strongly supports it", risk: 0.84, impact: 0.92, eval: evalRare}
}
func (t DualStackRoutePivot) ID() string                   { return t.impl().ID() }
func (t DualStackRoutePivot) Name() string                 { return t.impl().Name() }
func (t DualStackRoutePivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t DualStackRoutePivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t DualStackRoutePivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t DualStackRoutePivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t DualStackRoutePivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// All returns the diversified technique set for AC-03.
func All() []model.Technique {
	return []model.Technique{
		ReachabilityValidator{},
		ControlPathHeartbeat{},
		LatencyEnvelopeCheck{},
		TransitTrustPivot{},
		DualStackRoutePivot{},
	}
}

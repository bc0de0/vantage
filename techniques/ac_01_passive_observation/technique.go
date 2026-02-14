package ac_01_passive_observation

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

// PassiveDNSCollection favors passive or low-touch checks to reduce operational risk while building baselines.
// Risk profile: low confidence drift and low detection risk due to minimal interaction.
// Confidence rationale: high when little graph state exists because observations are easy to validate.
// Expected behavior: seeds follow-on discovery classes with foundational evidence.
type PassiveDNSCollection struct{}

func (t PassiveDNSCollection) impl() profileTechnique {
	return profileTechnique{id: "AC01PassiveDNSCollection", name: "PassiveDNSCollection", classID: "AC-01", summary: "collect passive DNS and certificate transparency observations", risk: 0.18, impact: 0.28, eval: evalMinimal}
}
func (t PassiveDNSCollection) ID() string                   { return t.impl().ID() }
func (t PassiveDNSCollection) Name() string                 { return t.impl().Name() }
func (t PassiveDNSCollection) ActionClassID() string        { return t.impl().ActionClassID() }
func (t PassiveDNSCollection) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t PassiveDNSCollection) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t PassiveDNSCollection) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t PassiveDNSCollection) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// PassiveSourceCorrelator captures high-confidence low-impact context to enrich precision targeting later.
// Risk profile: low risk and low footprint, suitable for steady-state reconnaissance.
// Confidence rationale: requires at least one evidence node to avoid blind guesses.
// Expected behavior: unlocks service and protocol follow-ups with improved context quality.
type PassiveSourceCorrelator struct{}

func (t PassiveSourceCorrelator) impl() profileTechnique {
	return profileTechnique{id: "AC01PassiveSourceCorrelator", name: "PassiveSourceCorrelator", classID: "AC-01", summary: "build structured low-noise context for later chaining", risk: 0.24, impact: 0.34, eval: evalObserved}
}
func (t PassiveSourceCorrelator) ID() string                   { return t.impl().ID() }
func (t PassiveSourceCorrelator) Name() string                 { return t.impl().Name() }
func (t PassiveSourceCorrelator) ActionClassID() string        { return t.impl().ActionClassID() }
func (t PassiveSourceCorrelator) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t PassiveSourceCorrelator) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t PassiveSourceCorrelator) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t PassiveSourceCorrelator) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// OrgExposureCatalog is pivot-heavy and links multiple graph hints to expose chained opportunities.
// Risk profile: medium risk due to directional probing that can trigger controls.
// Confidence rationale: requires evidence and hypothesis alignment to confirm a viable pivot.
// Expected behavior: unlocks credential, access, or lateral classes by revealing adjacency.
type OrgExposureCatalog struct{}

func (t OrgExposureCatalog) impl() profileTechnique {
	return profileTechnique{id: "AC01OrgExposureCatalog", name: "OrgExposureCatalog", classID: "AC-01", summary: "connect mid-stage observations into pivotable pathways", risk: 0.52, impact: 0.58, eval: evalPivot}
}
func (t OrgExposureCatalog) ID() string                   { return t.impl().ID() }
func (t OrgExposureCatalog) Name() string                 { return t.impl().Name() }
func (t OrgExposureCatalog) ActionClassID() string        { return t.impl().ActionClassID() }
func (t OrgExposureCatalog) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t OrgExposureCatalog) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t OrgExposureCatalog) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t OrgExposureCatalog) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// TrustChainPivot is a second pivot behavior focused on graph-link validation before escalation.
// Risk profile: medium-to-high because it exercises cross-node relationships.
// Confidence rationale: needs technique evidence plus supporting or enabling edges.
// Expected behavior: produces evidence that makes execution or impact classes evaluable.
type TrustChainPivot struct{}

func (t TrustChainPivot) impl() profileTechnique {
	return profileTechnique{id: "AC01TrustChainPivot", name: "TrustChainPivot", classID: "AC-01", summary: "stress pivot assumptions across linked graph paths", risk: 0.62, impact: 0.72, eval: evalGraphPivot}
}
func (t TrustChainPivot) ID() string                   { return t.impl().ID() }
func (t TrustChainPivot) Name() string                 { return t.impl().Name() }
func (t TrustChainPivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t TrustChainPivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t TrustChainPivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t TrustChainPivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t TrustChainPivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// ThirdPartySignalPivot models low-confidence high-impact behavior intended for rare but decisive opportunities.
// Risk profile: high operational and detection risk with potentially outsized downstream impact.
// Confidence rationale: only evaluates true for rare graph combinations across nodes and edges.
// Expected behavior: when triggered, it unlocks objective-oriented classes with strong score effects.
type ThirdPartySignalPivot struct{}

func (t ThirdPartySignalPivot) impl() profileTechnique {
	return profileTechnique{id: "AC01ThirdPartySignalPivot", name: "ThirdPartySignalPivot", classID: "AC-01", summary: "attempt rare high-consequence maneuver when graph strongly supports it", risk: 0.84, impact: 0.92, eval: evalRare}
}
func (t ThirdPartySignalPivot) ID() string                   { return t.impl().ID() }
func (t ThirdPartySignalPivot) Name() string                 { return t.impl().Name() }
func (t ThirdPartySignalPivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t ThirdPartySignalPivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t ThirdPartySignalPivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t ThirdPartySignalPivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t ThirdPartySignalPivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// All returns the diversified technique set for AC-01.
func All() []model.Technique {
	return []model.Technique{
		PassiveDNSCollection{},
		PassiveSourceCorrelator{},
		OrgExposureCatalog{},
		TrustChainPivot{},
		ThirdPartySignalPivot{},
	}
}

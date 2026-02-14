package ac_04_service_identification

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

// ServiceIdentifier favors passive or low-touch checks to reduce operational risk while building baselines.
// Risk profile: low confidence drift and low detection risk due to minimal interaction.
// Confidence rationale: high when little graph state exists because observations are easy to validate.
// Expected behavior: seeds follow-on discovery classes with foundational evidence.
type ServiceIdentifier struct{}

func (t ServiceIdentifier) impl() profileTechnique {
	return profileTechnique{id: "AC04ServiceIdentifier", name: "ServiceIdentifier", classID: "AC-04", summary: "identify exposed network services", risk: 0.18, impact: 0.28, eval: evalObserved}
}
func (t ServiceIdentifier) ID() string                   { return t.impl().ID() }
func (t ServiceIdentifier) Name() string                 { return t.impl().Name() }
func (t ServiceIdentifier) ActionClassID() string        { return t.impl().ActionClassID() }
func (t ServiceIdentifier) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t ServiceIdentifier) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t ServiceIdentifier) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t ServiceIdentifier) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// BannerRoleMapper captures high-confidence low-impact context to enrich precision targeting later.
// Risk profile: low risk and low footprint, suitable for steady-state reconnaissance.
// Confidence rationale: requires at least one evidence node to avoid blind guesses.
// Expected behavior: unlocks service and protocol follow-ups with improved context quality.
type BannerRoleMapper struct{}

func (t BannerRoleMapper) impl() profileTechnique {
	return profileTechnique{id: "AC04BannerRoleMapper", name: "BannerRoleMapper", classID: "AC-04", summary: "build structured low-noise context for later chaining", risk: 0.24, impact: 0.34, eval: evalObserved}
}
func (t BannerRoleMapper) ID() string                   { return t.impl().ID() }
func (t BannerRoleMapper) Name() string                 { return t.impl().Name() }
func (t BannerRoleMapper) ActionClassID() string        { return t.impl().ActionClassID() }
func (t BannerRoleMapper) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t BannerRoleMapper) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t BannerRoleMapper) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t BannerRoleMapper) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// PortRoleTriager is pivot-heavy and links multiple graph hints to expose chained opportunities.
// Risk profile: medium risk due to directional probing that can trigger controls.
// Confidence rationale: requires evidence and hypothesis alignment to confirm a viable pivot.
// Expected behavior: unlocks credential, access, or lateral classes by revealing adjacency.
type PortRoleTriager struct{}

func (t PortRoleTriager) impl() profileTechnique {
	return profileTechnique{id: "AC04PortRoleTriager", name: "PortRoleTriager", classID: "AC-04", summary: "connect mid-stage observations into pivotable pathways", risk: 0.52, impact: 0.58, eval: evalPivot}
}
func (t PortRoleTriager) ID() string                   { return t.impl().ID() }
func (t PortRoleTriager) Name() string                 { return t.impl().Name() }
func (t PortRoleTriager) ActionClassID() string        { return t.impl().ActionClassID() }
func (t PortRoleTriager) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t PortRoleTriager) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t PortRoleTriager) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t PortRoleTriager) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// ServiceDependencyPivot is a second pivot behavior focused on graph-link validation before escalation.
// Risk profile: medium-to-high because it exercises cross-node relationships.
// Confidence rationale: needs technique evidence plus supporting or enabling edges.
// Expected behavior: produces evidence that makes execution or impact classes evaluable.
type ServiceDependencyPivot struct{}

func (t ServiceDependencyPivot) impl() profileTechnique {
	return profileTechnique{id: "AC04ServiceDependencyPivot", name: "ServiceDependencyPivot", classID: "AC-04", summary: "stress pivot assumptions across linked graph paths", risk: 0.62, impact: 0.72, eval: evalGraphPivot}
}
func (t ServiceDependencyPivot) ID() string                   { return t.impl().ID() }
func (t ServiceDependencyPivot) Name() string                 { return t.impl().Name() }
func (t ServiceDependencyPivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t ServiceDependencyPivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t ServiceDependencyPivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t ServiceDependencyPivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t ServiceDependencyPivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// CrossTierBindingPivot models low-confidence high-impact behavior intended for rare but decisive opportunities.
// Risk profile: high operational and detection risk with potentially outsized downstream impact.
// Confidence rationale: only evaluates true for rare graph combinations across nodes and edges.
// Expected behavior: when triggered, it unlocks objective-oriented classes with strong score effects.
type CrossTierBindingPivot struct{}

func (t CrossTierBindingPivot) impl() profileTechnique {
	return profileTechnique{id: "AC04CrossTierBindingPivot", name: "CrossTierBindingPivot", classID: "AC-04", summary: "attempt rare high-consequence maneuver when graph strongly supports it", risk: 0.84, impact: 0.92, eval: evalRare}
}
func (t CrossTierBindingPivot) ID() string                   { return t.impl().ID() }
func (t CrossTierBindingPivot) Name() string                 { return t.impl().Name() }
func (t CrossTierBindingPivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t CrossTierBindingPivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t CrossTierBindingPivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t CrossTierBindingPivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t CrossTierBindingPivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// All returns the diversified technique set for AC-04.
func All() []model.Technique {
	return []model.Technique{
		ServiceIdentifier{},
		BannerRoleMapper{},
		PortRoleTriager{},
		ServiceDependencyPivot{},
		CrossTierBindingPivot{},
	}
}

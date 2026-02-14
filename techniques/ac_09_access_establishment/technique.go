package ac_09_access_establishment

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

// AccessEstablisher favors passive or low-touch checks to reduce operational risk while building baselines.
// Risk profile: low confidence drift and low detection risk due to minimal interaction.
// Confidence rationale: high when little graph state exists because observations are easy to validate.
// Expected behavior: seeds follow-on discovery classes with foundational evidence.
type AccessEstablisher struct{}

func (t AccessEstablisher) impl() profileTechnique {
	return profileTechnique{id: "AC09AccessEstablisher", name: "AccessEstablisher", classID: "AC-09", summary: "establish authenticated access with validated material", risk: 0.18, impact: 0.28, eval: evalObserved}
}
func (t AccessEstablisher) ID() string                   { return t.impl().ID() }
func (t AccessEstablisher) Name() string                 { return t.impl().Name() }
func (t AccessEstablisher) ActionClassID() string        { return t.impl().ActionClassID() }
func (t AccessEstablisher) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t AccessEstablisher) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t AccessEstablisher) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t AccessEstablisher) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// LeastPrivilegeSessionBootstrap captures high-confidence low-impact context to enrich precision targeting later.
// Risk profile: low risk and low footprint, suitable for steady-state reconnaissance.
// Confidence rationale: requires at least one evidence node to avoid blind guesses.
// Expected behavior: unlocks service and protocol follow-ups with improved context quality.
type LeastPrivilegeSessionBootstrap struct{}

func (t LeastPrivilegeSessionBootstrap) impl() profileTechnique {
	return profileTechnique{id: "AC09LeastPrivilegeSessionBootstrap", name: "LeastPrivilegeSessionBootstrap", classID: "AC-09", summary: "build structured low-noise context for later chaining", risk: 0.24, impact: 0.34, eval: evalObserved}
}
func (t LeastPrivilegeSessionBootstrap) ID() string                   { return t.impl().ID() }
func (t LeastPrivilegeSessionBootstrap) Name() string                 { return t.impl().Name() }
func (t LeastPrivilegeSessionBootstrap) ActionClassID() string        { return t.impl().ActionClassID() }
func (t LeastPrivilegeSessionBootstrap) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t LeastPrivilegeSessionBootstrap) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t LeastPrivilegeSessionBootstrap) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t LeastPrivilegeSessionBootstrap) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// EphemeralAccessTrial is pivot-heavy and links multiple graph hints to expose chained opportunities.
// Risk profile: medium risk due to directional probing that can trigger controls.
// Confidence rationale: requires evidence and hypothesis alignment to confirm a viable pivot.
// Expected behavior: unlocks credential, access, or lateral classes by revealing adjacency.
type EphemeralAccessTrial struct{}

func (t EphemeralAccessTrial) impl() profileTechnique {
	return profileTechnique{id: "AC09EphemeralAccessTrial", name: "EphemeralAccessTrial", classID: "AC-09", summary: "connect mid-stage observations into pivotable pathways", risk: 0.52, impact: 0.58, eval: evalPivot}
}
func (t EphemeralAccessTrial) ID() string                   { return t.impl().ID() }
func (t EphemeralAccessTrial) Name() string                 { return t.impl().Name() }
func (t EphemeralAccessTrial) ActionClassID() string        { return t.impl().ActionClassID() }
func (t EphemeralAccessTrial) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t EphemeralAccessTrial) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t EphemeralAccessTrial) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t EphemeralAccessTrial) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// SessionReusePivot is a second pivot behavior focused on graph-link validation before escalation.
// Risk profile: medium-to-high because it exercises cross-node relationships.
// Confidence rationale: needs technique evidence plus supporting or enabling edges.
// Expected behavior: produces evidence that makes execution or impact classes evaluable.
type SessionReusePivot struct{}

func (t SessionReusePivot) impl() profileTechnique {
	return profileTechnique{id: "AC09SessionReusePivot", name: "SessionReusePivot", classID: "AC-09", summary: "stress pivot assumptions across linked graph paths", risk: 0.62, impact: 0.72, eval: evalGraphPivot}
}
func (t SessionReusePivot) ID() string                   { return t.impl().ID() }
func (t SessionReusePivot) Name() string                 { return t.impl().Name() }
func (t SessionReusePivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t SessionReusePivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t SessionReusePivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t SessionReusePivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t SessionReusePivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// TrustPathPivot models low-confidence high-impact behavior intended for rare but decisive opportunities.
// Risk profile: high operational and detection risk with potentially outsized downstream impact.
// Confidence rationale: only evaluates true for rare graph combinations across nodes and edges.
// Expected behavior: when triggered, it unlocks objective-oriented classes with strong score effects.
type TrustPathPivot struct{}

func (t TrustPathPivot) impl() profileTechnique {
	return profileTechnique{id: "AC09TrustPathPivot", name: "TrustPathPivot", classID: "AC-09", summary: "attempt rare high-consequence maneuver when graph strongly supports it", risk: 0.84, impact: 0.92, eval: evalRare}
}
func (t TrustPathPivot) ID() string                   { return t.impl().ID() }
func (t TrustPathPivot) Name() string                 { return t.impl().Name() }
func (t TrustPathPivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t TrustPathPivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t TrustPathPivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t TrustPathPivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t TrustPathPivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// All returns the diversified technique set for AC-09.
func All() []model.Technique {
	return []model.Technique{
		AccessEstablisher{},
		LeastPrivilegeSessionBootstrap{},
		EphemeralAccessTrial{},
		SessionReusePivot{},
		TrustPathPivot{},
	}
}

package ac_07_auth_surface_analysis

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

// AuthSurfaceAnalyzer favors passive or low-touch checks to reduce operational risk while building baselines.
// Risk profile: low confidence drift and low detection risk due to minimal interaction.
// Confidence rationale: high when little graph state exists because observations are easy to validate.
// Expected behavior: seeds follow-on discovery classes with foundational evidence.
type AuthSurfaceAnalyzer struct{}

func (t AuthSurfaceAnalyzer) impl() profileTechnique {
	return profileTechnique{id: "AC07AuthSurfaceAnalyzer", name: "AuthSurfaceAnalyzer", classID: "AC-07", summary: "analyze authentication surfaces and login paths", risk: 0.18, impact: 0.28, eval: evalObserved}
}
func (t AuthSurfaceAnalyzer) ID() string                   { return t.impl().ID() }
func (t AuthSurfaceAnalyzer) Name() string                 { return t.impl().Name() }
func (t AuthSurfaceAnalyzer) ActionClassID() string        { return t.impl().ActionClassID() }
func (t AuthSurfaceAnalyzer) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t AuthSurfaceAnalyzer) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t AuthSurfaceAnalyzer) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t AuthSurfaceAnalyzer) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// LoginFlowClassifier captures high-confidence low-impact context to enrich precision targeting later.
// Risk profile: low risk and low footprint, suitable for steady-state reconnaissance.
// Confidence rationale: requires at least one evidence node to avoid blind guesses.
// Expected behavior: unlocks service and protocol follow-ups with improved context quality.
type LoginFlowClassifier struct{}

func (t LoginFlowClassifier) impl() profileTechnique {
	return profileTechnique{id: "AC07LoginFlowClassifier", name: "LoginFlowClassifier", classID: "AC-07", summary: "build structured low-noise context for later chaining", risk: 0.24, impact: 0.34, eval: evalObserved}
}
func (t LoginFlowClassifier) ID() string                   { return t.impl().ID() }
func (t LoginFlowClassifier) Name() string                 { return t.impl().Name() }
func (t LoginFlowClassifier) ActionClassID() string        { return t.impl().ActionClassID() }
func (t LoginFlowClassifier) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t LoginFlowClassifier) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t LoginFlowClassifier) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t LoginFlowClassifier) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// MFAChannelInventory is pivot-heavy and links multiple graph hints to expose chained opportunities.
// Risk profile: medium risk due to directional probing that can trigger controls.
// Confidence rationale: requires evidence and hypothesis alignment to confirm a viable pivot.
// Expected behavior: unlocks credential, access, or lateral classes by revealing adjacency.
type MFAChannelInventory struct{}

func (t MFAChannelInventory) impl() profileTechnique {
	return profileTechnique{id: "AC07MFAChannelInventory", name: "MFAChannelInventory", classID: "AC-07", summary: "connect mid-stage observations into pivotable pathways", risk: 0.52, impact: 0.58, eval: evalPivot}
}
func (t MFAChannelInventory) ID() string                   { return t.impl().ID() }
func (t MFAChannelInventory) Name() string                 { return t.impl().Name() }
func (t MFAChannelInventory) ActionClassID() string        { return t.impl().ActionClassID() }
func (t MFAChannelInventory) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t MFAChannelInventory) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t MFAChannelInventory) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t MFAChannelInventory) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// SessionBoundaryPivot is a second pivot behavior focused on graph-link validation before escalation.
// Risk profile: medium-to-high because it exercises cross-node relationships.
// Confidence rationale: needs technique evidence plus supporting or enabling edges.
// Expected behavior: produces evidence that makes execution or impact classes evaluable.
type SessionBoundaryPivot struct{}

func (t SessionBoundaryPivot) impl() profileTechnique {
	return profileTechnique{id: "AC07SessionBoundaryPivot", name: "SessionBoundaryPivot", classID: "AC-07", summary: "stress pivot assumptions across linked graph paths", risk: 0.62, impact: 0.72, eval: evalGraphPivot}
}
func (t SessionBoundaryPivot) ID() string                   { return t.impl().ID() }
func (t SessionBoundaryPivot) Name() string                 { return t.impl().Name() }
func (t SessionBoundaryPivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t SessionBoundaryPivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t SessionBoundaryPivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t SessionBoundaryPivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t SessionBoundaryPivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// IdentityFederationPivot models low-confidence high-impact behavior intended for rare but decisive opportunities.
// Risk profile: high operational and detection risk with potentially outsized downstream impact.
// Confidence rationale: only evaluates true for rare graph combinations across nodes and edges.
// Expected behavior: when triggered, it unlocks objective-oriented classes with strong score effects.
type IdentityFederationPivot struct{}

func (t IdentityFederationPivot) impl() profileTechnique {
	return profileTechnique{id: "AC07IdentityFederationPivot", name: "IdentityFederationPivot", classID: "AC-07", summary: "attempt rare high-consequence maneuver when graph strongly supports it", risk: 0.84, impact: 0.92, eval: evalRare}
}
func (t IdentityFederationPivot) ID() string                   { return t.impl().ID() }
func (t IdentityFederationPivot) Name() string                 { return t.impl().Name() }
func (t IdentityFederationPivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t IdentityFederationPivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t IdentityFederationPivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t IdentityFederationPivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t IdentityFederationPivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// All returns the diversified technique set for AC-07.
func All() []model.Technique {
	return []model.Technique{
		AuthSurfaceAnalyzer{},
		LoginFlowClassifier{},
		MFAChannelInventory{},
		SessionBoundaryPivot{},
		IdentityFederationPivot{},
	}
}

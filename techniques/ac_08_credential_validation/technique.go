package ac_08_credential_validation

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

// CredentialValidator favors passive or low-touch checks to reduce operational risk while building baselines.
// Risk profile: low confidence drift and low detection risk due to minimal interaction.
// Confidence rationale: high when little graph state exists because observations are easy to validate.
// Expected behavior: seeds follow-on discovery classes with foundational evidence.
type CredentialValidator struct{}

func (t CredentialValidator) impl() profileTechnique {
	return profileTechnique{id: "AC08CredentialValidator", name: "CredentialValidator", classID: "AC-08", summary: "validate credential material against target auth interfaces", risk: 0.18, impact: 0.28, eval: evalObserved}
}
func (t CredentialValidator) ID() string                   { return t.impl().ID() }
func (t CredentialValidator) Name() string                 { return t.impl().Name() }
func (t CredentialValidator) ActionClassID() string        { return t.impl().ActionClassID() }
func (t CredentialValidator) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t CredentialValidator) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t CredentialValidator) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t CredentialValidator) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// CredentialFormatLint captures high-confidence low-impact context to enrich precision targeting later.
// Risk profile: low risk and low footprint, suitable for steady-state reconnaissance.
// Confidence rationale: requires at least one evidence node to avoid blind guesses.
// Expected behavior: unlocks service and protocol follow-ups with improved context quality.
type CredentialFormatLint struct{}

func (t CredentialFormatLint) impl() profileTechnique {
	return profileTechnique{id: "AC08CredentialFormatLint", name: "CredentialFormatLint", classID: "AC-08", summary: "build structured low-noise context for later chaining", risk: 0.24, impact: 0.34, eval: evalObserved}
}
func (t CredentialFormatLint) ID() string                   { return t.impl().ID() }
func (t CredentialFormatLint) Name() string                 { return t.impl().Name() }
func (t CredentialFormatLint) ActionClassID() string        { return t.impl().ActionClassID() }
func (t CredentialFormatLint) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t CredentialFormatLint) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t CredentialFormatLint) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t CredentialFormatLint) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// LowRateCredentialCheck is pivot-heavy and links multiple graph hints to expose chained opportunities.
// Risk profile: medium risk due to directional probing that can trigger controls.
// Confidence rationale: requires evidence and hypothesis alignment to confirm a viable pivot.
// Expected behavior: unlocks credential, access, or lateral classes by revealing adjacency.
type LowRateCredentialCheck struct{}

func (t LowRateCredentialCheck) impl() profileTechnique {
	return profileTechnique{id: "AC08LowRateCredentialCheck", name: "LowRateCredentialCheck", classID: "AC-08", summary: "connect mid-stage observations into pivotable pathways", risk: 0.52, impact: 0.58, eval: evalPivot}
}
func (t LowRateCredentialCheck) ID() string                   { return t.impl().ID() }
func (t LowRateCredentialCheck) Name() string                 { return t.impl().Name() }
func (t LowRateCredentialCheck) ActionClassID() string        { return t.impl().ActionClassID() }
func (t LowRateCredentialCheck) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t LowRateCredentialCheck) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t LowRateCredentialCheck) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t LowRateCredentialCheck) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// PasswordReusePivot is a second pivot behavior focused on graph-link validation before escalation.
// Risk profile: medium-to-high because it exercises cross-node relationships.
// Confidence rationale: needs technique evidence plus supporting or enabling edges.
// Expected behavior: produces evidence that makes execution or impact classes evaluable.
type PasswordReusePivot struct{}

func (t PasswordReusePivot) impl() profileTechnique {
	return profileTechnique{id: "AC08PasswordReusePivot", name: "PasswordReusePivot", classID: "AC-08", summary: "stress pivot assumptions across linked graph paths", risk: 0.62, impact: 0.72, eval: evalGraphPivot}
}
func (t PasswordReusePivot) ID() string                   { return t.impl().ID() }
func (t PasswordReusePivot) Name() string                 { return t.impl().Name() }
func (t PasswordReusePivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t PasswordReusePivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t PasswordReusePivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t PasswordReusePivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t PasswordReusePivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// FederatedCredentialPivot models low-confidence high-impact behavior intended for rare but decisive opportunities.
// Risk profile: high operational and detection risk with potentially outsized downstream impact.
// Confidence rationale: only evaluates true for rare graph combinations across nodes and edges.
// Expected behavior: when triggered, it unlocks objective-oriented classes with strong score effects.
type FederatedCredentialPivot struct{}

func (t FederatedCredentialPivot) impl() profileTechnique {
	return profileTechnique{id: "AC08FederatedCredentialPivot", name: "FederatedCredentialPivot", classID: "AC-08", summary: "attempt rare high-consequence maneuver when graph strongly supports it", risk: 0.84, impact: 0.92, eval: evalRare}
}
func (t FederatedCredentialPivot) ID() string                   { return t.impl().ID() }
func (t FederatedCredentialPivot) Name() string                 { return t.impl().Name() }
func (t FederatedCredentialPivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t FederatedCredentialPivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t FederatedCredentialPivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t FederatedCredentialPivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t FederatedCredentialPivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// All returns the diversified technique set for AC-08.
func All() []model.Technique {
	return []model.Technique{
		CredentialValidator{},
		CredentialFormatLint{},
		LowRateCredentialCheck{},
		PasswordReusePivot{},
		FederatedCredentialPivot{},
	}
}

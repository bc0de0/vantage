package ac_14_impact_feasibility

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

// ImpactFeasibilityAssessor favors passive or low-touch checks to reduce operational risk while building baselines.
// Risk profile: low confidence drift and low detection risk due to minimal interaction.
// Confidence rationale: high when little graph state exists because observations are easy to validate.
// Expected behavior: seeds follow-on discovery classes with foundational evidence.
type ImpactFeasibilityAssessor struct{}

func (t ImpactFeasibilityAssessor) impl() profileTechnique {
	return profileTechnique{id: "AC14ImpactFeasibilityAssessor", name: "ImpactFeasibilityAssessor", classID: "AC-14", summary: "assess whether operational impact is feasible", risk: 0.18, impact: 0.28, eval: evalObserved}
}
func (t ImpactFeasibilityAssessor) ID() string                   { return t.impl().ID() }
func (t ImpactFeasibilityAssessor) Name() string                 { return t.impl().Name() }
func (t ImpactFeasibilityAssessor) ActionClassID() string        { return t.impl().ActionClassID() }
func (t ImpactFeasibilityAssessor) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t ImpactFeasibilityAssessor) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t ImpactFeasibilityAssessor) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t ImpactFeasibilityAssessor) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// ProcessFragilityReview captures high-confidence low-impact context to enrich precision targeting later.
// Risk profile: low risk and low footprint, suitable for steady-state reconnaissance.
// Confidence rationale: requires at least one evidence node to avoid blind guesses.
// Expected behavior: unlocks service and protocol follow-ups with improved context quality.
type ProcessFragilityReview struct{}

func (t ProcessFragilityReview) impl() profileTechnique {
	return profileTechnique{id: "AC14ProcessFragilityReview", name: "ProcessFragilityReview", classID: "AC-14", summary: "build structured low-noise context for later chaining", risk: 0.24, impact: 0.34, eval: evalObserved}
}
func (t ProcessFragilityReview) ID() string                   { return t.impl().ID() }
func (t ProcessFragilityReview) Name() string                 { return t.impl().Name() }
func (t ProcessFragilityReview) ActionClassID() string        { return t.impl().ActionClassID() }
func (t ProcessFragilityReview) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t ProcessFragilityReview) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t ProcessFragilityReview) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t ProcessFragilityReview) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// RecoveryWindowEstimate is pivot-heavy and links multiple graph hints to expose chained opportunities.
// Risk profile: medium risk due to directional probing that can trigger controls.
// Confidence rationale: requires evidence and hypothesis alignment to confirm a viable pivot.
// Expected behavior: unlocks credential, access, or lateral classes by revealing adjacency.
type RecoveryWindowEstimate struct{}

func (t RecoveryWindowEstimate) impl() profileTechnique {
	return profileTechnique{id: "AC14RecoveryWindowEstimate", name: "RecoveryWindowEstimate", classID: "AC-14", summary: "connect mid-stage observations into pivotable pathways", risk: 0.52, impact: 0.58, eval: evalPivot}
}
func (t RecoveryWindowEstimate) ID() string                   { return t.impl().ID() }
func (t RecoveryWindowEstimate) Name() string                 { return t.impl().Name() }
func (t RecoveryWindowEstimate) ActionClassID() string        { return t.impl().ActionClassID() }
func (t RecoveryWindowEstimate) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t RecoveryWindowEstimate) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t RecoveryWindowEstimate) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t RecoveryWindowEstimate) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// BusinessWorkflowPivot is a second pivot behavior focused on graph-link validation before escalation.
// Risk profile: medium-to-high because it exercises cross-node relationships.
// Confidence rationale: needs technique evidence plus supporting or enabling edges.
// Expected behavior: produces evidence that makes execution or impact classes evaluable.
type BusinessWorkflowPivot struct{}

func (t BusinessWorkflowPivot) impl() profileTechnique {
	return profileTechnique{id: "AC14BusinessWorkflowPivot", name: "BusinessWorkflowPivot", classID: "AC-14", summary: "stress pivot assumptions across linked graph paths", risk: 0.62, impact: 0.72, eval: evalGraphPivot}
}
func (t BusinessWorkflowPivot) ID() string                   { return t.impl().ID() }
func (t BusinessWorkflowPivot) Name() string                 { return t.impl().Name() }
func (t BusinessWorkflowPivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t BusinessWorkflowPivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t BusinessWorkflowPivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t BusinessWorkflowPivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t BusinessWorkflowPivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// DependencyCascadePivot models low-confidence high-impact behavior intended for rare but decisive opportunities.
// Risk profile: high operational and detection risk with potentially outsized downstream impact.
// Confidence rationale: only evaluates true for rare graph combinations across nodes and edges.
// Expected behavior: when triggered, it unlocks objective-oriented classes with strong score effects.
type DependencyCascadePivot struct{}

func (t DependencyCascadePivot) impl() profileTechnique {
	return profileTechnique{id: "AC14DependencyCascadePivot", name: "DependencyCascadePivot", classID: "AC-14", summary: "attempt rare high-consequence maneuver when graph strongly supports it", risk: 0.84, impact: 0.92, eval: evalRare}
}
func (t DependencyCascadePivot) ID() string                   { return t.impl().ID() }
func (t DependencyCascadePivot) Name() string                 { return t.impl().Name() }
func (t DependencyCascadePivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t DependencyCascadePivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t DependencyCascadePivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t DependencyCascadePivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t DependencyCascadePivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// All returns the diversified technique set for AC-14.
func All() []model.Technique {
	return []model.Technique{
		ImpactFeasibilityAssessor{},
		ProcessFragilityReview{},
		RecoveryWindowEstimate{},
		BusinessWorkflowPivot{},
		DependencyCascadePivot{},
	}
}

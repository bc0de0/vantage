package ac_10_privilege_assessment

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

// PrivilegeAssessor favors passive or low-touch checks to reduce operational risk while building baselines.
// Risk profile: low confidence drift and low detection risk due to minimal interaction.
// Confidence rationale: high when little graph state exists because observations are easy to validate.
// Expected behavior: seeds follow-on discovery classes with foundational evidence.
type PrivilegeAssessor struct{}

func (t PrivilegeAssessor) impl() profileTechnique {
	return profileTechnique{id: "AC10PrivilegeAssessor", name: "PrivilegeAssessor", classID: "AC-10", summary: "assess privilege level and escalation opportunities", risk: 0.18, impact: 0.28, eval: evalObserved}
}
func (t PrivilegeAssessor) ID() string                   { return t.impl().ID() }
func (t PrivilegeAssessor) Name() string                 { return t.impl().Name() }
func (t PrivilegeAssessor) ActionClassID() string        { return t.impl().ActionClassID() }
func (t PrivilegeAssessor) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t PrivilegeAssessor) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t PrivilegeAssessor) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t PrivilegeAssessor) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// RoleDriftSurvey captures high-confidence low-impact context to enrich precision targeting later.
// Risk profile: low risk and low footprint, suitable for steady-state reconnaissance.
// Confidence rationale: requires at least one evidence node to avoid blind guesses.
// Expected behavior: unlocks service and protocol follow-ups with improved context quality.
type RoleDriftSurvey struct{}

func (t RoleDriftSurvey) impl() profileTechnique {
	return profileTechnique{id: "AC10RoleDriftSurvey", name: "RoleDriftSurvey", classID: "AC-10", summary: "build structured low-noise context for later chaining", risk: 0.24, impact: 0.34, eval: evalObserved}
}
func (t RoleDriftSurvey) ID() string                   { return t.impl().ID() }
func (t RoleDriftSurvey) Name() string                 { return t.impl().Name() }
func (t RoleDriftSurvey) ActionClassID() string        { return t.impl().ActionClassID() }
func (t RoleDriftSurvey) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t RoleDriftSurvey) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t RoleDriftSurvey) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t RoleDriftSurvey) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// EntitlementConsistencyCheck is pivot-heavy and links multiple graph hints to expose chained opportunities.
// Risk profile: medium risk due to directional probing that can trigger controls.
// Confidence rationale: requires evidence and hypothesis alignment to confirm a viable pivot.
// Expected behavior: unlocks credential, access, or lateral classes by revealing adjacency.
type EntitlementConsistencyCheck struct{}

func (t EntitlementConsistencyCheck) impl() profileTechnique {
	return profileTechnique{id: "AC10EntitlementConsistencyCheck", name: "EntitlementConsistencyCheck", classID: "AC-10", summary: "connect mid-stage observations into pivotable pathways", risk: 0.52, impact: 0.58, eval: evalPivot}
}
func (t EntitlementConsistencyCheck) ID() string                   { return t.impl().ID() }
func (t EntitlementConsistencyCheck) Name() string                 { return t.impl().Name() }
func (t EntitlementConsistencyCheck) ActionClassID() string        { return t.impl().ActionClassID() }
func (t EntitlementConsistencyCheck) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t EntitlementConsistencyCheck) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t EntitlementConsistencyCheck) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t EntitlementConsistencyCheck) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// PrivilegeChainPivot is a second pivot behavior focused on graph-link validation before escalation.
// Risk profile: medium-to-high because it exercises cross-node relationships.
// Confidence rationale: needs technique evidence plus supporting or enabling edges.
// Expected behavior: produces evidence that makes execution or impact classes evaluable.
type PrivilegeChainPivot struct{}

func (t PrivilegeChainPivot) impl() profileTechnique {
	return profileTechnique{id: "AC10PrivilegeChainPivot", name: "PrivilegeChainPivot", classID: "AC-10", summary: "stress pivot assumptions across linked graph paths", risk: 0.62, impact: 0.72, eval: evalGraphPivot}
}
func (t PrivilegeChainPivot) ID() string                   { return t.impl().ID() }
func (t PrivilegeChainPivot) Name() string                 { return t.impl().Name() }
func (t PrivilegeChainPivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t PrivilegeChainPivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t PrivilegeChainPivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t PrivilegeChainPivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t PrivilegeChainPivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// ControlPlanePivot models low-confidence high-impact behavior intended for rare but decisive opportunities.
// Risk profile: high operational and detection risk with potentially outsized downstream impact.
// Confidence rationale: only evaluates true for rare graph combinations across nodes and edges.
// Expected behavior: when triggered, it unlocks objective-oriented classes with strong score effects.
type ControlPlanePivot struct{}

func (t ControlPlanePivot) impl() profileTechnique {
	return profileTechnique{id: "AC10ControlPlanePivot", name: "ControlPlanePivot", classID: "AC-10", summary: "attempt rare high-consequence maneuver when graph strongly supports it", risk: 0.84, impact: 0.92, eval: evalRare}
}
func (t ControlPlanePivot) ID() string                   { return t.impl().ID() }
func (t ControlPlanePivot) Name() string                 { return t.impl().Name() }
func (t ControlPlanePivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t ControlPlanePivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t ControlPlanePivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t ControlPlanePivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t ControlPlanePivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// All returns the diversified technique set for AC-10.
func All() []model.Technique {
	return []model.Technique{
		PrivilegeAssessor{},
		RoleDriftSurvey{},
		EntitlementConsistencyCheck{},
		PrivilegeChainPivot{},
		ControlPlanePivot{},
	}
}

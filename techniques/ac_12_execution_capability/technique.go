package ac_12_execution_capability

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

// ExecutionCapabilityValidator favors passive or low-touch checks to reduce operational risk while building baselines.
// Risk profile: low confidence drift and low detection risk due to minimal interaction.
// Confidence rationale: high when little graph state exists because observations are easy to validate.
// Expected behavior: seeds follow-on discovery classes with foundational evidence.
type ExecutionCapabilityValidator struct{}

func (t ExecutionCapabilityValidator) impl() profileTechnique {
	return profileTechnique{id: "AC12ExecutionCapabilityValidator", name: "ExecutionCapabilityValidator", classID: "AC-12", summary: "validate in-environment execution capability", risk: 0.18, impact: 0.28, eval: evalObserved}
}
func (t ExecutionCapabilityValidator) ID() string                   { return t.impl().ID() }
func (t ExecutionCapabilityValidator) Name() string                 { return t.impl().Name() }
func (t ExecutionCapabilityValidator) ActionClassID() string        { return t.impl().ActionClassID() }
func (t ExecutionCapabilityValidator) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t ExecutionCapabilityValidator) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t ExecutionCapabilityValidator) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t ExecutionCapabilityValidator) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// BenignCommandCanary captures high-confidence low-impact context to enrich precision targeting later.
// Risk profile: low risk and low footprint, suitable for steady-state reconnaissance.
// Confidence rationale: requires at least one evidence node to avoid blind guesses.
// Expected behavior: unlocks service and protocol follow-ups with improved context quality.
type BenignCommandCanary struct{}

func (t BenignCommandCanary) impl() profileTechnique {
	return profileTechnique{id: "AC12BenignCommandCanary", name: "BenignCommandCanary", classID: "AC-12", summary: "build structured low-noise context for later chaining", risk: 0.24, impact: 0.34, eval: evalObserved}
}
func (t BenignCommandCanary) ID() string                   { return t.impl().ID() }
func (t BenignCommandCanary) Name() string                 { return t.impl().Name() }
func (t BenignCommandCanary) ActionClassID() string        { return t.impl().ActionClassID() }
func (t BenignCommandCanary) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t BenignCommandCanary) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t BenignCommandCanary) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t BenignCommandCanary) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// RuntimeConstraintProbe is pivot-heavy and links multiple graph hints to expose chained opportunities.
// Risk profile: medium risk due to directional probing that can trigger controls.
// Confidence rationale: requires evidence and hypothesis alignment to confirm a viable pivot.
// Expected behavior: unlocks credential, access, or lateral classes by revealing adjacency.
type RuntimeConstraintProbe struct{}

func (t RuntimeConstraintProbe) impl() profileTechnique {
	return profileTechnique{id: "AC12RuntimeConstraintProbe", name: "RuntimeConstraintProbe", classID: "AC-12", summary: "connect mid-stage observations into pivotable pathways", risk: 0.52, impact: 0.58, eval: evalPivot}
}
func (t RuntimeConstraintProbe) ID() string                   { return t.impl().ID() }
func (t RuntimeConstraintProbe) Name() string                 { return t.impl().Name() }
func (t RuntimeConstraintProbe) ActionClassID() string        { return t.impl().ActionClassID() }
func (t RuntimeConstraintProbe) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t RuntimeConstraintProbe) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t RuntimeConstraintProbe) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t RuntimeConstraintProbe) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// ToolTransferPivot is a second pivot behavior focused on graph-link validation before escalation.
// Risk profile: medium-to-high because it exercises cross-node relationships.
// Confidence rationale: needs technique evidence plus supporting or enabling edges.
// Expected behavior: produces evidence that makes execution or impact classes evaluable.
type ToolTransferPivot struct{}

func (t ToolTransferPivot) impl() profileTechnique {
	return profileTechnique{id: "AC12ToolTransferPivot", name: "ToolTransferPivot", classID: "AC-12", summary: "stress pivot assumptions across linked graph paths", risk: 0.62, impact: 0.72, eval: evalGraphPivot}
}
func (t ToolTransferPivot) ID() string                   { return t.impl().ID() }
func (t ToolTransferPivot) Name() string                 { return t.impl().Name() }
func (t ToolTransferPivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t ToolTransferPivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t ToolTransferPivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t ToolTransferPivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t ToolTransferPivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// SchedulerPivot models low-confidence high-impact behavior intended for rare but decisive opportunities.
// Risk profile: high operational and detection risk with potentially outsized downstream impact.
// Confidence rationale: only evaluates true for rare graph combinations across nodes and edges.
// Expected behavior: when triggered, it unlocks objective-oriented classes with strong score effects.
type SchedulerPivot struct{}

func (t SchedulerPivot) impl() profileTechnique {
	return profileTechnique{id: "AC12SchedulerPivot", name: "SchedulerPivot", classID: "AC-12", summary: "attempt rare high-consequence maneuver when graph strongly supports it", risk: 0.84, impact: 0.92, eval: evalRare}
}
func (t SchedulerPivot) ID() string                   { return t.impl().ID() }
func (t SchedulerPivot) Name() string                 { return t.impl().Name() }
func (t SchedulerPivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t SchedulerPivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t SchedulerPivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t SchedulerPivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t SchedulerPivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// All returns the diversified technique set for AC-12.
func All() []model.Technique {
	return []model.Technique{
		ExecutionCapabilityValidator{},
		BenignCommandCanary{},
		RuntimeConstraintProbe{},
		ToolTransferPivot{},
		SchedulerPivot{},
	}
}

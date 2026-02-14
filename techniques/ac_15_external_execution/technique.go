package ac_15_external_execution

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

// ExternalExecutionCoordinator favors passive or low-touch checks to reduce operational risk while building baselines.
// Risk profile: low confidence drift and low detection risk due to minimal interaction.
// Confidence rationale: high when little graph state exists because observations are easy to validate.
// Expected behavior: seeds follow-on discovery classes with foundational evidence.
type ExternalExecutionCoordinator struct{}

func (t ExternalExecutionCoordinator) impl() profileTechnique {
	return profileTechnique{id: "AC15ExternalExecutionCoordinator", name: "ExternalExecutionCoordinator", classID: "AC-15", summary: "coordinate actions requiring external execution dependencies", risk: 0.18, impact: 0.28, eval: evalObserved}
}
func (t ExternalExecutionCoordinator) ID() string                   { return t.impl().ID() }
func (t ExternalExecutionCoordinator) Name() string                 { return t.impl().Name() }
func (t ExternalExecutionCoordinator) ActionClassID() string        { return t.impl().ActionClassID() }
func (t ExternalExecutionCoordinator) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t ExternalExecutionCoordinator) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t ExternalExecutionCoordinator) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t ExternalExecutionCoordinator) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// VendorExecutionReadiness captures high-confidence low-impact context to enrich precision targeting later.
// Risk profile: low risk and low footprint, suitable for steady-state reconnaissance.
// Confidence rationale: requires at least one evidence node to avoid blind guesses.
// Expected behavior: unlocks service and protocol follow-ups with improved context quality.
type VendorExecutionReadiness struct{}

func (t VendorExecutionReadiness) impl() profileTechnique {
	return profileTechnique{id: "AC15VendorExecutionReadiness", name: "VendorExecutionReadiness", classID: "AC-15", summary: "build structured low-noise context for later chaining", risk: 0.24, impact: 0.34, eval: evalObserved}
}
func (t VendorExecutionReadiness) ID() string                   { return t.impl().ID() }
func (t VendorExecutionReadiness) Name() string                 { return t.impl().Name() }
func (t VendorExecutionReadiness) ActionClassID() string        { return t.impl().ActionClassID() }
func (t VendorExecutionReadiness) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t VendorExecutionReadiness) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t VendorExecutionReadiness) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t VendorExecutionReadiness) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// OutsourceTaskValidation is pivot-heavy and links multiple graph hints to expose chained opportunities.
// Risk profile: medium risk due to directional probing that can trigger controls.
// Confidence rationale: requires evidence and hypothesis alignment to confirm a viable pivot.
// Expected behavior: unlocks credential, access, or lateral classes by revealing adjacency.
type OutsourceTaskValidation struct{}

func (t OutsourceTaskValidation) impl() profileTechnique {
	return profileTechnique{id: "AC15OutsourceTaskValidation", name: "OutsourceTaskValidation", classID: "AC-15", summary: "connect mid-stage observations into pivotable pathways", risk: 0.52, impact: 0.58, eval: evalPivot}
}
func (t OutsourceTaskValidation) ID() string                   { return t.impl().ID() }
func (t OutsourceTaskValidation) Name() string                 { return t.impl().Name() }
func (t OutsourceTaskValidation) ActionClassID() string        { return t.impl().ActionClassID() }
func (t OutsourceTaskValidation) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t OutsourceTaskValidation) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t OutsourceTaskValidation) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t OutsourceTaskValidation) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// SupplyChainPivot is a second pivot behavior focused on graph-link validation before escalation.
// Risk profile: medium-to-high because it exercises cross-node relationships.
// Confidence rationale: needs technique evidence plus supporting or enabling edges.
// Expected behavior: produces evidence that makes execution or impact classes evaluable.
type SupplyChainPivot struct{}

func (t SupplyChainPivot) impl() profileTechnique {
	return profileTechnique{id: "AC15SupplyChainPivot", name: "SupplyChainPivot", classID: "AC-15", summary: "stress pivot assumptions across linked graph paths", risk: 0.62, impact: 0.72, eval: evalGraphPivot}
}
func (t SupplyChainPivot) ID() string                   { return t.impl().ID() }
func (t SupplyChainPivot) Name() string                 { return t.impl().Name() }
func (t SupplyChainPivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t SupplyChainPivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t SupplyChainPivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t SupplyChainPivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t SupplyChainPivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// RemoteOpsPivot models low-confidence high-impact behavior intended for rare but decisive opportunities.
// Risk profile: high operational and detection risk with potentially outsized downstream impact.
// Confidence rationale: only evaluates true for rare graph combinations across nodes and edges.
// Expected behavior: when triggered, it unlocks objective-oriented classes with strong score effects.
type RemoteOpsPivot struct{}

func (t RemoteOpsPivot) impl() profileTechnique {
	return profileTechnique{id: "AC15RemoteOpsPivot", name: "RemoteOpsPivot", classID: "AC-15", summary: "attempt rare high-consequence maneuver when graph strongly supports it", risk: 0.84, impact: 0.92, eval: evalRare}
}
func (t RemoteOpsPivot) ID() string                   { return t.impl().ID() }
func (t RemoteOpsPivot) Name() string                 { return t.impl().Name() }
func (t RemoteOpsPivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t RemoteOpsPivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t RemoteOpsPivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t RemoteOpsPivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t RemoteOpsPivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// All returns the diversified technique set for AC-15.
func All() []model.Technique {
	return []model.Technique{
		ExternalExecutionCoordinator{},
		VendorExecutionReadiness{},
		OutsourceTaskValidation{},
		SupplyChainPivot{},
		RemoteOpsPivot{},
	}
}

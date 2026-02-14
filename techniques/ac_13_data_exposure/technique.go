package ac_13_data_exposure

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

// DataExposureVerifier favors passive or low-touch checks to reduce operational risk while building baselines.
// Risk profile: low confidence drift and low detection risk due to minimal interaction.
// Confidence rationale: high when little graph state exists because observations are easy to validate.
// Expected behavior: seeds follow-on discovery classes with foundational evidence.
type DataExposureVerifier struct{}

func (t DataExposureVerifier) impl() profileTechnique {
	return profileTechnique{id: "AC13DataExposureVerifier", name: "DataExposureVerifier", classID: "AC-13", summary: "verify accessible sensitive data exposure", risk: 0.18, impact: 0.28, eval: evalObserved}
}
func (t DataExposureVerifier) ID() string                   { return t.impl().ID() }
func (t DataExposureVerifier) Name() string                 { return t.impl().Name() }
func (t DataExposureVerifier) ActionClassID() string        { return t.impl().ActionClassID() }
func (t DataExposureVerifier) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t DataExposureVerifier) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t DataExposureVerifier) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t DataExposureVerifier) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// PublicDataSampling captures high-confidence low-impact context to enrich precision targeting later.
// Risk profile: low risk and low footprint, suitable for steady-state reconnaissance.
// Confidence rationale: requires at least one evidence node to avoid blind guesses.
// Expected behavior: unlocks service and protocol follow-ups with improved context quality.
type PublicDataSampling struct{}

func (t PublicDataSampling) impl() profileTechnique {
	return profileTechnique{id: "AC13PublicDataSampling", name: "PublicDataSampling", classID: "AC-13", summary: "build structured low-noise context for later chaining", risk: 0.24, impact: 0.34, eval: evalObserved}
}
func (t PublicDataSampling) ID() string                   { return t.impl().ID() }
func (t PublicDataSampling) Name() string                 { return t.impl().Name() }
func (t PublicDataSampling) ActionClassID() string        { return t.impl().ActionClassID() }
func (t PublicDataSampling) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t PublicDataSampling) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t PublicDataSampling) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t PublicDataSampling) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// SchemaVisibilityCheck is pivot-heavy and links multiple graph hints to expose chained opportunities.
// Risk profile: medium risk due to directional probing that can trigger controls.
// Confidence rationale: requires evidence and hypothesis alignment to confirm a viable pivot.
// Expected behavior: unlocks credential, access, or lateral classes by revealing adjacency.
type SchemaVisibilityCheck struct{}

func (t SchemaVisibilityCheck) impl() profileTechnique {
	return profileTechnique{id: "AC13SchemaVisibilityCheck", name: "SchemaVisibilityCheck", classID: "AC-13", summary: "connect mid-stage observations into pivotable pathways", risk: 0.52, impact: 0.58, eval: evalPivot}
}
func (t SchemaVisibilityCheck) ID() string                   { return t.impl().ID() }
func (t SchemaVisibilityCheck) Name() string                 { return t.impl().Name() }
func (t SchemaVisibilityCheck) ActionClassID() string        { return t.impl().ActionClassID() }
func (t SchemaVisibilityCheck) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t SchemaVisibilityCheck) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t SchemaVisibilityCheck) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t SchemaVisibilityCheck) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// DataLinkagePivot is a second pivot behavior focused on graph-link validation before escalation.
// Risk profile: medium-to-high because it exercises cross-node relationships.
// Confidence rationale: needs technique evidence plus supporting or enabling edges.
// Expected behavior: produces evidence that makes execution or impact classes evaluable.
type DataLinkagePivot struct{}

func (t DataLinkagePivot) impl() profileTechnique {
	return profileTechnique{id: "AC13DataLinkagePivot", name: "DataLinkagePivot", classID: "AC-13", summary: "stress pivot assumptions across linked graph paths", risk: 0.62, impact: 0.72, eval: evalGraphPivot}
}
func (t DataLinkagePivot) ID() string                   { return t.impl().ID() }
func (t DataLinkagePivot) Name() string                 { return t.impl().Name() }
func (t DataLinkagePivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t DataLinkagePivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t DataLinkagePivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t DataLinkagePivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t DataLinkagePivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// BackupChannelPivot models low-confidence high-impact behavior intended for rare but decisive opportunities.
// Risk profile: high operational and detection risk with potentially outsized downstream impact.
// Confidence rationale: only evaluates true for rare graph combinations across nodes and edges.
// Expected behavior: when triggered, it unlocks objective-oriented classes with strong score effects.
type BackupChannelPivot struct{}

func (t BackupChannelPivot) impl() profileTechnique {
	return profileTechnique{id: "AC13BackupChannelPivot", name: "BackupChannelPivot", classID: "AC-13", summary: "attempt rare high-consequence maneuver when graph strongly supports it", risk: 0.84, impact: 0.92, eval: evalRare}
}
func (t BackupChannelPivot) ID() string                   { return t.impl().ID() }
func (t BackupChannelPivot) Name() string                 { return t.impl().Name() }
func (t BackupChannelPivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t BackupChannelPivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t BackupChannelPivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t BackupChannelPivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t BackupChannelPivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// All returns the diversified technique set for AC-13.
func All() []model.Technique {
	return []model.Technique{
		DataExposureVerifier{},
		PublicDataSampling{},
		SchemaVisibilityCheck{},
		DataLinkagePivot{},
		BackupChannelPivot{},
	}
}

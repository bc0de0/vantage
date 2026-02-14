package ac_06_version_enumeration

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

// VersionEnumerator favors passive or low-touch checks to reduce operational risk while building baselines.
// Risk profile: low confidence drift and low detection risk due to minimal interaction.
// Confidence rationale: high when little graph state exists because observations are easy to validate.
// Expected behavior: seeds follow-on discovery classes with foundational evidence.
type VersionEnumerator struct{}

func (t VersionEnumerator) impl() profileTechnique {
	return profileTechnique{id: "AC06VersionEnumerator", name: "VersionEnumerator", classID: "AC-06", summary: "enumerate service versions and capabilities", risk: 0.18, impact: 0.28, eval: evalObserved}
}
func (t VersionEnumerator) ID() string                   { return t.impl().ID() }
func (t VersionEnumerator) Name() string                 { return t.impl().Name() }
func (t VersionEnumerator) ActionClassID() string        { return t.impl().ActionClassID() }
func (t VersionEnumerator) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t VersionEnumerator) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t VersionEnumerator) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t VersionEnumerator) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// PatchCadenceSnapshot captures high-confidence low-impact context to enrich precision targeting later.
// Risk profile: low risk and low footprint, suitable for steady-state reconnaissance.
// Confidence rationale: requires at least one evidence node to avoid blind guesses.
// Expected behavior: unlocks service and protocol follow-ups with improved context quality.
type PatchCadenceSnapshot struct{}

func (t PatchCadenceSnapshot) impl() profileTechnique {
	return profileTechnique{id: "AC06PatchCadenceSnapshot", name: "PatchCadenceSnapshot", classID: "AC-06", summary: "build structured low-noise context for later chaining", risk: 0.24, impact: 0.34, eval: evalObserved}
}
func (t PatchCadenceSnapshot) ID() string                   { return t.impl().ID() }
func (t PatchCadenceSnapshot) Name() string                 { return t.impl().Name() }
func (t PatchCadenceSnapshot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t PatchCadenceSnapshot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t PatchCadenceSnapshot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t PatchCadenceSnapshot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t PatchCadenceSnapshot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// BuildFingerprintSampler is pivot-heavy and links multiple graph hints to expose chained opportunities.
// Risk profile: medium risk due to directional probing that can trigger controls.
// Confidence rationale: requires evidence and hypothesis alignment to confirm a viable pivot.
// Expected behavior: unlocks credential, access, or lateral classes by revealing adjacency.
type BuildFingerprintSampler struct{}

func (t BuildFingerprintSampler) impl() profileTechnique {
	return profileTechnique{id: "AC06BuildFingerprintSampler", name: "BuildFingerprintSampler", classID: "AC-06", summary: "connect mid-stage observations into pivotable pathways", risk: 0.52, impact: 0.58, eval: evalPivot}
}
func (t BuildFingerprintSampler) ID() string                   { return t.impl().ID() }
func (t BuildFingerprintSampler) Name() string                 { return t.impl().Name() }
func (t BuildFingerprintSampler) ActionClassID() string        { return t.impl().ActionClassID() }
func (t BuildFingerprintSampler) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t BuildFingerprintSampler) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t BuildFingerprintSampler) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t BuildFingerprintSampler) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// ChangelogPivot is a second pivot behavior focused on graph-link validation before escalation.
// Risk profile: medium-to-high because it exercises cross-node relationships.
// Confidence rationale: needs technique evidence plus supporting or enabling edges.
// Expected behavior: produces evidence that makes execution or impact classes evaluable.
type ChangelogPivot struct{}

func (t ChangelogPivot) impl() profileTechnique {
	return profileTechnique{id: "AC06ChangelogPivot", name: "ChangelogPivot", classID: "AC-06", summary: "stress pivot assumptions across linked graph paths", risk: 0.62, impact: 0.72, eval: evalGraphPivot}
}
func (t ChangelogPivot) ID() string                   { return t.impl().ID() }
func (t ChangelogPivot) Name() string                 { return t.impl().Name() }
func (t ChangelogPivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t ChangelogPivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t ChangelogPivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t ChangelogPivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t ChangelogPivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// DependencyLineagePivot models low-confidence high-impact behavior intended for rare but decisive opportunities.
// Risk profile: high operational and detection risk with potentially outsized downstream impact.
// Confidence rationale: only evaluates true for rare graph combinations across nodes and edges.
// Expected behavior: when triggered, it unlocks objective-oriented classes with strong score effects.
type DependencyLineagePivot struct{}

func (t DependencyLineagePivot) impl() profileTechnique {
	return profileTechnique{id: "AC06DependencyLineagePivot", name: "DependencyLineagePivot", classID: "AC-06", summary: "attempt rare high-consequence maneuver when graph strongly supports it", risk: 0.84, impact: 0.92, eval: evalRare}
}
func (t DependencyLineagePivot) ID() string                   { return t.impl().ID() }
func (t DependencyLineagePivot) Name() string                 { return t.impl().Name() }
func (t DependencyLineagePivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t DependencyLineagePivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t DependencyLineagePivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t DependencyLineagePivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t DependencyLineagePivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// All returns the diversified technique set for AC-06.
func All() []model.Technique {
	return []model.Technique{
		VersionEnumerator{},
		PatchCadenceSnapshot{},
		BuildFingerprintSampler{},
		ChangelogPivot{},
		DependencyLineagePivot{},
	}
}

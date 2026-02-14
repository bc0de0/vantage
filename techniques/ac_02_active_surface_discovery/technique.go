package ac_02_active_surface_discovery

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

// SurfaceProbe favors passive or low-touch checks to reduce operational risk while building baselines.
// Risk profile: low confidence drift and low detection risk due to minimal interaction.
// Confidence rationale: high when little graph state exists because observations are easy to validate.
// Expected behavior: seeds follow-on discovery classes with foundational evidence.
type SurfaceProbe struct{}

func (t SurfaceProbe) impl() profileTechnique {
	return profileTechnique{id: "AC02SurfaceProbe", name: "SurfaceProbe", classID: "AC-02", summary: "probe target surface for reachable hosts and endpoints", risk: 0.18, impact: 0.28, eval: evalObserved}
}
func (t SurfaceProbe) ID() string                   { return t.impl().ID() }
func (t SurfaceProbe) Name() string                 { return t.impl().Name() }
func (t SurfaceProbe) ActionClassID() string        { return t.impl().ActionClassID() }
func (t SurfaceProbe) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t SurfaceProbe) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t SurfaceProbe) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t SurfaceProbe) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// AssetCensusSweep captures high-confidence low-impact context to enrich precision targeting later.
// Risk profile: low risk and low footprint, suitable for steady-state reconnaissance.
// Confidence rationale: requires at least one evidence node to avoid blind guesses.
// Expected behavior: unlocks service and protocol follow-ups with improved context quality.
type AssetCensusSweep struct{}

func (t AssetCensusSweep) impl() profileTechnique {
	return profileTechnique{id: "AC02AssetCensusSweep", name: "AssetCensusSweep", classID: "AC-02", summary: "build structured low-noise context for later chaining", risk: 0.24, impact: 0.34, eval: evalObserved}
}
func (t AssetCensusSweep) ID() string                   { return t.impl().ID() }
func (t AssetCensusSweep) Name() string                 { return t.impl().Name() }
func (t AssetCensusSweep) ActionClassID() string        { return t.impl().ActionClassID() }
func (t AssetCensusSweep) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t AssetCensusSweep) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t AssetCensusSweep) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t AssetCensusSweep) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// InternetEdgeSampler is pivot-heavy and links multiple graph hints to expose chained opportunities.
// Risk profile: medium risk due to directional probing that can trigger controls.
// Confidence rationale: requires evidence and hypothesis alignment to confirm a viable pivot.
// Expected behavior: unlocks credential, access, or lateral classes by revealing adjacency.
type InternetEdgeSampler struct{}

func (t InternetEdgeSampler) impl() profileTechnique {
	return profileTechnique{id: "AC02InternetEdgeSampler", name: "InternetEdgeSampler", classID: "AC-02", summary: "connect mid-stage observations into pivotable pathways", risk: 0.52, impact: 0.58, eval: evalPivot}
}
func (t InternetEdgeSampler) ID() string                   { return t.impl().ID() }
func (t InternetEdgeSampler) Name() string                 { return t.impl().Name() }
func (t InternetEdgeSampler) ActionClassID() string        { return t.impl().ActionClassID() }
func (t InternetEdgeSampler) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t InternetEdgeSampler) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t InternetEdgeSampler) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t InternetEdgeSampler) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// AdjacentSubnetPivot is a second pivot behavior focused on graph-link validation before escalation.
// Risk profile: medium-to-high because it exercises cross-node relationships.
// Confidence rationale: needs technique evidence plus supporting or enabling edges.
// Expected behavior: produces evidence that makes execution or impact classes evaluable.
type AdjacentSubnetPivot struct{}

func (t AdjacentSubnetPivot) impl() profileTechnique {
	return profileTechnique{id: "AC02AdjacentSubnetPivot", name: "AdjacentSubnetPivot", classID: "AC-02", summary: "stress pivot assumptions across linked graph paths", risk: 0.62, impact: 0.72, eval: evalGraphPivot}
}
func (t AdjacentSubnetPivot) ID() string                   { return t.impl().ID() }
func (t AdjacentSubnetPivot) Name() string                 { return t.impl().Name() }
func (t AdjacentSubnetPivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t AdjacentSubnetPivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t AdjacentSubnetPivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t AdjacentSubnetPivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t AdjacentSubnetPivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// ShadowAssetPivot models low-confidence high-impact behavior intended for rare but decisive opportunities.
// Risk profile: high operational and detection risk with potentially outsized downstream impact.
// Confidence rationale: only evaluates true for rare graph combinations across nodes and edges.
// Expected behavior: when triggered, it unlocks objective-oriented classes with strong score effects.
type ShadowAssetPivot struct{}

func (t ShadowAssetPivot) impl() profileTechnique {
	return profileTechnique{id: "AC02ShadowAssetPivot", name: "ShadowAssetPivot", classID: "AC-02", summary: "attempt rare high-consequence maneuver when graph strongly supports it", risk: 0.84, impact: 0.92, eval: evalRare}
}
func (t ShadowAssetPivot) ID() string                   { return t.impl().ID() }
func (t ShadowAssetPivot) Name() string                 { return t.impl().Name() }
func (t ShadowAssetPivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t ShadowAssetPivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t ShadowAssetPivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t ShadowAssetPivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t ShadowAssetPivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// All returns the diversified technique set for AC-02.
func All() []model.Technique {
	return []model.Technique{
		SurfaceProbe{},
		AssetCensusSweep{},
		InternetEdgeSampler{},
		AdjacentSubnetPivot{},
		ShadowAssetPivot{},
	}
}

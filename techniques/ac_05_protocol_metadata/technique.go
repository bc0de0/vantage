package ac_05_protocol_metadata

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

// ProtocolMetadataInspector favors passive or low-touch checks to reduce operational risk while building baselines.
// Risk profile: low confidence drift and low detection risk due to minimal interaction.
// Confidence rationale: high when little graph state exists because observations are easy to validate.
// Expected behavior: seeds follow-on discovery classes with foundational evidence.
type ProtocolMetadataInspector struct{}

func (t ProtocolMetadataInspector) impl() profileTechnique {
	return profileTechnique{id: "AC05ProtocolMetadataInspector", name: "ProtocolMetadataInspector", classID: "AC-05", summary: "inspect protocol banners and metadata", risk: 0.18, impact: 0.28, eval: evalObserved}
}
func (t ProtocolMetadataInspector) ID() string                   { return t.impl().ID() }
func (t ProtocolMetadataInspector) Name() string                 { return t.impl().Name() }
func (t ProtocolMetadataInspector) ActionClassID() string        { return t.impl().ActionClassID() }
func (t ProtocolMetadataInspector) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t ProtocolMetadataInspector) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t ProtocolMetadataInspector) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t ProtocolMetadataInspector) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// HandshakeParameterAudit captures high-confidence low-impact context to enrich precision targeting later.
// Risk profile: low risk and low footprint, suitable for steady-state reconnaissance.
// Confidence rationale: requires at least one evidence node to avoid blind guesses.
// Expected behavior: unlocks service and protocol follow-ups with improved context quality.
type HandshakeParameterAudit struct{}

func (t HandshakeParameterAudit) impl() profileTechnique {
	return profileTechnique{id: "AC05HandshakeParameterAudit", name: "HandshakeParameterAudit", classID: "AC-05", summary: "build structured low-noise context for later chaining", risk: 0.24, impact: 0.34, eval: evalObserved}
}
func (t HandshakeParameterAudit) ID() string                   { return t.impl().ID() }
func (t HandshakeParameterAudit) Name() string                 { return t.impl().Name() }
func (t HandshakeParameterAudit) ActionClassID() string        { return t.impl().ActionClassID() }
func (t HandshakeParameterAudit) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t HandshakeParameterAudit) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t HandshakeParameterAudit) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t HandshakeParameterAudit) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// CipherPreferenceSampler is pivot-heavy and links multiple graph hints to expose chained opportunities.
// Risk profile: medium risk due to directional probing that can trigger controls.
// Confidence rationale: requires evidence and hypothesis alignment to confirm a viable pivot.
// Expected behavior: unlocks credential, access, or lateral classes by revealing adjacency.
type CipherPreferenceSampler struct{}

func (t CipherPreferenceSampler) impl() profileTechnique {
	return profileTechnique{id: "AC05CipherPreferenceSampler", name: "CipherPreferenceSampler", classID: "AC-05", summary: "connect mid-stage observations into pivotable pathways", risk: 0.52, impact: 0.58, eval: evalPivot}
}
func (t CipherPreferenceSampler) ID() string                   { return t.impl().ID() }
func (t CipherPreferenceSampler) Name() string                 { return t.impl().Name() }
func (t CipherPreferenceSampler) ActionClassID() string        { return t.impl().ActionClassID() }
func (t CipherPreferenceSampler) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t CipherPreferenceSampler) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t CipherPreferenceSampler) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t CipherPreferenceSampler) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// ProtocolDowngradePivot is a second pivot behavior focused on graph-link validation before escalation.
// Risk profile: medium-to-high because it exercises cross-node relationships.
// Confidence rationale: needs technique evidence plus supporting or enabling edges.
// Expected behavior: produces evidence that makes execution or impact classes evaluable.
type ProtocolDowngradePivot struct{}

func (t ProtocolDowngradePivot) impl() profileTechnique {
	return profileTechnique{id: "AC05ProtocolDowngradePivot", name: "ProtocolDowngradePivot", classID: "AC-05", summary: "stress pivot assumptions across linked graph paths", risk: 0.62, impact: 0.72, eval: evalGraphPivot}
}
func (t ProtocolDowngradePivot) ID() string                   { return t.impl().ID() }
func (t ProtocolDowngradePivot) Name() string                 { return t.impl().Name() }
func (t ProtocolDowngradePivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t ProtocolDowngradePivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t ProtocolDowngradePivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t ProtocolDowngradePivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t ProtocolDowngradePivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// MetadataLeakPivot models low-confidence high-impact behavior intended for rare but decisive opportunities.
// Risk profile: high operational and detection risk with potentially outsized downstream impact.
// Confidence rationale: only evaluates true for rare graph combinations across nodes and edges.
// Expected behavior: when triggered, it unlocks objective-oriented classes with strong score effects.
type MetadataLeakPivot struct{}

func (t MetadataLeakPivot) impl() profileTechnique {
	return profileTechnique{id: "AC05MetadataLeakPivot", name: "MetadataLeakPivot", classID: "AC-05", summary: "attempt rare high-consequence maneuver when graph strongly supports it", risk: 0.84, impact: 0.92, eval: evalRare}
}
func (t MetadataLeakPivot) ID() string                   { return t.impl().ID() }
func (t MetadataLeakPivot) Name() string                 { return t.impl().Name() }
func (t MetadataLeakPivot) ActionClassID() string        { return t.impl().ActionClassID() }
func (t MetadataLeakPivot) Evaluate(g *model.Graph) bool { return t.impl().Evaluate(g) }
func (t MetadataLeakPivot) Execute(ctx context.Context, g *model.Graph) (model.Evidence, error) {
	return t.impl().Execute(ctx, g)
}
func (t MetadataLeakPivot) RiskModifier() float64   { return t.impl().RiskModifier() }
func (t MetadataLeakPivot) ImpactModifier() float64 { return t.impl().ImpactModifier() }

// All returns the diversified technique set for AC-05.
func All() []model.Technique {
	return []model.Technique{
		ProtocolMetadataInspector{},
		HandshakeParameterAudit{},
		CipherPreferenceSampler{},
		ProtocolDowngradePivot{},
		MetadataLeakPivot{},
	}
}

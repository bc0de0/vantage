package ac_08_credential_validation

import (
	"context"

	"vantage/techniques/model"
)

// CredentialValidator implements action class AC-08.
type CredentialValidator struct{}

func (t CredentialValidator) ID() string            { return "CredentialValidator" }
func (t CredentialValidator) Name() string          { return "CredentialValidator" }
func (t CredentialValidator) ActionClassID() string { return "AC-08" }
func (t CredentialValidator) Evaluate(graph *model.Graph) bool {
	return graph != nil && graph.TechniqueNodes > 0
}
func (t CredentialValidator) Execute(_ context.Context, _ *model.Graph) (model.Evidence, error) {
	return model.Evidence{TechniqueID: t.ID(), Summary: "validate credential material against target auth interfaces", Success: true}, nil
}
func (t CredentialValidator) RiskModifier() float64   { return 0.8 }
func (t CredentialValidator) ImpactModifier() float64 { return 0.8 }

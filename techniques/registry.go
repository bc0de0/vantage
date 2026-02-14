package techniques

import (
	"fmt"

	"vantage/techniques/ac_01_passive_observation"
	"vantage/techniques/ac_02_active_surface_discovery"
	"vantage/techniques/ac_03_reachability_validation"
	"vantage/techniques/ac_04_service_identification"
	"vantage/techniques/ac_05_protocol_metadata"
	"vantage/techniques/ac_06_version_enumeration"
	"vantage/techniques/ac_07_auth_surface_analysis"
	"vantage/techniques/ac_08_credential_validation"
	"vantage/techniques/ac_09_access_establishment"
	"vantage/techniques/ac_10_privilege_assessment"
	"vantage/techniques/ac_11_lateral_reachability"
	"vantage/techniques/ac_12_execution_capability"
	"vantage/techniques/ac_13_data_exposure"
	"vantage/techniques/ac_14_impact_feasibility"
	"vantage/techniques/ac_15_external_execution"
)

// RegisterAll returns the complete static technique set keyed by Technique.ID().
func RegisterAll() map[string]Technique {
	all := []Technique{
		ac_01_passive_observation.PassiveDNSCollection{},
		ac_02_active_surface_discovery.SurfaceProbe{},
		ac_03_reachability_validation.ReachabilityValidator{},
		ac_04_service_identification.ServiceIdentifier{},
		ac_05_protocol_metadata.ProtocolMetadataInspector{},
		ac_06_version_enumeration.VersionEnumerator{},
		ac_07_auth_surface_analysis.AuthSurfaceAnalyzer{},
		ac_08_credential_validation.CredentialValidator{},
		ac_09_access_establishment.AccessEstablisher{},
		ac_10_privilege_assessment.PrivilegeAssessor{},
		ac_11_lateral_reachability.LateralReachabilityAnalyzer{},
		ac_12_execution_capability.ExecutionCapabilityValidator{},
		ac_13_data_exposure.DataExposureVerifier{},
		ac_14_impact_feasibility.ImpactFeasibilityAssessor{},
		ac_15_external_execution.ExternalExecutionCoordinator{},
	}

	registry := make(map[string]Technique, len(all))
	for _, t := range all {
		if _, exists := registry[t.ID()]; exists {
			panic(fmt.Sprintf("duplicate technique id: %s", t.ID()))
		}
		registry[t.ID()] = t
	}
	return registry
}

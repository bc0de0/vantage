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
	all := make([]Technique, 0, 75)
	all = append(all, ac_01_passive_observation.All()...)
	all = append(all, ac_02_active_surface_discovery.All()...)
	all = append(all, ac_03_reachability_validation.All()...)
	all = append(all, ac_04_service_identification.All()...)
	all = append(all, ac_05_protocol_metadata.All()...)
	all = append(all, ac_06_version_enumeration.All()...)
	all = append(all, ac_07_auth_surface_analysis.All()...)
	all = append(all, ac_08_credential_validation.All()...)
	all = append(all, ac_09_access_establishment.All()...)
	all = append(all, ac_10_privilege_assessment.All()...)
	all = append(all, ac_11_lateral_reachability.All()...)
	all = append(all, ac_12_execution_capability.All()...)
	all = append(all, ac_13_data_exposure.All()...)
	all = append(all, ac_14_impact_feasibility.All()...)
	all = append(all, ac_15_external_execution.All()...)

	registry := make(map[string]Technique, len(all))
	for _, t := range all {
		if _, exists := registry[t.ID()]; exists {
			panic(fmt.Sprintf("duplicate technique id: %s", t.ID()))
		}
		registry[t.ID()] = t
	}
	return registry
}

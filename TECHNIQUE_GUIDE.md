# Technique Guide

This module applies a **MECE** design:

- **Mutually Exclusive:** every technique binds to exactly one canonical ActionClass ID.
- **Collectively Exhaustive:** all ActionClass IDs from `action-classes-normalized/` have at least one technique implementation.

| ActionClass ID | Technique | Description | Risk Modifier | Impact Modifier |
|---|---|---|---:|---:|
| AC-01 | PassiveDNSCollection | Collect passive DNS and CT evidence. | 0.20 | 0.40 |
| AC-02 | SurfaceProbe | Probe target surface for reachable hosts/endpoints. | 0.50 | 0.60 |
| AC-03 | ReachabilityValidator | Validate host reachability. | 0.40 | 0.60 |
| AC-04 | ServiceIdentifier | Identify exposed services. | 0.50 | 0.70 |
| AC-05 | ProtocolMetadataInspector | Inspect protocol metadata/banners. | 0.40 | 0.70 |
| AC-06 | VersionEnumerator | Enumerate versions and capabilities. | 0.60 | 0.80 |
| AC-07 | AuthSurfaceAnalyzer | Analyze authentication surfaces. | 0.60 | 0.70 |
| AC-08 | CredentialValidator | Validate credential material. | 0.80 | 0.80 |
| AC-09 | AccessEstablisher | Establish authenticated access. | 0.90 | 0.90 |
| AC-10 | PrivilegeAssessor | Assess obtained privileges. | 0.70 | 0.85 |
| AC-11 | LateralReachabilityAnalyzer | Assess lateral movement reachability. | 0.80 | 0.85 |
| AC-12 | ExecutionCapabilityValidator | Validate execution capability. | 0.85 | 0.90 |
| AC-13 | DataExposureVerifier | Verify data exposure conditions. | 0.70 | 0.95 |
| AC-14 | ImpactFeasibilityAssessor | Assess impact feasibility. | 0.90 | 0.95 |
| AC-15 | ExternalExecutionCoordinator | Coordinate externally-required execution. | 0.75 | 0.80 |

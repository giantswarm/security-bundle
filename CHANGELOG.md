# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Renamed `clusterName` value to `clusterID`, aligning with other bundles and the values provided by `cluster-operator`.
- Updated `Kyverno-Policies` version to 0.19.0 to enable restricted PSS by default.

## [0.14.3] - 2023-05-17

### Changed

- Update to `kyverno` (app) version 0.14.5, introducing a webhooks cleanup job when uninstalling Kyverno chart.

## [0.14.2] - 2023-05-09

### Changed

- Update to `trivy` (app) version 0.8.1, introducing a `PolicyException` and making `VerticalPodAutoscaler` configurable.

## [0.14.1] - 2023-05-04

### Changed

- Update to `falco` (app) version 0.5.2, moving the `PolicyExceptions` to the `giantswarm` namespace.

## [0.14.0] - 2023-05-03

### Changed

- Changed the value type of configmap/secret values to object rather than multi line string.
- Update to `kyverno` (app) version 0.14.4, restricting the namespaces where `PolicyExceptions` may be created.
- Update to `trivy-operator` (app) version 0.4.0, adding a `Cilium Network Policy` for `trivy-operator`.

## [0.13.0] - 2023-04-12

### Changed

- Renamed `security-pack` to `security-bundle`.
- Changed default installation namespace from `security-pack` to `security-bundle`.
- Renamed `kyverno-policies` value to `kyvernoPolicies`.
- Renamed `falco-app` chart name to `falco`.
- Removed `security-bundle` from the `playground` catalog. **Users must now install `security-bundle` from the `giantswarm` catalog**.
- Removed `starboard` (app). **Users must now install `trivy-operator` as a replacement**.
- Update to `kyverno` (app) version 0.14.3, containing upstream `kyverno` version 1.9.2.
- Update to `kyverno-policies` (app) version 0.18.1, containing upstream `kyverno-policies` version 1.7.5.
- Update to `trivy-operator` (app) version 0.3.7, containing upstream `trivy-operator` version 0.12.0.
- Update to `trivy` (app) version 0.8.0, containing upstream `trivy` version 0.37.2.
- Update to `starboard-exporter` (app) version 0.7.3.
- Update to `falco` (app) version 0.5.1.

## [0.12.0] - 2023-02-08

### Changed

- Renames `falco` app `chartName` & bumps version to `0.5.0`.
- Promotes `jiralert` app to the Giant Swarm catalog.

## [0.11.0] - 2023-01-11

### Changed

- Move `kyverno` (app) to its own namespace.
- Update to `starboard-exporter` version 0.7.0.
- Remove `security-pack-helper`. Its current logic is no longer needed with Kyverno 1.8.0 and above.

## [0.10.0] - 2022-12-21

### Added
- Add icon url to `Chart.yaml`

### Changed

- Update to `Falco` (app) version 0.4.3, this version adds support for `VerticalPodAutoscaler` (0.4.2), and makes use of the `falco-driverless` image.
- Update to `Trivy` (app) version 0.7.1, this version adds support for `VerticalPodAutoscaler`.
- Update to `trivy-operator` (app) version 0.3.2, containing upstream `trivy-operator` version 0.7.1.
- Update to `Kyverno` (app) version 0.13.1, containing upstream `Kyverno` version 1.8.4.

## [0.9.0] - 2022-11-10

### Changed

- Update to `Trivy` (app) version 0.7.0, this version changes the `trivy-app` chart name to `trivy`.
- Update to `trivy-operator` (app) version 0.2.0, this version changes the `serverUrl` value from `trivy-app` to `trivy`.
- Update to `Kyverno` (app) version 0.11.8, this version includes a `CiliumNetworkPolicy` for the `kyverno-crd-install` job.
- Update to `Falco` (app) version 0.4.0, containing upstream `falco` version 0.33.0 as well as `falco-exporter` upstream version 0.8.0 and `falcosidekick` upstream version 2.26.0.

## [0.8.1] - 2022-10-25

### Changed

- Update to `Kyverno` (app) version 0.11.6, containing upstream policy-reporter version 2.10.1.
- Update to `starboard-exporter` version 0.6.2, enabling custom metric logic via `ServiceMonitor` relabeling configuration.

## [0.8.0] - 2022-09-27

### Added

- Push to `giantswarm` catalog in addition to `playground` catalog.
- Add (currently optional) `security-pack-helper` app.

### Changed

- Update to `Kyverno` (app) version 0.11.2, containing upstream version 1.7.3.
- Update to `Trivy` (app) version 0.6.0, containing upstream version 0.30.4.
- Update to `starboard-exporter` version 0.6.0.
- Update to `trivy-operator` (app) version 0.2.0, containing upstream version 0.2.1.

### Removed

- Disable `starboard-app` by default. Starboard has been replaced with Trivy Operator. Starboard can still be enabled in the pack configuration values.

## [0.7.0] - 2022-09-06

## Added

- Enable optional installation of `trivy-operator`.

## [0.6.0] - 2022-08-17

### Changed

- Update to Kyverno (app) version 0.11.0 containing upstream version 1.7.2 and resilience improvements.
- Update to `starboard-exporter` version 0.5.1, adding support for selectively enabling/disabling metrics for each report type.
- Update to Trivy (app) version 0.4.0 containing upstream version 0.28.1.

## [0.5.0] - 2022-08-01

### Changed

- Update `starboard-app` version 0.8.0, adding a `PriorityLevelConfig` and `FlowSchema` for starboard's API requests.

## [0.4.0] - 2022-07-08

### Changed

- Update to `starboard-exporter` version 0.5.0, adding support for exporting Starboard `CISKubeBenchReport` resources.
- Update to `jiralert` version 0.0.3, adding support for autoresolving tickets.
- Fix `Jiralert` catalog.

## [0.3.1] - 2022-06-10

### Added

- Adding `cluster` flag to the `kgs` command generating App CR.

### Changed

- Update to Kyverno (app) version 0.10.1 which fixes the CRDs deployed to match upstream version 1.6.2.

## [0.3.0] - 2022-05-17

### Added

- Add `jiralert-app` 0.0.2 to the security pack for creating Jira tickets.

### Changed

- Update to `starboard-exporter` version 0.4.1, adding spread re-queueing of reports by +/- 10% by default to help smooth resource utilization.

## [0.2.0] - 2022-04-25

### Changed

- Update to Starboard (app) version 0.7.1 containing upstream version 0.15.3. This release introduces support for `ClusterComplianceReport` generation and includes a benchmark for the NSA + CISA Kubernetes Hardening Guide.
- Update to `starboard-exporter` version 0.4.0, adding support for exporting Starboard `ConfigAuditReport` resources.
- Update to Trivy (app) version 0.3.0 containing upstream version 0.25.0.

## [0.1.0] - 2022-04-07

### Added

- Enable `kyverno` installation by default.
- Update to Kyverno (app) version 0.10.0 containing upstream version 1.6.2.
- Add `kyverno-policies` 0.17.1 to the security pack for enforcing Kubernetes Pod Security Standards (PSS).

### Changed

- Use Falco app version 0.3.2.

## [0.0.1] - 2022-03-24

- Initial release containing (Giant Swarm apps) Falco 0.3.1, Kyverno 0.9.1, Starboard 0.6.0, Starboard exporter 0.3.1, and Trivy 0.2.0.

[Unreleased]: https://github.com/giantswarm/security-bundle/compare/v0.14.3...HEAD
[0.14.3]: https://github.com/giantswarm/security-bundle/compare/v0.14.2...v0.14.3
[0.14.2]: https://github.com/giantswarm/security-bundle/compare/v0.14.1...v0.14.2
[0.14.1]: https://github.com/giantswarm/security-bundle/compare/v0.14.0...v0.14.1
[0.14.0]: https://github.com/giantswarm/security-bundle/compare/v0.13.0...v0.14.0
[0.13.0]: https://github.com/giantswarm/security-bundle/compare/v0.12.0...v0.13.0
[0.12.0]: https://github.com/giantswarm/security-pack/compare/v0.11.0...v0.12.0
[0.11.0]: https://github.com/giantswarm/security-pack/compare/v0.10.0...v0.11.0
[0.10.0]: https://github.com/giantswarm/security-pack/compare/v0.9.0...v0.10.0
[0.9.0]: https://github.com/giantswarm/security-pack/compare/v0.8.1...v0.9.0
[0.8.1]: https://github.com/giantswarm/security-pack/compare/v0.8.0...v0.8.1
[0.8.0]: https://github.com/giantswarm/security-pack/compare/v0.7.0...v0.8.0
[0.7.0]: https://github.com/giantswarm/security-pack/compare/v0.6.0...v0.7.0
[0.6.0]: https://github.com/giantswarm/security-pack/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/giantswarm/security-pack/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/giantswarm/security-pack/compare/v0.3.1...v0.4.0
[0.3.1]: https://github.com/giantswarm/security-pack/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/giantswarm/security-pack/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/giantswarm/security-pack/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/giantswarm/security-pack/compare/v0.0.1...v0.1.0
[0.0.1]: https://github.com/giantswarm/security-pack/releases/tag/v0.0.1

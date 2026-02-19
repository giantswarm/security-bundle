# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Add `io.giantswarm.application.audience` and `io.giantswarm.application.managed` chart annotations for Backstage visibility.

### Changed

- Migrate sub-apps from App CRs to Flux HelmRelease CRs.
- Add pre-upgrade migration hook to safely clean up Chart CRs during the transition.
- Update `kyverno` (app) to v0.24.0.
  - This release includes a Kyverno upstream update. Please refer to the following Release Notes from upstream for the latest changes:
    - https://github.com/kyverno/kyverno/releases/tag/v1.17.0
- Update `kyverno-crds` (app) to v1.17.0.
- Update `kyverno-policies` (app) to v0.25.0.
- Update `kubescape` (app) to v0.0.6.
- Migrate chart annotations to OCI-compatible format (change `application.giantswarm.io/team` to `io.giantswarm.application.team`, `application.giantswarm.io/app-type` to `io.giantswarm.application.app-type`).

## [1.17.0] - 2026-01-29

### Changed

- Update `kyverno` (app) to v0.23.0.
- Update `kyverno-crds` (app) to v1.16.0.
- Update `reports-server` (app) to v0.1.0.
- Update `cloudnative-pg` (app) to v0.0.13.
- Update `kubescape` (app) to v0.0.5.
- Update `starboard-exporter` (app) to v1.0.2.

## [1.16.1] - 2025-12-18

### Changed

- Add missing dependency to all apps.

## [1.16.0] - 2025-12-18

### Changed

- Allow to set multiple dependencies on the depends-on annotation.
- Rename `edgedb` to `gel`.
- Update `cloudnative-pg` (app) to v0.0.12.
- Update `gel` (app) to v1.0.1.

## [1.15.0] - 2025-11-06

### Added

- Add `kubescape` (app) version v0.0.4.

### Changed

- Update `kyverno` (app) to v0.21.1.
- Update `kyverno-crds` (app) to v1.15.0.

### Notes

This release includes a Kyverno upstream update. Please refer to the following Release Notes from upstream for the latest changes:

## [1.14.0] - 2025-10-28

### Changed

- Update `kyverno` (app) to v0.20.1.
- Update `kyverno-crds` (app) to v1.14.0.
- Update `kyverno-policies` (app) to v0.24.0.
- Update `reports-server` (app) to v0.0.3.

### Notes

This release includes a Kyverno upstream update. Please refer to the following Release Notes from upstream for the latest changes:

## [1.13.1] - 2025-10-21

### Changed

- Revert previous `kyverno` update ([#536](https://github.com/giantswarm/security-bundle/pull/536), [#531](https://github.com/giantswarm/security-bundle/pull/531), [#538](https://github.com/giantswarm/security-bundle/pull/538)).
- Update `kyverno-policy-operator` (app) to v0.1.6.

### Notes

This release reverts the Kyverno update performed in [1.13.0], but keeps the following app updates:
- Update `trivy-operator` (app) to v0.12.1.
- Update `trivy` (app) to v0.14.1.
- Update `falco` (app) to v0.11.0.

## [1.13.0] - 2025-10-08

### Changed

- Update `kyverno` (app) to v0.20.0.
- Update `kyverno-crds` (app) to v1.14.0.
- Update `kyverno-policies` (app) to v0.24.0.
- Update `kyverno-policy-operator` (app) to v0.1.5.
- Update `trivy-operator` (app) to v0.12.1.
- Update `trivy` (app) to v0.14.1.
- Update `falco` (app) to v0.11.0.

### Notes

This release includes a Kyverno upstream update. Please refer to the following Release Notes from upstream for the latest changes:

## [1.12.0] - 2025-07-30

### Changed

- Update `trivy-operator` (app) to v0.11.1.
- Update `trivy` (app) to v0.14.0.
- Update `falco` (app) to v0.10.1.
- Update `cloudnative-pg` (app) to v0.0.10.
- Update `starboard-exporter` (app) to v0.8.2.
- Updated E2E tests to use apptest-framework v1.14.0

## [1.11.0] - 2025-05-28

### Added

- Add `policy-api-crds` app to manage Policy API CRDs.

### Changed

- Update `trivy` (app) to v0.13.4.
- Update `cloudnative-pg` (app) to v0.0.7.
- Update `starboard-exporter` (app) to v0.8.1.
- Update `kyverno-policy-operator` (app) to v0.0.11.
- Update `cloudnative-pg` (app) to v0.0.9.

## [1.10.1] - 2025-04-29

### Changed

- Update `kyverno-crds` (app) to v1.13.1.

### Notes

**Note:** Kyverno `PolicyExceptions` (API group `kyverno.io`) versions `v2alpha1` and `v2beta1` are deprecated and will be removed in the next Kyverno minor release (v1.14). Please update all Kyverno PolicyExceptions to `v2`. No action is required for Giant Swarm Policy API `PolicyExceptions` (API group `policy.giantswarm.io`), which are handled automatically.

## [1.10.0] - 2025-02-28

### Added

- Add e2e tests for the `security-bundle` and all is components

### Changed

- Update `kyverno` (app) to v0.19.0.
- Update `kyverno-crds` (app) to v1.13.0.
- Update `kyverno-policies` (app) to v0.23.0.
- Update `edgedb` (app) to v0.1.0.
- Update `falco` (app) to v0.10.0.
- Update `trivy` (app) to v0.13.2.

## [1.9.1] - 2024-11-13

### Changed

- Update `trivy-operator` (app) to v0.10.3.
- Update `trivy` (app) to v0.13.1.

## [1.9.0] - 2024-10-30

### Changed

- Update `kyverno` (app) to v0.18.1.
- Update `kyverno-crds` (app) to v1.12.0.
- Update `kyverno-policies` (app) to v0.21.0.
- Update `starboard-exporter` (app) to v0.8.0.
- Update `trivy-operator` (app) to v0.10.2.
- Update `trivy` (app) to v0.13.0.
- Update `falco` (app) to v0.9.1.

### Breaking Changes

**Note:** When upgrading to this security-bundle version with Falco enabled, the Falco App will fail to upgrade due to a breaking change in the upstream chart. To finish the upgrade, disable, then re-enable the Falco App by setting `apps.falco.enabled=[false|true]` [in the security-bundle user values Config Map](https://github.com/giantswarm/security-bundle/tree/main?tab=readme-ov-file#configuring).

## [1.8.2] - 2024-08-29

### Changed

- Update `cloudnative-pg` (app) to v0.0.6.
- Update `trivy-operator` (app) to v0.10.0.
- Update `kyverno-policy-operator` (app) to v0.0.8.
- Update `kyverno` (app) to v0.17.16.

## [1.8.1] - 2024-08-14

### Changed

- Update `trivy-operator` (app) to v0.9.1.

## [1.8.0] - 2024-07-19

### Added

- Add `kyverno-crds` app to handle Kyverno CRD install.

### Changed

- Update `kyverno` (app) to v0.17.15. This version disables the CRD install job in favor of `kyverno-crds` App.

## [1.7.1] - 2024-06-13

### Changed

- Update `kyverno` (app) to v0.17.14.
- Update `starboard-exporter` (app) to v0.7.11.

## [1.7.0] - 2024-06-06

### Added

- Add `cloudnative-pg`, `edgedb`, and `reports-server` apps (disabled).

### Changed

- Update `trivy` (app) to v0.12.0.
- Update `trivy-operator` (app) to v0.9.0.
- Update `cloudnative-pg` (app) to v0.0.5.

## [1.6.7] - 2024-05-23

### Changed

- Update `kyverno` (app) to v0.17.13.

## [1.6.6] - 2024-05-23

### Changed

- Update `kyverno` (app) to v0.17.12.
- Update `trivy-operator` (app) to v0.8.1.

## [1.6.5] - 2024-05-07

### Changed

- Update `starboard-exporter` (app) to v0.7.10.

## [1.6.4] - 2024-05-03

### Added

- Add `extraConfigs` settings for the bundle Apps.

### Changed

- Update `trivy` (app) to v0.11.0.
- Update `falco` (app) to v0.8.1.
- Update `kyverno` (app) to v0.17.10.
- Update `starboard-exporter` (app) to v0.7.9.

## [1.6.3] - 2024-04-04

### Changed

- Update `kyverno` (app) to v0.17.9.
- Update `trivy` (app) to v0.10.1.

## [1.6.2] - 2024-02-22

### Changed

- Update to kyverno (app) to v0.17.6.

## [1.6.1] - 2024-02-08

### Changed

- Update to kyverno (app) to v0.17.5.
- Update to exception-recommender (app) to v0.1.1.
- Update to trivy-operator (app) to v0.7.2.

## [1.6.0] - 2024-01-26

### Changed

- Update to exception-recommender (app) to v0.1.0.
- Update to falco (app) to v0.8.0.
- Update to kyverno-policy-operator (app) version v0.0.7.
- Update to kyverno (app) version v0.17.2.
- Update to starboard-exporter (app) version v0.7.8.
- Update to trivy-operator (app) to v0.7.0.
- Update to trivy (app) to v0.10.0.

## [1.5.0] - 2023-12-19

### Added

- Add a `global.namespace` value for automatically setting all app namespaces.

### Changed

- Update to exception-recommender (app) to v0.0.7.
- Update to falco (app) to v0.7.1.
- Update to jiralert (app) version v0.1.3.
- Update to kyverno-policy-operator (app) version v0.0.6.
- Update to starboard-exporter (app) version v0.7.7.
- Update to trivy-operator (app) to v0.5.1.

## [1.4.2] - 2023-12-07

### Added

- Add `options` to individual app settings to allow custom timeout values.

### Changed

- Update to `kyverno-app` (app) version v0.16.4.
- Update to `kyverno-policies` (app) version v0.20.2.
- Update to `exception-recommender` (app) to v0.0.6.
- Update to `starboard-exporter` (app) version v0.7.5.

## [1.4.1] - 2023-11-29

### Changed

- Update to `kyverno` (app) version 0.16.3.

## [1.4.0] - 2023-11-16

### Changed

- Revert namespace change of `exception-recommender` and `kyverno-policy-operator`.
- Update to `kyverno` (app) version 0.16.2.
- Update to `kyverno-policy-operator` (app) v0.0.5.
- Update to `exception-recommender` (app) v0.0.3.
- Update to `falco` (app) to v0.7.0.

## [1.3.0] - 2023-10-31

### Changed

- Update to `kyverno` (app) version 0.16.1.
- Update to `kyverno-policy-operator` (app) version 0.0.4.
- Update to `trivy-operator` (app) version 0.5.0.
- Update to `trivy` (app) version 0.9.0.
- Update to `jiralert` (app) version 0.1.2.
- Update to `falco` (app) version 0.6.7.
- Disable PSPs for `falco-exporter`.

## [1.2.0] - 2023-10-24

### Changed

- Change `kyverno-policy-operator` and `exception-recommender` namespaces to `policy-exceptions`.
- Update to `kyverno-policy-operator` (app) version 0.0.3.
- Update to `kyverno` (app) version 0.16.0.

## [1.1.0] - 2023-10-10

### Added

- Add `exception-recommender` (app) to the security bundle to create Giant Swarm PolicyException recommendations.
- Add `kyverno-policy-operator` (app) to the security bundle to automatically create Kyverno PolicyExceptions from Giant Swarm PolicyExceptions.

### Changed

- Update to `kyverno` (app) upstream version 1.10.2. *Note:* This update includes breaking changes in the values structure, please check the [migration docs](https://github.com/giantswarm/kyverno-app/tree/main/helm/kyverno/charts/kyverno#new-chart-values) before upgrading.
- Update to `kyverno-policies` (app) version 0.20.1.
- Update to `trivy-operator` (app) to version 0.4.1.
- Update to `trivy` (app) version 0.8.3.
- Update to `falco` (app) version 0.6.5.

## [1.0.2] - 2023-07-05

### Added

- Add `depends-on` annotation to set App depenencies.

### Changed

- Update to `kyverno` (app) version 0.14.10, introducing Cilium Network Policy support.
- Update to `kyverno-policies` (app) version 0.20.0.

## [1.0.1] - 2023-06-26

### Changed

- Update to `kyverno` (app) version 0.14.9, introducing the `scrapeTimeout` and `interval` configurable values for Policy Reporter `ServiceMonitors`.

## [1.0.0] - 2023-06-22

### Changed

- Set `Kyverno Policies` config to `Enforce` mode.
- Disable PSPs for all Apps from `userConfig`.
- Updated `kyverno` (App) to version 0.14.8, which modifies the PSPs logic so it can be disabled for every component.

## [0.16.0] - 2023-06-13

### Changed

- Update to `kyverno` (app) version 0.14.7, introducing exception mechanisms for `chart-operator` and restricting wildcards for Kinds.
- Disabled the default apps `falco`, `trivy`, `trivy-operator` and `starboard-exporter`. This apps can be manually enabled.

## [0.15.0] - 2023-05-31

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

### Added

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

### Added

- Initial release containing (Giant Swarm apps) Falco 0.3.1, Kyverno 0.9.1, Starboard 0.6.0, Starboard exporter 0.3.1, and Trivy 0.2.0.

[Unreleased]: https://github.com/giantswarm/security-bundle/compare/v1.17.0...HEAD
[1.17.0]: https://github.com/giantswarm/security-bundle/compare/v1.16.1...v1.17.0
[1.16.1]: https://github.com/giantswarm/security-bundle/compare/v1.16.0...v1.16.1
[1.16.0]: https://github.com/giantswarm/security-bundle/compare/v1.15.0...v1.16.0
[1.15.0]: https://github.com/giantswarm/security-bundle/compare/v1.14.0...v1.15.0
[1.14.0]: https://github.com/giantswarm/security-bundle/compare/v1.13.1...v1.14.0
[1.13.1]: https://github.com/giantswarm/security-bundle/compare/v1.13.0...v1.13.1
[1.13.0]: https://github.com/giantswarm/security-bundle/compare/v1.12.0...v1.13.0
[1.12.0]: https://github.com/giantswarm/security-bundle/compare/v1.11.0...v1.12.0
[1.11.0]: https://github.com/giantswarm/security-bundle/compare/v1.10.1...v1.11.0
[1.10.1]: https://github.com/giantswarm/security-bundle/compare/v1.10.0...v1.10.1
[1.10.0]: https://github.com/giantswarm/security-bundle/compare/v1.9.1...v1.10.0
[1.9.1]: https://github.com/giantswarm/security-bundle/compare/v1.9.0...v1.9.1
[1.9.0]: https://github.com/giantswarm/security-bundle/compare/v1.8.2...v1.9.0
[1.8.2]: https://github.com/giantswarm/security-bundle/compare/v1.8.1...v1.8.2
[1.8.1]: https://github.com/giantswarm/security-bundle/compare/v1.8.0...v1.8.1
[1.8.0]: https://github.com/giantswarm/security-bundle/compare/v1.7.1...v1.8.0
[1.7.1]: https://github.com/giantswarm/security-bundle/compare/v1.7.0...v1.7.1
[1.7.0]: https://github.com/giantswarm/security-bundle/compare/v1.6.7...v1.7.0
[1.6.7]: https://github.com/giantswarm/security-bundle/compare/v1.6.6...v1.6.7
[1.6.6]: https://github.com/giantswarm/security-bundle/compare/v1.6.5...v1.6.6
[1.6.5]: https://github.com/giantswarm/security-bundle/compare/v1.6.4...v1.6.5
[1.6.4]: https://github.com/giantswarm/security-bundle/compare/v1.6.3...v1.6.4
[1.6.3]: https://github.com/giantswarm/security-bundle/compare/v1.6.2...v1.6.3
[1.6.2]: https://github.com/giantswarm/security-bundle/compare/v1.6.1...v1.6.2
[1.6.1]: https://github.com/giantswarm/security-bundle/compare/v1.6.0...v1.6.1
[1.6.0]: https://github.com/giantswarm/security-bundle/compare/v1.5.0...v1.6.0
[1.5.0]: https://github.com/giantswarm/security-bundle/compare/v1.4.2...v1.5.0
[1.4.2]: https://github.com/giantswarm/security-bundle/compare/v1.4.1...v1.4.2
[1.4.1]: https://github.com/giantswarm/security-bundle/compare/v1.4.0...v1.4.1
[1.4.0]: https://github.com/giantswarm/security-bundle/compare/v1.3.0...v1.4.0
[1.3.0]: https://github.com/giantswarm/security-bundle/compare/v1.2.0...v1.3.0
[1.2.0]: https://github.com/giantswarm/security-bundle/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/giantswarm/security-bundle/compare/v1.0.2...v1.1.0
[1.0.2]: https://github.com/giantswarm/security-bundle/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/giantswarm/security-bundle/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/giantswarm/security-bundle/compare/v0.17.1...v1.0.0
[0.17.1]: https://github.com/giantswarm/security-bundle/compare/v0.17.0...v0.17.1
[0.17.0]: https://github.com/giantswarm/security-bundle/compare/v0.16.0...v0.17.0
[0.16.0]: https://github.com/giantswarm/security-bundle/compare/v0.15.0...v0.16.0
[0.15.0]: https://github.com/giantswarm/security-bundle/compare/v0.14.3...v0.15.0
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

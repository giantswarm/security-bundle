# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/giantswarm/security-pack/compare/v0.3.1...HEAD
[0.3.1]: https://github.com/giantswarm/security-pack/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/giantswarm/security-pack/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/giantswarm/security-pack/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/giantswarm/security-pack/compare/v0.0.1...v0.1.0
[0.0.1]: https://github.com/giantswarm/security-pack/releases/tag/v0.0.1

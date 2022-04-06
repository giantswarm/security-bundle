# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Enable `kyverno` installation by default.
- Update to Kyverno (app) version 0.10.0 containing upstream version 1.6.2.
- Add `kyverno-policies` to the security pack for enforcing Kubernetes Pod Security Standards (PSS).

### Changed

- Use Falco app version 0.3.2.

## [0.0.1] - 2022-03-24

- Initial release containing (Giant Swarm apps) Falco 0.3.1, Kyverno 0.9.1, Starboard 0.6.0, Starboard exporter 0.3.1, and Trivy 0.2.0.

[Unreleased]: https://github.com/giantswarm/security-pack/compare/v0.0.1...HEAD
[0.0.1]: https://github.com/giantswarm/security-pack/releases/tag/v0.0.1

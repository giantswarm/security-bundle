# Security Bundle Observability Labels

This document describes the observability labels feature in the security-bundle that enables log collection from security components.

## Overview

The security-bundle can automatically add the `app.kubernetes.io/part-of: "observability"` label to security components to enable log collection in Loki. This feature is based on the Giant Swarm observability documentation: https://docs.giantswarm.io/tutorials/observability/data-ingestion/logs/

## Features

- **Automatic labeling**: Adds observability labels to deployed security components
- **Selective control**: Each component can be individually enabled/disabled
- **Security-aware**: Excludes sensitive components by default
- **Post-deployment**: Uses Helm hooks to apply labels after components are deployed
- **Error handling**: Gracefully handles missing components

## Configuration

### Basic Usage

Enable observability for all security components:

```yaml
global:
  observability:
    enabled: true
```

### Advanced Configuration

Selectively enable/disable components:

```yaml
global:
  observability:
    enabled: true
    components:
      kyverno:
        admissionController: true
        backgroundController: true
        cleanupController: true
        policyReporter: true
        policyReporterPlugin: true
        policyReporterUI: true
        reportsController: true
      securityBundle:
        policyOperator: true
      giantswarm:
        athena: true
        dex: true
        falcoExporter: true
        falcoSidekick: true
        falcoMetacollector: true
        trivy: true
        trivyOperator: true
```

## Components Covered

### Kyverno Namespace
- `kyverno-admission-controller` (Deployment)
- `kyverno-background-controller` (Deployment)
- `kyverno-cleanup-controller` (Deployment)
- `kyverno-policy-reporter` (Deployment)
- `kyverno-policy-reporter-kyverno-plugin` (Deployment)
- `kyverno-policy-reporter-ui` (Deployment)
- `kyverno-reports-controller` (Deployment)

### Security-Bundle Namespace
- `kyverno-policy-operator` (Deployment)

### Giantswarm Namespace
- `athena` (Deployment)
- `dex` (Deployment)
- `falco-exporter` (DaemonSet)
- `falco-sidekick` (Deployment)
- `falco-k8s-metacollector` (Deployment)
- `trivy` (StatefulSet)
- `trivy-operator` (Deployment)

## Security Considerations

The following components are **excluded** from log collection for security reasons:

- **Falco**: Contains customer security data and sensitive information
- **Starboard-exporter**: Contains sensitive vulnerability data

These exclusions are built into the template and cannot be overridden to maintain security.

## Examples

See `examples/observability-logging.yaml` for a complete configuration example.

## Troubleshooting

### Job Fails to Run

Check the job logs:
```bash
kubectl logs -n <namespace> job/security-bundle-add-observability-labels
```

### Components Not Labeled

Verify the component exists:
```bash
kubectl get deployment <component-name> -n <namespace>
```

Check if the label was applied:
```bash
kubectl get deployment <component-name> -n <namespace> -o yaml | grep "app.kubernetes.io/part-of"
```

### Disable Observability

To disable the feature entirely:
```yaml
global:
  observability:
    enabled: false
```
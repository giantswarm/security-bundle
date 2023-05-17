## Examples

This directory contains example `user-values` files.

### Kyverno and Kyverno Policies with PSS Restricted policies

This file disables all default apps excluding Kyverno and Kyverno policies. Additionally, it sets the `kyvernoPolicies` user values to enable the `restricted` PSS policies in `Audit` mode.

**_The `clusterName` and `organization` values need to be updated with your cluster and organization._**

The Security Bundle App CR can then be templated with the following command:
```shell
kubectl gs template app \
--catalog giantswarm \
--name security-bundle \
--app-name demo01-security-bundle \ # Replace with cluster ID
--in-cluster \
--cluster-name demo01 \ # Replace with cluster ID
--target-namespace demo01 \ # Replace with cluster ID
--version 0.14.3 \ # Replace with security-bundle version
--user-configmap kyverno-pss.yaml > outerApp.yaml
```

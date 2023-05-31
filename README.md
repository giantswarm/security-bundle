[![CircleCI](https://circleci.com/gh/giantswarm/security-bundle.svg?style=shield)](https://circleci.com/gh/giantswarm/security-bundle)

# Giant Swarm Security Bundle

Giant Swarm offers a [managed security bundle][security-bundle] which provides an unintrusive baseline for security observability and enforcement in Kubernetes clusters. This App is a convenient wrapper containing multiple other Apps which make up the security bundle. See our full [App Bundle reference][app-bundle] to learn more.

By default, installing the security bundle in a cluster includes:

- Falco, from our [`falco-app`][falco-app]
- Kyverno, from our [`kyverno-app`][kyverno-app]
  - our [`kyverno-policies`][kyverno-policies] app for Kubernetes Pod Security Standards (PSS)
- Trivy, from our [`trivy-app`][trivy-app]
- Trivy Operator, from our [`trivy-operator`][trivy-operator-app] app
  - our [`starboard-exporter`][starboard-exporter] for exposing metrics

Some optional components are also installable from this bundle, including:

- Jiralert, from our [`jiralert-app`][jiralert-app], for automatically creating Jira issues from security findings

Previous versions of the pack included Starboard, from our [`starboard-app`][starboard-app]. Starboard has been deprecated in favor of Trivy Operator, and we have removed the Starboard app from this app bundle as of v0.13.0.

Apps can be selectively enabled or disabled using the `enabled` setting for that app in the `security-bundle` Helm values.

More information and configuration options can be found in each app repository.

**Note:** There is a known issue when uninstalling Kyverno where some resources may not be removed properly. The cause is still under investigation, but a [recent update](https://github.com/kyverno/kyverno/issues/3111) has improved the situation so that Kyverno's webhooks are properly removed, meaning the remnants are unnecessary but should otherwise have no negative effect on the cluster. If installing Kyverno through this App, you may need to manually remove some lingering Kyverno resources if you subsequently choose to remove Kyverno.

## Installing

:warning: **Existing `security-pack` users must delete the old `security-pack` CR first before installing the bundle.** It is not possible to update directly from a `security-pack` to a `security-bundle` App CR by renaming it.

### Updating from `security-pack`

To change an existing `security-pack` install to a `security-bundle`, the following changes must be made:
- any overrides to the Apps `configMap` or `secret` inside the `userConfig` key must be switched from string to object.
- any overrides to the `kyverno-policies` App in the `security-bundle` App values configuration must be replaced with equivalent overrides under the `kyvernoPolicies` key. The key `kyverno-policies` has been renamed to `kyvernoPolicies` only to simplify its usage in Helm. The name of the `kyverno-policies` App itself is unchanged.
- if using the default installation namespace (`security-pack`), any logic which depends on that namespace must be updated to reference the new default namespaces (`security-bundle`). If setting a custom installation namespace, no change is required.
- if the existing `security-pack` App CR is installing from the `playground` catalog, the catalog must be changed to `giantswarm`. The `security-bundle` will not be pushed to the `playground` catalog.
- after the above changes have been made, the old `security-pack` CR must be deleted before the new `security-bundle` CR can then be created.


This "App of Apps" method is rather new and our UX tooling is still catching up, so our normal App installation methods may or may not work for you depending on your management cluster and component versions.

The currently recommended way to install the security bundle is:

1. Create `user-values.yaml` containing the name of the cluster where the Apps should be installed, and the organization where that cluster is running:

    ```yaml
    clusterID: demo1
    organization: demo-team
    ```

2. Use `kubectl gs` to template the "outer" App CR:

    ```shell
    $ kubectl gs template app \
    --catalog giantswarm \
    --name security-bundle \
    --app-name demo01-security-bundle \
    --in-cluster \
    --cluster-name demo01 \
    --target-namespace demo01 \
    --version 0.0.1 \
    --user-configmap user-values.yaml > outerApp.yaml
    ```

3. Apply the generated App CR and ConfigMap to the management cluster:

    ```shell
    $ kubectl --context=<your-mc> apply -f outerApp.yaml
    configmap/security-bundle-userconfig created
    app.application.giantswarm.io/security-bundle created
    ```

Support for these methods are not yet officially supported, but may still work:

1. [Using our web interface](https://docs.giantswarm.io/ui-api/web/app-platform/#installing-an-app)
2. [Using our API](https://docs.giantswarm.io/api/#operation/createClusterAppV5)

### **Important**

If you are not using `kubectl gs` plugin, plese remember to ensure the correct label: `app-operator.giantswarm.io/version: 0.0.0` is set on the App CR. Missing this configuration will result with stuck deployment of an app.

When naming the App CR, please make sure the name is unique within the Management Cluster, using just `security-bundle`
name for two or more App CRs may lead to unexpected behavior. It is recommended to use cluster name as a prefix or suffix,
for example `demo01-security-bundle` or `security-bundle-demo1`.

## Configuring

### values.yaml

**This is an example of a values file you could upload using our web interface.**

```yaml
# values.yaml
clusterID: demo1
organization: demo-team
```

### Sample App CR and ConfigMap for the management cluster

If you have access to the Kubernetes API on the management cluster, you could create
the App CR and ConfigMap directly.

Here is an example that would install the app to
workload cluster `abc12`:

```yaml
# appCR.yaml
apiVersion: application.giantswarm.io/v1alpha1
kind: App
metadata:
  labels:
    app-operator.giantswarm.io/version: 0.0.0
    giantswarm.io/cluster: demo1
  name: security-bundle
  namespace: demo1
spec:
  catalog: giantswarm
  kubeConfig:
    inCluster: true
  name: security-bundle
  namespace: demo1
  userConfig:
    configMap:
      name: security-bundle-userconfig
      namespace: demo1
  version: 0.0.1
```

```yaml
# user-values-configmap.yaml
apiVersion: v1
data:
  values: |
    clusterID: demo1
    organization: giantswarm
kind: ConfigMap
metadata:
  creationTimestamp: null
  name: security-bundle-userconfig
  namespace: demo1
```

See our [full reference page on how to configure applications](https://docs.giantswarm.io/app-platform/app-configuration/) for more details.

[app-bundle]: https://docs.giantswarm.io/getting-started/app-platform/app-bundle/
[falco-app]: https://github.com/giantswarm/falco-app
[jiralert-app]: https://github.com/giantswarm/jiralert-app
[kyverno-app]: https://github.com/giantswarm/kyverno-app
[kyverno-policies]: https://github.com/giantswarm/kyverno-policies/
[security-bundle]: https://docs.giantswarm.io/app-platform/apps/security/
[starboard-app]: https://github.com/giantswarm/starboard-app
[starboard-exporter]: https://github.com/giantswarm/starboard-exporter/
[trivy-app]: https://github.com/giantswarm/trivy-app/
[trivy-operator-app]: https://github.com/giantswarm/trivy-operator-app

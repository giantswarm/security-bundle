[![CircleCI](https://circleci.com/gh/giantswarm/security-bundle.svg?style=shield)](https://circleci.com/gh/giantswarm/security-bundle)

# Giant Swarm Security Bundle

Giant Swarm offers a [managed security bundle][security-bundle] which provides an unintrusive baseline for security observability and enforcement in Kubernetes clusters. This App is a convenient wrapper containing multiple other Apps which make up the security bundle. See our full [App Bundle reference][app-bundle] to learn more.

By default, installing the security bundle in a cluster includes:

- Kyverno, from our [`kyverno-app`][kyverno-app]
  - our [`kyverno-policies`][kyverno-policies] app for Kubernetes Pod Security Standards (PSS)

- components supporting our Policy API features
  - our [`kyverno-policy-operator`][kyverno-policy-operator] app, which orchestrates our managed Kyverno policies.
  - our [`exception-recommender`][exception-recommender] app, which recommends possible Giant Swarm PolicyExceptions for non-compliant workloads.

Some optional components are also installable from this bundle, including:

- Falco, from our [`falco-app`][falco-app]
- Trivy, from our [`trivy-app`][trivy-app]
- Trivy Operator, from our [`trivy-operator`][trivy-operator-app] app
  - our [`starboard-exporter`][starboard-exporter] for exposing metrics
- Jiralert, from our [`jiralert-app`][jiralert-app], for automatically creating Jira issues from security findings

Several additional components deployable from the bundle are under development for future platform features. These are platform-internal and not intended for direct customer use, but can be optionally enabled to test upcoming improvements to report storage.

- CloudNative PG, from our [`cloudnative-pg-app`][cnpg-app]
- EdgeDB, from our [`edgedb-app`][edgedb-app]
- Kyverno Reports Server, from our [`reports-server-app`][reports-server-app]

Previous versions of the pack included Starboard, from our [`starboard-app`][starboard-app]. Starboard has been deprecated in favor of Trivy Operator, and we have removed the Starboard app from this app bundle as of v0.13.0.

Apps can be selectively enabled or disabled using the `enabled` setting for that app in the `security-bundle` Helm values.

More information and configuration options can be found in each app repository.

## Installing

:warning: **In version `v1.0.0` PSPs are disabled by default. Clusters running versions older than `1.25.0` must enable the PSPs in the userconfig of the `values.yaml` file before installing the Security Bundle or use older `v0.x.x` versions.**

### Compatibility Matrix

| Bundle Version  | K8s Version  | GS Release  | Branch  | PSS Policy State  | PSPs installed |
|:---:|:---:|:---:|:---:|:---:|:---:|
| v1.x.x  |  >= v1.25.0 | >= v20.0.0  | `main`  | enforce  | no  |
| v1.x.x  |  v1.24.x | >= v19.3.0, < v20.0.0  | `main`  | enforce  | no |
| v0.x.x  | < v1.25.0 | >= v19.1.0, < v19.3.0  | `legacy`  | audit  | yes  |

### Upgrading from a self-managed to a preinstalled `security-bundle`

The `security-bundle` is now being installed by default in new Giant Swarm cluster versions.

When upgrading from a cluster where the bundle was not preinstalled, it is possible that the installation of the bundle will fail if the bundle itself, or one of its apps (like Kyverno) was installed as an optional app prior to the upgrade.

We are working on an automated way to resolve this condition, but due to technical limitations and variation between how customers manage Apps (e.g. gitops) we currently recommend uninstalling any customer-installed `security-bundle`, `kyverno-app`, and `kyverno-policies` Apps installed in a cluster when upgrading to a version containing the bundle by default. 

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
[cnpg-app]: https://github.com/giantswarm/cloudnative-pg-app
[edgedb-app]: https://github.com/giantswarm/edgedb-app
[exception-recommender]: https://github.com/giantswarm/exception-recommender
[falco-app]: https://github.com/giantswarm/falco-app
[jiralert-app]: https://github.com/giantswarm/jiralert-app
[kyverno-app]: https://github.com/giantswarm/kyverno-app
[kyverno-policies]: https://github.com/giantswarm/kyverno-policies/
[kyverno-policy-operator]: https://github.com/giantswarm/kyverno-policy-operator
[reports-server-app]: https://github.com/giantswarm/reports-server-app
[security-bundle]: https://docs.giantswarm.io/app-platform/apps/security/
[starboard-app]: https://github.com/giantswarm/starboard-app
[starboard-exporter]: https://github.com/giantswarm/starboard-exporter/
[trivy-app]: https://github.com/giantswarm/trivy-app/
[trivy-operator-app]: https://github.com/giantswarm/trivy-operator-app

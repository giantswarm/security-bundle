[![CircleCI](https://circleci.com/gh/giantswarm/security-pack.svg?style=shield)](https://circleci.com/gh/giantswarm/security-pack)

# security-pack chart

Giant Swarm offers a security-pack App which can be installed in workload clusters. This App is a convenient wrapper containing multiple other Apps composing our security pack.

**Note:** There is a [known issue](https://github.com/kyverno/kyverno/issues/3111) when uninstalling `kyverno` which is pending upstream release. If installing `kyverno` through this App, you may need to manually Kyverno's `validatingwebhookconfigurations` and `mutatingwebhookconfigurations` if you subsequently uninstall the App prior to the 1.7.0 release.

## Installing

This "App of Apps" method is rather new and our UX tooling is still catching up, so our normal App installation methods may or may not work for you depending on your management cluster and component versions.

The currently recommended way to install the security pack is:

1. Create `user-values.yaml` containing the name of the cluster where the Apps should be installed, and the organization where that cluster is running:

    ```yaml
    clusterName: demo1
    organization: demo-team
    ```

2. Use `kubectl gs` to template the "outer" App CR:

    ```shell
    $ kubectl gs template app \
    --catalog giantswarm \
    --name security-pack \
    --in-cluster \
    --namespace demo1 \ 
    --version 0.0.1 \
    --user-configmap user-values.yaml > outerApp.yaml
    ```

3. Apply the generated App CR and ConfigMap to the management cluster:

    ```shell
    $ kubectl --context=<your-mc> apply -f outerApp.yaml
    configmap/security-pack-userconfig created
    app.application.giantswarm.io/security-pack created
    ```

Support for these methods are not yet officially supported, but may still work:

1. [Using our web interface](https://docs.giantswarm.io/ui-api/web/app-platform/#installing-an-app)
2. [Using our API](https://docs.giantswarm.io/api/#operation/createClusterAppV5)

## Configuring

### values.yaml

**This is an example of a values file you could upload using our web interface.**

```yaml
# values.yaml
clusterName: demo1
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
  name: security-pack
  namespace: demo1
spec:
  catalog: giantswarm
  kubeConfig:
    inCluster: true
  name: security-pack
  namespace: demo1
  userConfig:
    configMap:
      name: security-pack-userconfig
      namespace: demo1
  version: 0.0.1
```

```yaml
# user-values-configmap.yaml
apiVersion: v1
data:
  values: |
    clusterName: demo1
    organization: giantswarm
kind: ConfigMap
metadata:
  creationTimestamp: null
  name: security-pack-userconfig
  namespace: demo1
```

See our [full reference page on how to configure applications](https://docs.giantswarm.io/app-platform/app-configuration/) for more details.

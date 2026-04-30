package basic

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/giantswarm/apiextensions-application/api/v1alpha1"
	"github.com/giantswarm/apptest-framework/v4/pkg/state"
	"github.com/giantswarm/apptest-framework/v4/pkg/suite"
	clusterclient "github.com/giantswarm/clustertest/v4/pkg/client"
	"github.com/giantswarm/clustertest/v4/pkg/logger"
	"github.com/giantswarm/clustertest/v4/pkg/wait"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

const (
	isUpgrade            = false
	parentReadyTimeout   = 5 * time.Minute
	childrenReadyTimeout = 10 * time.Minute
	deploymentsTimeout   = 10 * time.Minute
	pollingInterval      = 10 * time.Second
	instanceLabel        = "app.kubernetes.io/instance"
	valuesFile           = "./values.yaml"
)

type resourceKind string

const (
	deploymentKind  resourceKind = "Deployment"
	statefulSetKind resourceKind = "StatefulSet"
	daemonSetKind   resourceKind = "DaemonSet"
)

type appCheck struct {
	appKey   string // key under .apps in values.yaml
	instance string // app.kubernetes.io/instance label value
	kind     resourceKind
}

var appChecks = []appCheck{
	{appKey: "kyverno", instance: "kyverno", kind: deploymentKind},
	{appKey: "kubescape", instance: "kubescape", kind: deploymentKind},
	{appKey: "trivy", instance: "trivy", kind: statefulSetKind},
	{appKey: "trivyOperator", instance: "trivy-operator", kind: deploymentKind},
	{appKey: "starboardExporter", instance: "starboard-exporter", kind: deploymentKind},
	{appKey: "falco", instance: "falco", kind: daemonSetKind},
}

func TestBasic(t *testing.T) {
	enabled := mustReadEnabledApps(t, valuesFile)

	suite.New().
		WithIsUpgrade(isUpgrade).
		WithValuesFile(valuesFile).
		Tests(func() {
			It("should deploy security-bundle and all enabled child Apps", func() {
				ctx := context.Background()
				mc := state.GetFramework().MC()
				cluster := state.GetCluster()
				orgNamespace := cluster.GetNamespace()
				parentName := fmt.Sprintf("%s-security-bundle", cluster.Name)

				By(fmt.Sprintf("Waiting for parent App %s to be deployed", parentName))
				Eventually(wait.IsAppDeployed(ctx, mc, parentName, orgNamespace)).
					WithTimeout(parentReadyTimeout).
					WithPolling(pollingInterval).
					Should(BeTrue())

				By("Listing child Apps managed by the security-bundle")
				appList := &v1alpha1.AppList{}
				err := mc.List(ctx, appList,
					client.InNamespace(orgNamespace),
					client.MatchingLabels{"giantswarm.io/managed-by": parentName},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(appList.Items).NotTo(BeEmpty(), "expected at least one child App managed by %s", parentName)

				children := make([]types.NamespacedName, 0, len(appList.Items))
				for _, app := range appList.Items {
					children = append(children, types.NamespacedName{Name: app.Name, Namespace: app.Namespace})
				}

				By(fmt.Sprintf("Waiting for %d child Apps to be deployed", len(children)))
				Eventually(wait.IsAllAppDeployed(ctx, mc, children)).
					WithTimeout(childrenReadyTimeout).
					WithPolling(pollingInterval).
					Should(BeTrue())
			})

			for _, check := range appChecks {
				if !enabled[check.appKey] {
					continue
				}
				It(fmt.Sprintf("should have all %s %ss running and ready on the workload cluster", check.instance, check.kind), func() {
					ctx := context.Background()
					cluster := state.GetCluster()

					wcClient, err := state.GetFramework().WC(cluster.Name)
					Expect(err).NotTo(HaveOccurred())

					selector := client.MatchingLabels{instanceLabel: check.instance}

					By(fmt.Sprintf("Waiting for all %ss matching %s=%s to be ready", check.kind, instanceLabel, check.instance))
					Eventually(func() error {
						return checkResourceReady(ctx, wcClient, check.instance, check.kind, selector)
					}).WithTimeout(deploymentsTimeout).WithPolling(pollingInterval).Should(Succeed())
				})
			}
		}).
		Run(t, "Basic Test")
}

func checkResourceReady(ctx context.Context, c *clusterclient.Client, app string, kind resourceKind, selector client.MatchingLabels) error {
	switch kind {
	case deploymentKind:
		return deploymentsReady(ctx, c, app, selector)
	case statefulSetKind:
		return statefulSetsReady(ctx, c, app, selector)
	case daemonSetKind:
		return daemonSetsReady(ctx, c, app, selector)
	default:
		return fmt.Errorf("unsupported resource kind %q", kind)
	}
}

func deploymentsReady(ctx context.Context, c *clusterclient.Client, app string, selector client.MatchingLabels) error {
	list := &appsv1.DeploymentList{}
	if err := c.List(ctx, list, selector); err != nil {
		return err
	}
	if len(list.Items) == 0 {
		return fmt.Errorf("no deployments found for app %q matching %v", app, selector)
	}
	for _, d := range list.Items {
		desired := int32(1)
		if d.Spec.Replicas != nil {
			desired = *d.Spec.Replicas
		}
		if d.Status.AvailableReplicas != desired {
			logger.Log("deployment %s/%s from app %s is not yet ready: %d/%d available replicas",
				d.Namespace, d.Name, app, d.Status.AvailableReplicas, desired)
			return fmt.Errorf("deployment %s/%s has %d/%d available replicas",
				d.Namespace, d.Name, d.Status.AvailableReplicas, desired)
		}
		logger.Log("deployment %s/%s from app %s is ready", d.Namespace, d.Name, app)
	}
	return nil
}

func statefulSetsReady(ctx context.Context, c *clusterclient.Client, app string, selector client.MatchingLabels) error {
	list := &appsv1.StatefulSetList{}
	if err := c.List(ctx, list, selector); err != nil {
		return err
	}
	if len(list.Items) == 0 {
		return fmt.Errorf("no statefulsets found for app %q matching %v", app, selector)
	}
	for _, s := range list.Items {
		desired := int32(1)
		if s.Spec.Replicas != nil {
			desired = *s.Spec.Replicas
		}
		if s.Status.AvailableReplicas != desired {
			logger.Log("statefulset %s/%s from app %s is not yet ready: %d/%d available replicas",
				s.Namespace, s.Name, app, s.Status.AvailableReplicas, desired)
			return fmt.Errorf("statefulset %s/%s has %d/%d available replicas",
				s.Namespace, s.Name, s.Status.AvailableReplicas, desired)
		}
		logger.Log("statefulset %s/%s from app %s is ready", s.Namespace, s.Name, app)
	}
	return nil
}

func daemonSetsReady(ctx context.Context, c *clusterclient.Client, app string, selector client.MatchingLabels) error {
	list := &appsv1.DaemonSetList{}
	if err := c.List(ctx, list, selector); err != nil {
		return err
	}
	if len(list.Items) == 0 {
		return fmt.Errorf("no daemonsets found for app %q matching %v", app, selector)
	}
	for _, d := range list.Items {
		if d.Status.NumberReady != d.Status.DesiredNumberScheduled {
			logger.Log("daemonset %s/%s from app %s is not yet ready: %d/%d pods ready",
				d.Namespace, d.Name, app, d.Status.NumberReady, d.Status.DesiredNumberScheduled)
			return fmt.Errorf("daemonset %s/%s has %d/%d ready pods",
				d.Namespace, d.Name, d.Status.NumberReady, d.Status.DesiredNumberScheduled)
		}
		logger.Log("daemonset %s/%s from app %s is ready", d.Namespace, d.Name, app)
	}
	return nil
}

func mustReadEnabledApps(t *testing.T, path string) map[string]bool {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var v struct {
		Apps map[string]struct {
			Enabled bool `json:"enabled"`
		} `json:"apps"`
	}
	if err := yaml.Unmarshal(data, &v); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	enabled := make(map[string]bool, len(v.Apps))
	for k, a := range v.Apps {
		enabled[k] = a.Enabled
	}
	return enabled
}

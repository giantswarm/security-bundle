package basic

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	helmv2 "github.com/fluxcd/helm-controller/api/v2"
	"github.com/giantswarm/apptest-framework/v5/pkg/state"
	"github.com/giantswarm/apptest-framework/v5/pkg/suite"
	clusterclient "github.com/giantswarm/clustertest/v5/pkg/client"
	"github.com/giantswarm/clustertest/v5/pkg/failurehandler"
	"github.com/giantswarm/clustertest/v5/pkg/helmrelease"
	"github.com/giantswarm/clustertest/v5/pkg/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

const (
	isUpgrade          = false
	parentReadyTimeout = 5 * time.Minute
	// kyverno alone has been observed taking ~13 minutes to install on a fresh
	// cluster, so the budget for the whole set has to be well above that.
	childrenReadyTimeout = 25 * time.Minute
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
			It("should deploy security-bundle and all enabled child HelmReleases", func() {
				ctx := context.Background()
				framework := state.GetFramework()
				mc := framework.MC()
				cluster := state.GetCluster()
				orgNamespace := cluster.GetNamespace()
				parentName := fmt.Sprintf("%s-security-bundle", cluster.Name)

				By(fmt.Sprintf("Waiting for the security-bundle %s to be deployed", parentName))
				Eventually(helmrelease.IsAppOrHelmReleaseReady(ctx, mc, parentName, orgNamespace)).
					WithTimeout(parentReadyTimeout).
					WithPolling(pollingInterval).
					Should(BeTrue(), failurehandler.HelmReleasesNotReady(framework, cluster))

				By("Listing child HelmReleases managed by the security-bundle")
				hrList := &helmv2.HelmReleaseList{}
				err := mc.List(ctx, hrList,
					client.InNamespace(orgNamespace),
					client.MatchingLabels{"giantswarm.io/managed-by": parentName},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(hrList.Items).NotTo(BeEmpty(), "expected at least one child HelmRelease managed by %s", parentName)

				children := make([]types.NamespacedName, 0, len(hrList.Items))
				for _, hr := range hrList.Items {
					children = append(children, types.NamespacedName{Name: hr.Name, Namespace: hr.Namespace})
				}

				By(fmt.Sprintf("Waiting for %d child HelmReleases to be ready", len(children)))
				Eventually(helmrelease.AreAllReady(ctx, mc, children)).
					WithTimeout(childrenReadyTimeout).
					WithPolling(pollingInterval).
					Should(Succeed(), failurehandler.HelmReleasesNotReady(framework, cluster))
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

			// Functional scenarios exercising the security apps against a real
			// workload. They share one namespace and one compliant Deployment,
			// so they run Ordered and tear down together in AfterAll rather
			// than per-spec.
			Describe("security app scenarios", Ordered, func() {
				AfterAll(cleanupScenarioResources)

				registerScenarioNamespace()
				registerKyvernoPolicyScenarios()
				registerTrivyScenarios()
				registerPolicyExceptionScenarios()
			})
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

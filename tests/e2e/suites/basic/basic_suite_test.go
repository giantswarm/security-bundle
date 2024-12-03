package basic

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/giantswarm/apptest-framework/pkg/config"
	"github.com/giantswarm/apptest-framework/pkg/state"
	"github.com/giantswarm/apptest-framework/pkg/suite"
	"github.com/giantswarm/clustertest/pkg/wait"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	batchv1 "k8s.io/api/batch/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "k8s.io/api/apps/v1"
)

const (
	isUpgrade        = false
	appReadyTimeout  = 10 * time.Minute
	appReadyInterval = 5 * time.Second
)

var components = []string{
	"security-bundle",
	"exception-recommender",
	"falco",
	"kyverno-crds",
	"kyverno",
	"kyverno-policy-operator",
	"kyverno-policies",
	"starboard-exporter",
	"trivy",
	"trivy-operator",
}

func TestBasic(t *testing.T) {
	suite.New(config.MustLoad("../../config.yaml")).
		WithIsUpgrade(isUpgrade).
		WithValuesFile("./values.yaml").
		Tests(func() {
			It("should deploy all security bundle App CRs successfully", func() {
				Expect(state.GetCluster()).NotTo(BeNil(), "cluster state should be initialized")
				Expect(state.GetCluster().Organization).NotTo(BeNil(), "organization should be available")

				namespace := state.GetCluster().Organization.GetNamespace()

				By("Verifying all App CRs are deployed")
				for _, component := range components {
					appName := fmt.Sprintf("%s-%s", state.GetCluster().Name, component)
					Eventually(wait.IsAppDeployed(context.Background(),
						state.GetFramework().MC(),
						appName,
						namespace)).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(BeTrue(), fmt.Sprintf("%s should be deployed", component))
				}
			})
			It("should have all components running and ready", func() {
				wcClient, err := state.GetFramework().WC(state.GetCluster().Name)
				Expect(err).NotTo(HaveOccurred(), "should get workload cluster client")

				componentConfigs := map[string]struct {
					namespace string
					kind      string
					name      string
				}{
					// security-bundle namespace
					"exception-recommender":   {namespace: "security-bundle", kind: "Deployment", name: "exception-recommender"},
					"falco":                   {namespace: "security-bundle", kind: "DaemonSet", name: "falco"},
					"falco-exporter":          {namespace: "security-bundle", kind: "DaemonSet", name: "falco-falco-exporter"},
					"falco-sidekick":          {namespace: "security-bundle", kind: "Deployment", name: "falco-falcosidekick"},
					"falco-metacollector":     {namespace: "security-bundle", kind: "Deployment", name: "falco-k8s-metacollector"},
					"kyverno-policy-operator": {namespace: "security-bundle", kind: "Deployment", name: "kyverno-policy-operator"},
					"starboard-exporter":      {namespace: "security-bundle", kind: "Deployment", name: "starboard-exporter"},
					"trivy":                   {namespace: "security-bundle", kind: "StatefulSet", name: "trivy"},
					"trivy-operator":          {namespace: "security-bundle", kind: "Deployment", name: "trivy-operator"},

					// kyverno namespace
					"kyverno":                       {namespace: "kyverno", kind: "Deployment", name: "kyverno-admission-controller"},
					"kyverno-background-controller": {namespace: "kyverno", kind: "Deployment", name: "kyverno-background-controller"},
					"kyverno-cleanup-controller":    {namespace: "kyverno", kind: "Deployment", name: "kyverno-cleanup-controller"},
					"kyverno-kyverno-plugin":        {namespace: "kyverno", kind: "Deployment", name: "kyverno-kyverno-plugin"},
					"kyverno-policy-reporter":       {namespace: "kyverno", kind: "Deployment", name: "kyverno-policy-reporter"},
					"kyverno-reports-controller":    {namespace: "kyverno", kind: "Deployment", name: "kyverno-reports-controller"},
					"kyverno-ui":                    {namespace: "kyverno", kind: "Deployment", name: "kyverno-ui"},

					// kyverno-webhook
					"kyverno-webhook": {namespace: "kyverno", kind: "MutatingWebhookConfiguration", name: "kyverno-policy-mutating-webhook-cfg"},

					// jobs
					"kyverno-cleanup-cluster-ephemeral-reports": {namespace: "kyverno", kind: "CronJob", name: "kyverno-cleanup-cluster-ephemeral-reports"},
					"kyverno-cleanup-ephemeral-reports":         {namespace: "kyverno", kind: "CronJob", name: "kyverno-cleanup-ephemeral-reports"},
					"kyverno-cleanup-update-requests":           {namespace: "kyverno", kind: "CronJob", name: "kyverno-cleanup-update-requests"},
				}

				for component, config := range componentConfigs {
					By(fmt.Sprintf("Checking %s %s", component, config.kind))
					Eventually(func() bool {
						var ready, replicas int32
						switch config.kind {
						case "Deployment":
							deployment := &appsv1.Deployment{}
							err := wcClient.Get(context.Background(), client.ObjectKey{Namespace: config.namespace, Name: config.name}, deployment)
							if err != nil {
								return false
							}
							ready = deployment.Status.ReadyReplicas
							replicas = deployment.Status.Replicas
						case "DaemonSet":
							ds := &appsv1.DaemonSet{}
							err := wcClient.Get(context.Background(), client.ObjectKey{Namespace: config.namespace, Name: config.name}, ds)
							if err != nil {
								return false
							}
							ready = ds.Status.NumberReady
							replicas = ds.Status.DesiredNumberScheduled
						case "StatefulSet":
							sts := &appsv1.StatefulSet{}
							err := wcClient.Get(context.Background(), client.ObjectKey{Namespace: config.namespace, Name: config.name}, sts)
							if err != nil {
								return false
							}
							ready = sts.Status.ReadyReplicas
							replicas = sts.Status.Replicas
						case "MutatingWebhookConfiguration":
							webhook := &admissionregistrationv1.MutatingWebhookConfiguration{}
							err := wcClient.Get(context.Background(), client.ObjectKey{Name: config.name}, webhook)
							return err == nil
						case "CronJob":
							cj := &batchv1.CronJob{}
							err := wcClient.Get(context.Background(), client.ObjectKey{Namespace: config.namespace, Name: config.name}, cj)
							return err == nil
						}
						return ready == replicas && replicas > 0

					}).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(BeTrue(), fmt.Sprintf("%s %s should be ready", component, config.kind))
				}
			})
		}).
		Run(t, "Basic Test")
}

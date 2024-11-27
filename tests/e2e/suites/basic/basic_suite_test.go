package basic

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/giantswarm/apptest-framework/pkg/config"
	"github.com/giantswarm/apptest-framework/pkg/state"
	"github.com/giantswarm/apptest-framework/pkg/suite"
	"github.com/giantswarm/clustertest/pkg/wait"

	appsv1 "k8s.io/api/apps/v1"
)

const (
	isUpgrade        = false
	appReadyTimeout  = 10 * time.Minute
	appReadyInterval = 5 * time.Second
)

var components = []string{
	"security-bundle", // Main bundle app
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
			It("should deploy all security bundle components successfully", func() {
				Expect(state.GetCluster()).NotTo(BeNil(), "cluster state should be initialized")
				Expect(state.GetCluster().Organization).NotTo(BeNil(), "organization should be available")

				namespace := state.GetCluster().Organization.GetNamespace()

				By("Verifying all components are deployed")
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

				// Verify workload cluster components
				for _, component := range []string{"kyverno", "trivy", "falco"} {
					By(fmt.Sprintf("Checking %s deployments", component))
					Eventually(func() bool {
						// List deployments in the component's namespace
						deploymentList := &appsv1.DeploymentList{}
						err := wcClient.List(context.Background(), deploymentList, client.InNamespace("namespace-name"))
						if err != nil {
							return false
						}
						// Find the component's deployment and check its status
						for _, d := range deploymentList.Items {
							if d.Name == component {
								return d.Status.ReadyReplicas == d.Status.Replicas && d.Status.Replicas > 0
							}
						}
						return false
					}).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(BeTrue(), fmt.Sprintf("%s deployment should be ready", component))
				}
			})
		}).
		Run(t, "Basic Test")
}

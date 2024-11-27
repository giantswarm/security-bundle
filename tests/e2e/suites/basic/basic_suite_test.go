package basic

import (
	"context"
	"fmt"
	"strings"
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

				componentConfigs := map[string]struct {
					namespace string
					kind      string
					prefix    string
				}{
					"kyverno":            {namespace: "kyverno", kind: "Deployment", prefix: "kyverno"},
					"falco":              {namespace: "giantswarm", kind: "Deployment", prefix: "falco"},
					"trivy":              {namespace: "giantswarm", kind: "Deployment", prefix: "trivy-operator"},
					"trivy-statefulset":  {namespace: "giantswarm", kind: "StatefulSet", prefix: "trivy"},
					"starboard-exporter": {namespace: "giantswarm", kind: "Deployment", prefix: "starboard-exporter"},
				}

				for component, config := range componentConfigs {
					By(fmt.Sprintf("Checking %s %s", component, config.kind))
					Eventually(func() bool {
						var ready, replicas int32
						switch config.kind {
						case "Deployment":
							deploymentList := &appsv1.DeploymentList{}
							err := wcClient.List(context.Background(), deploymentList, client.InNamespace(config.namespace))
							if err != nil {
								return false
							}
							for _, d := range deploymentList.Items {
								if strings.HasPrefix(d.Name, config.prefix) {
									ready = d.Status.ReadyReplicas
									replicas = d.Status.Replicas
									break
								}
							}
						case "StatefulSet":
							stsList := &appsv1.StatefulSetList{}
							err := wcClient.List(context.Background(), stsList, client.InNamespace(config.namespace))
							if err != nil {
								return false
							}
							for _, s := range stsList.Items {
								if strings.HasPrefix(s.Name, config.prefix) {
									ready = s.Status.ReadyReplicas
									replicas = s.Status.Replicas
									break
								}
							}
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

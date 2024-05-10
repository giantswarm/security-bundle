package basic

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/giantswarm/apptest-framework/pkg/config"
	"github.com/giantswarm/apptest-framework/pkg/state"
	"github.com/giantswarm/apptest-framework/pkg/suite"
	"github.com/giantswarm/clustertest/pkg/client"
	"github.com/giantswarm/clustertest/pkg/organization"
	"github.com/giantswarm/clustertest/pkg/wait"
)

const (
	isUpgrade = false
)

func TestBasic(t *testing.T) {
	const (
		timeout          = time.Second * 10
		duration         = time.Second * 10
		interval         = time.Millisecond * 250
		appReadyTimeout  = 10 * time.Minute
		appReadyInterval = 5 * time.Second
		kyvernoNamespace = "kyverno"
		bundleNamespace  = "security-bundle"
	)

	suite.New(config.MustLoad("../../config.yaml")).
		// The namespace to install the app into within the workload cluster
		WithInstallNamespace(state.GetCluster().Organization.GetNamespace()).
		// If this is an upgrade test or not.
		// If true, the suite will first install the latest released version of the app before upgrading to the test version
		WithIsUpgrade(isUpgrade).
		WithValuesFile("./values.yaml").
		BeforeUpgrade(func() {
			// Perform any checks between installing the latest released version
			// and upgrading it to the version to test
			// E.g. ensure that the initial install has completed and has settled before upgrading
		}).
		Tests(func() {
			var org *organization.Org

			BeforeSuite(func() {
				org = state.GetCluster().Organization
			})

			Describe("Check Apps status", func() {

				It("should have kyverno, kyverno-policies and kyverno-policy-operator deplyoed", func() {
					Eventually(wait.IsAppDeployed(state.GetContext(), state.GetFramework().MC(), fmt.Sprintf("%s-kyverno", state.GetCluster().Name), org.GetNamespace())).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(BeTrue())

					Eventually(wait.IsAppDeployed(state.GetContext(), state.GetFramework().MC(), fmt.Sprintf("%s-kyverno-policies", state.GetCluster().Name), org.GetNamespace())).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(BeTrue())

					Eventually(wait.IsAppDeployed(state.GetContext(), state.GetFramework().MC(), fmt.Sprintf("%s-kyverno-policy-operator", state.GetCluster().Name), org.GetNamespace())).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(BeTrue())
				})

				It("should have trivy, trivy-operator and starboard-exporter deplyoed", func() {
					Eventually(wait.IsAppDeployed(state.GetContext(), state.GetFramework().MC(), fmt.Sprintf("%s-trivy", state.GetCluster().Name), org.GetNamespace())).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(BeTrue())

					Eventually(wait.IsAppDeployed(state.GetContext(), state.GetFramework().MC(), fmt.Sprintf("%s-trivy-operator", state.GetCluster().Name), org.GetNamespace())).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(BeTrue())

					Eventually(wait.IsAppDeployed(state.GetContext(), state.GetFramework().MC(), fmt.Sprintf("%s-starboard-exporter", state.GetCluster().Name), org.GetNamespace())).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(BeTrue())
				})

				It("should have falco deplyoed", func() {
					Eventually(wait.IsAppDeployed(state.GetContext(), state.GetFramework().MC(), fmt.Sprintf("%s-falco", state.GetCluster().Name), org.GetNamespace())).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(BeTrue())
				})
			})

			Describe("Check that Apps are running", func() {
				var checkDeployment v1.Deployment
				var wcClient *client.Client
				var err error
				var ctx context.Context

				BeforeEach(func() {
					checkDeployment = v1.Deployment{}
				})

				BeforeAll(func() {
					wcClient, err = state.GetFramework().WC(state.GetCluster().Name)
					Expect(err).ToNot(HaveOccurred())

					ctx = state.GetContext()
				})

				It("should have kyverno-admission running", func() {
					kyvernoAdmissionDeploymentName := "kyverno-admission-controller"

					kyvernoDeploymentLookup := types.NamespacedName{Name: kyvernoAdmissionDeploymentName, Namespace: kyvernoNamespace}

					By("checking if the kyverno-admission-controller Deployment is satisfied")

					Eventually(func() bool {
						err := wcClient.Get(ctx, kyvernoDeploymentLookup, &checkDeployment)
						if err != nil {
							fmt.Printf("Unable to get %s Deployment", kyvernoAdmissionDeploymentName)
							return false
						}

						if checkDeployment.Status.ReadyReplicas >= 3 {
							fmt.Printf("%s Deployment has %d replicas ready", kyvernoAdmissionDeploymentName, checkDeployment.Status.ReadyReplicas)
							return true
						} else {
							fmt.Printf("%s Deployment is not yet satisfied: Has %d replicas ready", kyvernoAdmissionDeploymentName, checkDeployment.Status.ReadyReplicas)
							return false
						}
					}).
						WithTimeout(timeout).
						WithPolling(interval).
						Should(BeTrue())
				})
			})
		}).
		Run(t, "Basic Test")
}

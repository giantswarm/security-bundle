package basic

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/giantswarm/apptest-framework/pkg/config"
	"github.com/giantswarm/apptest-framework/pkg/state"
	"github.com/giantswarm/apptest-framework/pkg/suite"
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
		WithInstallNamespace("default").
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
			Describe("Check Apps status", func() {
				var org *organization.Org

				mcClient := state.GetFramework().MC()
				clusterName := state.GetCluster().Name
				ctx := state.GetContext()

				BeforeEach(func() {
					org = state.GetCluster().Organization
				})

				It("should have kyverno, kyverno-policies and kyverno-policy-operator deplyoed", func() {
					Eventually(wait.IsAppDeployed(ctx, mcClient, fmt.Sprintf("%s-kyverno", clusterName), org.GetNamespace())).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(BeTrue())

					Eventually(wait.IsAppDeployed(ctx, mcClient, fmt.Sprintf("%s-kyverno-policies", clusterName), org.GetNamespace())).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(BeTrue())

					Eventually(wait.IsAppDeployed(ctx, mcClient, fmt.Sprintf("%s-kyverno-policy-operator", clusterName), org.GetNamespace())).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(BeTrue())
				})

				It("should have trivy, trivy-operator and starboard-exporter deplyoed", func() {
					Eventually(wait.IsAppDeployed(ctx, mcClient, fmt.Sprintf("%s-trivy", clusterName), org.GetNamespace())).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(BeTrue())

					Eventually(wait.IsAppDeployed(ctx, mcClient, fmt.Sprintf("%s-trivy-operator", clusterName), org.GetNamespace())).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(BeTrue())

					Eventually(wait.IsAppDeployed(ctx, mcClient, fmt.Sprintf("%s-starboard-exporter", clusterName), org.GetNamespace())).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(BeTrue())
				})

				It("should have falco deplyoed", func() {
					Eventually(wait.IsAppDeployed(ctx, mcClient, fmt.Sprintf("%s-falco", clusterName), org.GetNamespace())).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(BeTrue())
				})
			})
			// Describe("Check that Apps are running", func() {
			// 	It("should have kyverno-admission running", func() {
			// 		kyvernoAdmissionDeploymentName := "kyverno-admission-controller"
			// 		kyvernoAdmissionDeployment := v1.Deployment{}

			// 		kyvernoDeploymentLookup := types.NamespacedName{Name: kyvernoAdmissionDeploymentName, Namespace: kyvernoNamespace}
			// 		clusterName := state.GetCluster().Name

			// 		By("checking if the kyverno-admission-controller Deployment is satisfied")

			// 		Eventually(func() bool {
			// 			wcClient, err := state.GetFramework().WC(clusterName)
			// 			if err != nil {
			// 				fmt.Println("Unable to get Workload Cluster client")
			// 				return false
			// 			}
			// 			err = wcClient.Get(state.GetContext(), kyvernoDeploymentLookup, &kyvernoAdmissionDeployment)
			// 			if err != nil {
			// 				fmt.Printf("Unable to get %s Deployment", kyvernoAdmissionDeploymentName)
			// 				return false
			// 			}

			// 			if kyvernoAdmissionDeployment.Status.ReadyReplicas >= 3 {
			// 				fmt.Printf("%s Deployment has %d replicas ready", kyvernoAdmissionDeploymentName, kyvernoAdmissionDeployment.Status.ReadyReplicas)
			// 				return true
			// 			} else {
			// 				fmt.Printf("%s Deployment is not yet satisfied: Has %d replicas ready", kyvernoAdmissionDeploymentName, kyvernoAdmissionDeployment.Status.ReadyReplicas)
			// 				return false
			// 			}
			// 		}).
			// 			WithTimeout(timeout).
			// 			WithPolling(interval).
			// 			Should(BeTrue())
			// 	})
			// })
		}).
		Run(t, "Basic Test")
}

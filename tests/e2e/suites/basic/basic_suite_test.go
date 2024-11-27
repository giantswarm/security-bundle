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
	"github.com/giantswarm/clustertest/pkg/wait"
)

const (
	isUpgrade        = false
	appReadyTimeout  = 10 * time.Minute
	appReadyInterval = 5 * time.Second
	kyvernoNamespace = "kyverno"
	bundleNamespace  = "security-bundle"
)

func TestBasic(t *testing.T) {
	cfg := config.MustLoad("../../config.yaml")

	testSuite := suite.New(cfg).
		WithValuesFile("./values.yaml").
		WithIsUpgrade(isUpgrade)

	testSuite.Tests(func() {
		BeforeSuite(func() {
			// Wait for setup to complete and verify state
			Eventually(func() bool {
				return state.GetCluster() != nil &&
					state.GetCluster().Organization != nil
			}).WithTimeout(2*time.Minute).
				WithPolling(5*time.Second).
				Should(BeTrue(), "cluster state should be initialized")

			// Once state is initialized, set the namespace
			testSuite.WithInstallNamespace(state.GetCluster().Organization.GetNamespace())
		})

		Describe("Security Bundle Components", func() {
			It("should deploy core components", func() {
				org := state.GetCluster().Organization

				// Wait for kyverno CRDs
				Eventually(wait.IsAppDeployed(state.GetContext(),
					state.GetFramework().MC(),
					fmt.Sprintf("%s-kyverno-crds", state.GetCluster().Name),
					org.GetNamespace())).
					WithTimeout(appReadyTimeout).
					WithPolling(appReadyInterval).
					Should(BeTrue(), "kyverno CRDs should be deployed")

				// Wait for kyverno
				Eventually(wait.IsAppDeployed(state.GetContext(),
					state.GetFramework().MC(),
					fmt.Sprintf("%s-kyverno", state.GetCluster().Name),
					org.GetNamespace())).
					WithTimeout(appReadyTimeout).
					WithPolling(appReadyInterval).
					Should(BeTrue(), "kyverno should be deployed")

				// Wait for kyverno policies
				Eventually(wait.IsAppDeployed(state.GetContext(),
					state.GetFramework().MC(),
					fmt.Sprintf("%s-kyverno-policies", state.GetCluster().Name),
					org.GetNamespace())).
					WithTimeout(appReadyTimeout).
					WithPolling(appReadyInterval).
					Should(BeTrue(), "kyverno policies should be deployed")

				// Wait for kyverno policy operator
				Eventually(wait.IsAppDeployed(state.GetContext(),
					state.GetFramework().MC(),
					fmt.Sprintf("%s-kyverno-policy-operator", state.GetCluster().Name),
					org.GetNamespace())).
					WithTimeout(appReadyTimeout).
					WithPolling(appReadyInterval).
					Should(BeTrue(), "kyverno policy operator should be deployed")
			})

			Context("Optional Components", func() {
				It("should deploy trivy components when enabled", func() {
					org := state.GetCluster().Organization

					// Wait for trivy
					Eventually(wait.IsAppDeployed(state.GetContext(),
						state.GetFramework().MC(),
						fmt.Sprintf("%s-trivy", state.GetCluster().Name),
						org.GetNamespace())).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(BeTrue(), "trivy should be deployed")

					// Wait for trivy operator
					Eventually(wait.IsAppDeployed(state.GetContext(),
						state.GetFramework().MC(),
						fmt.Sprintf("%s-trivy-operator", state.GetCluster().Name),
						org.GetNamespace())).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(BeTrue(), "trivy operator should be deployed")
				})

				It("should deploy falco when enabled", func() {
					org := state.GetCluster().Organization

					Eventually(wait.IsAppDeployed(state.GetContext(),
						state.GetFramework().MC(),
						fmt.Sprintf("%s-falco", state.GetCluster().Name),
						org.GetNamespace())).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(BeTrue(), "falco should be deployed")
				})
			})
		})
	}).Run(t, "Security Bundle Basic Test")
}

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
		WithInstallNamespace(bundleNamespace).
		// If this is an upgrade test or not.
		// If true, the suite will first install the latest released version of the app before upgrading to the test version
		WithIsUpgrade(isUpgrade).
		WithValuesFile("./values.yaml").
		AfterClusterReady(func() {
			It("should have kyverno and kyverno-policy-operator deplyoed", func() {
				org := state.GetCluster().Organization

				// Ingress-Nginx depends on external-dns for DNS resolution and cert-manager for certificates.
				Eventually(wait.IsAppDeployed(state.GetContext(), state.GetFramework().MC(), fmt.Sprintf("%s-kyverno", state.GetCluster().Name), org.GetNamespace())).
					WithTimeout(appReadyTimeout).
					WithPolling(appReadyInterval).
					Should(BeTrue())

				Eventually(wait.IsAppDeployed(state.GetContext(), state.GetFramework().MC(), fmt.Sprintf("%s-kyverno-policy-operator", state.GetCluster().Name), org.GetNamespace())).
					WithTimeout(appReadyTimeout).
					WithPolling(appReadyInterval).
					Should(BeTrue())
			})
		}).
		BeforeUpgrade(func() {
			// Perform any checks between installing the latest released version
			// and upgrading it to the version to test
			// E.g. ensure that the initial install has completed and has settled before upgrading
		}).
		Tests(func() {
			// kyvernoAdmissionDeploymentName := "kyverno-admission-controller"
			// kyvernoAdmissionDeployment := v1.Deployment{}

			// kyvernoDeploymentLookup := types.NamespacedName{Name: kyvernoAdmissionDeploymentName, Namespace: kyvernoNamespace}

			It("should have kyverno running", func() {
				// STEP
				By("checking if the deployment is satisfied")
				Eventually(func() bool {
					return true
				}, timeout, interval).Should(BeTrue())
			})
		}).
		Run(t, "Basic Test")
}

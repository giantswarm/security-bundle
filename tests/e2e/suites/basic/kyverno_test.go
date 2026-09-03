package basic

import (
	"context"
	"fmt"
	"time"

	"github.com/giantswarm/clustertest/v5/pkg/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const workloadReadyTimeout = 5 * time.Minute

// registerKyvernoPolicyScenarios covers the restricted Pod Security Standards
// shipped by kyverno-policies. The bundle sets `validationFailureAction:
// Enforce` for them via userConfig, so a violating workload is rejected at
// admission rather than merely reported.
func registerKyvernoPolicyScenarios() {
	It("should reject a Deployment whose pods run as root", func() {
		ctx := context.Background()
		wc := wcClient()

		By(fmt.Sprintf("Creating Deployment %s, which runs as root", violatingWorkload))
		err := wc.Create(ctx, workloadDeployment(violatingWorkload, false))

		Expect(err).To(HaveOccurred(), "expected kyverno to reject a Deployment running as root")
		// Kyverno's autogen rules match the Deployment directly, so the
		// rejection lands on this create rather than on a downstream pod.
		for _, policy := range violatedPolicies {
			Expect(err.Error()).To(ContainSubstring(policy),
				"expected the admission error to name the %s policy", policy)
		}

		logger.Log("Deployment %s was rejected: %s", violatingWorkload, err)
	})

	It("should admit a Deployment that complies with the restricted Pod Security Standards", func() {
		ctx := context.Background()
		wc := wcClient()

		By(fmt.Sprintf("Creating Deployment %s", compliantWorkload))
		deployment := workloadDeployment(compliantWorkload, true)
		Expect(wc.Create(ctx, deployment)).To(Succeed())

		DeferCleanup(func() {
			Expect(deleteIfExists(context.Background(), wcClient(), deployment)).To(Succeed())
		})

		By("Waiting for the workload to become available")
		Eventually(func() error {
			return deploymentReady(ctx, wc, compliantWorkload)
		}).WithTimeout(workloadReadyTimeout).WithPolling(pollingInterval).Should(Succeed())
	})
}

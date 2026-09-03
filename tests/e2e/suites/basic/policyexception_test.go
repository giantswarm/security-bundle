package basic

import (
	"context"
	"fmt"
	"time"

	"github.com/giantswarm/clustertest/v5/pkg/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	policyExceptionName    = "e2e-runs-as-root-exception"
	policyExceptionTimeout = 3 * time.Minute

	// exceptionEffectiveTimeout covers kyverno refreshing its exception cache
	// after the operator creates the Kyverno PolicyException.
	exceptionEffectiveTimeout = 3 * time.Minute

	managedByLabel = "app.kubernetes.io/managed-by"
	managedByValue = "kyverno-policy-operator"
)

// registerPolicyExceptionScenarios covers kyverno-policy-operator translating a
// Giant Swarm PolicyException into a Kyverno one, and then confirms the
// translation actually exempts the workload rather than only producing a CR.
func registerPolicyExceptionScenarios() {
	It("should translate a Giant Swarm PolicyException into a Kyverno PolicyException", func() {
		ctx := context.Background()
		wc := wcClient()

		// The operator resolves each referenced policy from its cache of
		// installed ClusterPolicies and requeues indefinitely on an unknown
		// name, so these must be policies kyverno-policies actually ships.
		gsException := newUnstructured(gsPolicyExceptionGVK)
		gsException.SetName(policyExceptionName)
		gsException.SetNamespace(policyExceptionNamespace)
		Expect(unstructured.SetNestedStringSlice(gsException.Object, violatedPolicies, "spec", "policies")).To(Succeed())
		Expect(unstructured.SetNestedSlice(gsException.Object, []interface{}{
			map[string]interface{}{
				"kind":       "Deployment",
				"names":      []interface{}{violatingWorkload},
				"namespaces": []interface{}{scenarioNamespace},
			},
		}, "spec", "targets")).To(Succeed())

		By(fmt.Sprintf("Creating Giant Swarm PolicyException %s/%s", policyExceptionNamespace, policyExceptionName))
		Expect(wc.Create(ctx, gsException)).To(Succeed())

		DeferCleanup(func() {
			// The Kyverno PolicyException is owned by this one, so deleting it
			// takes both away.
			Expect(deleteIfExists(context.Background(), wcClient(), gsException)).To(Succeed())
		})

		By("Waiting for the operator to create the corresponding Kyverno PolicyException")
		kyvernoException := newUnstructured(kyvernoPolicyExceptionGVK)
		Eventually(func() error {
			return wc.Get(ctx, client.ObjectKey{
				Namespace: policyExceptionNamespace,
				Name:      policyExceptionName,
			}, kyvernoException)
		}).WithTimeout(policyExceptionTimeout).WithPolling(pollingInterval).Should(Succeed())

		By("Checking the generated PolicyException")
		Expect(kyvernoException.GetLabels()).To(HaveKeyWithValue(managedByLabel, managedByValue))

		owners := kyvernoException.GetOwnerReferences()
		Expect(owners).To(HaveLen(1), "expected the Kyverno PolicyException to be owned by its Giant Swarm counterpart")
		Expect(owners[0].Kind).To(Equal(gsPolicyExceptionGVK.Kind))
		Expect(owners[0].Name).To(Equal(policyExceptionName))

		Expect(exceptedPolicies(kyvernoException)).To(ConsistOf(violatedPolicies),
			"expected every requested policy to be carried over into spec.exceptions")

		logger.Log("Kyverno PolicyException %s/%s excepts %v",
			policyExceptionNamespace, policyExceptionName, violatedPolicies)
	})

	It("should admit the previously rejected Deployment once the exception exists", func() {
		ctx := context.Background()
		wc := wcClient()

		deployment := workloadDeployment(violatingWorkload, false)

		DeferCleanup(func() {
			Expect(deleteIfExists(context.Background(), wcClient(), deployment)).To(Succeed())
		})

		// Retried because kyverno may not have picked the new exception up yet;
		// a rejected create leaves nothing behind, so retrying is safe.
		By(fmt.Sprintf("Recreating Deployment %s, now covered by the exception", violatingWorkload))
		Eventually(func() error {
			return wc.Create(ctx, deployment)
		}).WithTimeout(exceptionEffectiveTimeout).WithPolling(pollingInterval).Should(Succeed())

		By("Waiting for the exempted workload to become available")
		Eventually(func() error {
			return deploymentReady(ctx, wc, violatingWorkload)
		}).WithTimeout(workloadReadyTimeout).WithPolling(pollingInterval).Should(Succeed())
	})
}

// exceptedPolicies returns the policy names listed in spec.exceptions.
func exceptedPolicies(exception *unstructured.Unstructured) []string {
	entries, found, err := unstructured.NestedSlice(exception.Object, "spec", "exceptions")
	if err != nil || !found {
		return nil
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		exceptionEntry, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		if name, ok := exceptionEntry["policyName"].(string); ok {
			names = append(names, name)
		}
	}
	return names
}

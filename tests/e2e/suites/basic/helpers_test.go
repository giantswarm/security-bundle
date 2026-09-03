package basic

import (
	"context"
	"fmt"

	"github.com/giantswarm/apptest-framework/v5/pkg/state"
	clusterclient "github.com/giantswarm/clustertest/v5/pkg/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Scenario resources all live in a single throwaway namespace on the workload
// cluster, except the PolicyExceptions which the operator pins to its own
// destination namespace.
const (
	scenarioNamespace = "security-bundle-e2e"

	// policyExceptionNamespace is the `destinationNamespace` the
	// kyverno-policy-operator chart configures, so the generated Kyverno
	// PolicyException always lands here regardless of where its Giant Swarm
	// counterpart is created.
	policyExceptionNamespace = "policy-exceptions"

	// scanTarget is pinned rather than floating on :latest so the
	// disallow-latest-tag best-practices policy stays quiet, and busybox ships
	// the wget used to scrape the starboard-exporter metrics endpoint.
	scanTargetImage      = "gsoci.azurecr.io/giantswarm/busybox:1.38.0"
	scanTargetRepository = "giantswarm/busybox"

	// compliantWorkload doubles as the kyverno "should be admitted" case, the
	// image trivy-operator scans, and the pod the metrics scrape execs from.
	compliantWorkload = "e2e-compliant"
	violatingWorkload = "e2e-runs-as-root"

	workloadContainer = "test"
	workloadLabel     = "e2e-workload"
)

// The restricted Pod Security Standards policies the violating workload breaks.
// It is otherwise identical to the compliant one, so this list is exactly what
// the PolicyException has to cover for the workload to be admitted.
var violatedPolicies = []string{
	"require-run-as-nonroot",
	"require-run-as-non-root-user",
}

// CRDs are addressed as unstructured to keep the kyverno, trivy-operator and
// policy-api modules out of this module's dependency graph.
var (
	vulnerabilityReportListGVK = schema.GroupVersionKind{
		Group:   "aquasecurity.github.io",
		Version: "v1alpha1",
		Kind:    "VulnerabilityReportList",
	}
	gsPolicyExceptionGVK = schema.GroupVersionKind{
		Group:   "policy.giantswarm.io",
		Version: "v1alpha1",
		Kind:    "PolicyException",
	}
	kyvernoPolicyExceptionGVK = schema.GroupVersionKind{
		Group:   "kyverno.io",
		Version: "v2",
		Kind:    "PolicyException",
	}
)

// registerScenarioNamespace creates the namespace the scenarios below share and
// tears it down, along with everything in it, once the suite finishes.
func registerScenarioNamespace() {
	It("should create the scenario namespace on the workload cluster", func() {
		ctx := context.Background()
		wc := wcClient()

		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: scenarioNamespace},
		}

		By(fmt.Sprintf("Creating namespace %s", scenarioNamespace))
		err := wc.Create(ctx, namespace)
		if apierrors.IsAlreadyExists(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())

		DeferCleanup(func() {
			Expect(deleteIfExists(context.Background(), wcClient(), namespace)).To(Succeed())
		})
	})
}

// wcClient returns a client for the workload cluster under test.
func wcClient() *clusterclient.Client {
	c, err := state.GetFramework().WC(state.GetCluster().Name)
	Expect(err).NotTo(HaveOccurred())
	return c
}

// workloadDeployment builds a single-replica Deployment running the scan target
// image. When compliant is false the pod runs as root, which is the only
// difference between the two - every other restricted PSS control stays
// satisfied so the failure is attributable to the run-as-root policies alone.
func workloadDeployment(name string, compliant bool) *appsv1.Deployment {
	replicas := int32(1)
	allowPrivilegeEscalation := false
	runAsNonRoot := compliant
	user := int64(35)
	if !compliant {
		user = 0
	}

	selector := map[string]string{workloadLabel: name}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: scenarioNamespace,
			Labels:    selector,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: selector},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: selector},
				Spec: corev1.PodSpec{
					SecurityContext: &corev1.PodSecurityContext{
						RunAsNonRoot: &runAsNonRoot,
						RunAsUser:    &user,
						RunAsGroup:   &user,
						SeccompProfile: &corev1.SeccompProfile{
							Type: corev1.SeccompProfileTypeRuntimeDefault,
						},
					},
					Containers: []corev1.Container{
						{
							Name:  workloadContainer,
							Image: scanTargetImage,
							// busybox sleep only understands seconds.
							Command: []string{"sleep", "99999999"},
							SecurityContext: &corev1.SecurityContext{
								AllowPrivilegeEscalation: &allowPrivilegeEscalation,
								Capabilities: &corev1.Capabilities{
									Drop: []corev1.Capability{"ALL"},
								},
							},
						},
					},
				},
			},
		},
	}
}

// deploymentReady reports whether every replica of the named Deployment is available.
func deploymentReady(ctx context.Context, c *clusterclient.Client, name string) error {
	deployment := &appsv1.Deployment{}
	if err := c.Get(ctx, client.ObjectKey{Namespace: scenarioNamespace, Name: name}, deployment); err != nil {
		return err
	}

	desired := int32(1)
	if deployment.Spec.Replicas != nil {
		desired = *deployment.Spec.Replicas
	}
	if deployment.Status.AvailableReplicas != desired {
		return fmt.Errorf("deployment %s/%s has %d/%d available replicas",
			scenarioNamespace, name, deployment.Status.AvailableReplicas, desired)
	}
	return nil
}

// runningPod returns the name of a running pod belonging to the named workload.
func runningPod(ctx context.Context, c *clusterclient.Client, workload string) (string, error) {
	pods := &corev1.PodList{}
	err := c.List(ctx, pods,
		client.InNamespace(scenarioNamespace),
		client.MatchingLabels{workloadLabel: workload},
	)
	if err != nil {
		return "", err
	}

	for _, pod := range pods.Items {
		if pod.Status.Phase == corev1.PodRunning {
			return pod.Name, nil
		}
	}
	return "", fmt.Errorf("no running pod found for workload %q in namespace %s", workload, scenarioNamespace)
}

// newUnstructured returns an empty object stamped with the given kind.
func newUnstructured(gvk schema.GroupVersionKind) *unstructured.Unstructured {
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(gvk)
	return obj
}

// deleteIfExists removes an object, tolerating one that was never created or
// was already swept up by an owner reference.
func deleteIfExists(ctx context.Context, c *clusterclient.Client, obj client.Object) error {
	err := c.Delete(ctx, obj)
	if apierrors.IsNotFound(err) {
		return nil
	}
	return err
}

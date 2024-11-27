package basic

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/giantswarm/apptest-framework/pkg/config"
	"github.com/giantswarm/apptest-framework/pkg/state"
	"github.com/giantswarm/apptest-framework/pkg/suite"
	"github.com/giantswarm/clustertest/pkg/wait"

	helmv2beta1 "github.com/fluxcd/helm-controller/api/v2beta1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	isUpgrade        = false
	appReadyTimeout  = 5 * time.Minute
	appReadyInterval = 5 * time.Second
	testNamespace    = "security-bundle-test"
)

func TestBasic(t *testing.T) {
	if os.Getenv("E2E_KUBECONFIG") == "" {
		t.Fatal("E2E_KUBECONFIG environment variable must be set")
	}

	suite.New(config.MustLoad("../../config.yaml")).
		WithIsUpgrade(isUpgrade).
		WithValuesFile("./values.yaml").
		WithInstallNamespace(testNamespace).
		AfterClusterReady(func() {
			// Create test namespace if it doesn't exist
			ctx := context.Background()
			if state.GetFramework() != nil && state.GetFramework().MC() != nil {
				ns := &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: testNamespace,
					},
				}
				err := state.GetFramework().MC().Create(ctx, ns)
				if err != nil && !k8serrors.IsAlreadyExists(err) {
					Fail(fmt.Sprintf("Failed to create test namespace: %v", err))
				}
			}
		}).
		Tests(func() {
			It("should deploy kyverno core components", func() {
				components := []string{"kyverno-crds", "kyverno", "kyverno-policies", "kyverno-policy-operator"}

				for _, component := range components {
					appName := fmt.Sprintf("security-bundle-%s", component)

					By(fmt.Sprintf("Checking HelmRelease for %s", component))
					Eventually(func() (bool, error) {
						mcKubeClient := state.GetFramework().MC()
						release := &helmv2beta1.HelmRelease{}
						err := mcKubeClient.Get(context.Background(), types.NamespacedName{
							Name:      appName,
							Namespace: testNamespace,
						}, release)
						if err != nil {
							return false, err
						}

						for _, c := range release.Status.Conditions {
							if c.Type == "Ready" {
								if c.Status == "True" {
									return true, nil
								} else {
									return false, errors.New(fmt.Sprintf("HelmRelease not ready [%s]: %s", c.Reason, c.Message))
								}
							}
						}

						return false, errors.New("HelmRelease not ready")
					}).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(BeTrue(), fmt.Sprintf("%s HelmRelease should be ready", component))

					By(fmt.Sprintf("Verifying %s is deployed", component))
					Eventually(wait.IsAppDeployed(context.Background(),
						state.GetFramework().MC(),
						appName,
						testNamespace)).
						WithTimeout(appReadyTimeout).
						WithPolling(appReadyInterval).
						Should(BeTrue(), fmt.Sprintf("%s should be deployed", component))
				}
			})
		}).
		AfterSuite(func() {
			// Cleanup namespace if needed
			if state.GetFramework() != nil && state.GetFramework().MC() != nil {
				ctx := context.Background()
				ns := &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: testNamespace,
					},
				}
				err := state.GetFramework().MC().Delete(ctx, ns)
				if err != nil && !k8serrors.IsNotFound(err) {
					fmt.Printf("Warning: Failed to delete test namespace: %v\n", err)
				}
			}
		}).
		Run(t, "Basic Test")
}

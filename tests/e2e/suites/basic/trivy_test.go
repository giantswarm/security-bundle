package basic

import (
	"context"
	"fmt"
	"strings"
	"time"

	clusterclient "github.com/giantswarm/clustertest/v5/pkg/client"
	"github.com/giantswarm/clustertest/v5/pkg/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// Scanning goes through a scan job plus the trivy server, so this is
	// deliberately the most generous timeout in the suite.
	vulnerabilityReportTimeout = 15 * time.Minute
	metricsTimeout             = 5 * time.Minute

	starboardExporterInstance = "starboard-exporter"
	starboardMetricsPort      = 8080

	// severityCountMetric is emitted per image and severity from the report
	// summary, and carries image_repository in its label set.
	severityCountMetric = "starboard_exporter_vulnerabilityreport_image_vulnerability_severity_count"
)

// registerTrivyScenarios checks that trivy-operator scans a workload deployed
// into an arbitrary namespace, and that starboard-exporter turns the resulting
// report into metrics. trivy-operator runs with an empty targetNamespaces, so
// it watches all namespaces and needs no extra configuration here.
func registerTrivyScenarios() {
	It("should create a VulnerabilityReport for the deployed workload image", func() {
		ctx := context.Background()
		wc := wcClient()

		By(fmt.Sprintf("Waiting for a VulnerabilityReport covering %s", scanTargetRepository))
		Eventually(func() error {
			_, err := findVulnerabilityReport(ctx, wc)
			return err
		}).WithTimeout(vulnerabilityReportTimeout).WithPolling(pollingInterval).Should(Succeed())
	})

	It("should expose starboard-exporter metrics for that VulnerabilityReport", func() {
		ctx := context.Background()
		wc := wcClient()

		By("Locating the starboard-exporter metrics Service")
		services := &corev1.ServiceList{}
		Expect(wc.List(ctx, services, client.MatchingLabels{instanceLabel: starboardExporterInstance})).To(Succeed())
		Expect(services.Items).NotTo(BeEmpty(),
			"expected a Service labelled %s=%s", instanceLabel, starboardExporterInstance)
		exporter := services.Items[0]

		By("Finding a pod to scrape from")
		var podName string
		Eventually(func() error {
			var err error
			podName, err = runningPod(ctx, wc, compliantWorkload)
			return err
		}).WithTimeout(workloadReadyTimeout).WithPolling(pollingInterval).Should(Succeed())

		metricsURL := fmt.Sprintf("http://%s.%s.svc:%d/metrics", exporter.Name, exporter.Namespace, starboardMetricsPort)

		// The exporter image is distroless, so the scrape runs from the
		// workload pod instead. Its NetworkPolicy allows ingress on the metrics
		// port from anywhere in the cluster.
		By(fmt.Sprintf("Scraping %s from pod %s", metricsURL, podName))
		Eventually(func() error {
			stdout, stderr, err := wc.ExecInPod(ctx, podName, scenarioNamespace, workloadContainer,
				[]string{"wget", "-q", "-O-", metricsURL})
			if err != nil {
				return fmt.Errorf("scraping %s failed: %w (stderr: %q)", metricsURL, err, stderr)
			}

			if sample, found := findMetricSample(stdout, severityCountMetric, scanTargetRepository); found {
				logger.Log("Found metric sample: %s", sample)
				return nil
			}
			return fmt.Errorf("no %s sample with image_repository=%q in the exporter output",
				severityCountMetric, scanTargetRepository)
		}).WithTimeout(metricsTimeout).WithPolling(pollingInterval).Should(Succeed())
	})
}

// findVulnerabilityReport looks for a report in the scenario namespace whose
// scanned artifact is the workload image.
func findVulnerabilityReport(ctx context.Context, c *clusterclient.Client) (*unstructured.Unstructured, error) {
	reports := &unstructured.UnstructuredList{}
	reports.SetGroupVersionKind(vulnerabilityReportListGVK)

	if err := c.List(ctx, reports, client.InNamespace(scenarioNamespace)); err != nil {
		return nil, err
	}

	for i := range reports.Items {
		report := reports.Items[i]
		repository, found, err := unstructured.NestedString(report.Object, "report", "artifact", "repository")
		if err != nil || !found {
			continue
		}
		if repository == scanTargetRepository {
			logger.Log("Found VulnerabilityReport %s/%s for %s", report.GetNamespace(), report.GetName(), repository)
			return &report, nil
		}
	}

	return nil, fmt.Errorf("no VulnerabilityReport for repository %q in namespace %s",
		scanTargetRepository, scenarioNamespace)
}

// findMetricSample returns the first exposition line for the given metric whose
// image_repository label matches, ignoring HELP and TYPE comments.
func findMetricSample(exposition, metric, repository string) (string, bool) {
	labelMatch := fmt.Sprintf("image_repository=%q", repository)

	for _, line := range strings.Split(exposition, "\n") {
		if strings.HasPrefix(line, "#") || !strings.HasPrefix(line, metric) {
			continue
		}
		if strings.Contains(line, labelMatch) {
			return strings.TrimSpace(line), true
		}
	}
	return "", false
}

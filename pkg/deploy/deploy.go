package deploy

import (
	"github.com/justmike1/deployer/pkg/config"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

func HelmChart(clusterName, chartURI, namespace, valuesFile string, repoURL ...string) {
	chartName := extractChartName(chartURI)
	if chartName == "" {
		log.Fatalf("Invalid chart URI: %s", chartURI)
	}

	kubeconfigContent, err := config.GetKubeconfigContent(clusterName)
	if err != nil {
		log.Fatalf("Failed to get kubeconfig for cluster %s: %v", clusterName, err)
	}
	kubeconfigPath, err := config.CreateTempKubeconfigFile(kubeconfigContent)
	if err != nil {
		log.Fatalf("Failed to create temporary kubeconfig file: %v", err)
	}
	defer os.Remove(kubeconfigPath)

	args := []string{
		"upgrade", "--install", chartName, chartURI,
		"-n", namespace, "--create-namespace",
		"--kubeconfig", kubeconfigPath,
	}
	if len(repoURL) > 0 {
		args = append(args, "--repo", repoURL[0])
	}
	if valuesFile != "" {
		args = append(args, "-f", valuesFile)
	}

	cmd := exec.Command("helm", args...)
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()

	log.Printf("Running Helm command: helm %v\n", args)
	if err := cmd.Run(); err != nil {
		log.Fatalf("Helm install failed: %v", err)
	}
}

func extractChartName(chartURI string) string {
	noScheme := strings.TrimPrefix(chartURI, "oci://")
	return path.Base(noScheme)
}

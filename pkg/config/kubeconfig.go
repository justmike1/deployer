package config

import (
	"fmt"
	"os"
	"os/exec"
)

func GetKubeconfigContent(clusterName string) (string, error) {
	cmd := exec.Command("k3d", "kubeconfig", "get", clusterName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get kubeconfig for cluster %s: %w", clusterName, err)
	}

	return string(output), nil
}

func CreateTempKubeconfigFile(kubeconfigContent string) (string, error) {
	tmpFile, err := os.CreateTemp("", "kubeconfig-*.yaml")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary kubeconfig file: %w", err)
	}

	if _, err := tmpFile.Write([]byte(kubeconfigContent)); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to write kubeconfig content to temporary file: %w", err)
	}

	tmpFile.Close()
	return tmpFile.Name(), nil
}

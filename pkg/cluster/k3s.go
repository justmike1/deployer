package cluster

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func K3dCluster(clusterName string) {
	if isK3dInstalled() {
		log.Println("k3d is already installed.")
	} else {
		log.Println("k3d is not installed. Starting installation...")
		installK3d()
	}

	if isK3dClusterRunning(clusterName) {
		log.Printf("%v cluster is already running.\n", clusterName)
	} else {
		log.Printf("%v cluster is not running. Starting a new cluster...\n", clusterName)
		createK3dCluster(clusterName)
	}
}

func isK3dInstalled() bool {
	_, err := exec.LookPath("k3d")
	return err == nil
}

func installK3d() {
	cmd := exec.Command("sh", "-c", "curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to install k3d: %v", err)
	}

	log.Println("k3d has been successfully installed.")
}

func isK3dClusterRunning(clusterName string) bool {
	cmd := exec.Command("k3d", "cluster", "list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to check k3d clusters: %v", err)
		return false
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] == clusterName {
			return true
		}
	}
	return false
}

func createK3dCluster(clusterName string) {
	cmd := exec.Command("k3d", "cluster", "create", clusterName,
		"--port", "80:80@loadbalancer",
		"--port", "443:443@loadbalancer",
		"--port", "30080:80@loadbalancer", // NodePort for HTTP
		"--port", "30443:443@loadbalancer", // NodePort for HTTPS
		"--agents", "1",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to create k3d cluster: %v", err)
	}

	log.Printf("%v cluster has been successfully created.", clusterName)

	outputPath := "/tmp/kubeconfig.yaml"
	exportCmd := exec.Command("k3d", "kubeconfig", "get", clusterName)
	kubeconfig, err := exportCmd.Output()
	if err != nil {
		log.Fatalf("Failed to export kubeconfig: %v", err)
	}
	if err := os.WriteFile(outputPath, kubeconfig, 0644); err != nil {
		log.Fatalf("Failed to write kubeconfig to %s: %v", outputPath, err)
	}
	log.Printf("Wrote kubeconfig to %s", outputPath)
}

func DestroyK3dCluster(clusterName string) {
	if !isK3dInstalled() {
		log.Println("k3d is not installed.")
		return
	}

	if !isK3dClusterRunning(clusterName) {
		log.Printf("Cluster %s is not running or doesn't exist.\n", clusterName)
		return
	}

	log.Printf("Deleting k3d cluster %s...\n", clusterName)
	cmd := exec.Command("k3d", "cluster", "delete", clusterName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to delete k3d cluster: %v", err)
	}

	log.Printf("Cluster %s has been successfully deleted.\n", clusterName)
}

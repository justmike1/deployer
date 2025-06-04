package main

import (
	"flag"
	"github.com/justmike1/deployer/pkg/cluster"
	"github.com/justmike1/deployer/pkg/config"
	"github.com/justmike1/deployer/pkg/deploy"
	"github.com/justmike1/deployer/pkg/setup"
	"github.com/justmike1/deployer/pkg/status"
	"log"
	"os"
)

const defaultClusterName = "local"
const defaultNamespace = "default"

func main() {
	if len(os.Args) < 2 {
		log.Println("No command provided. Use 'help' command to see available commands.")
		os.Exit(1)
	}

	cmdOption := config.ParseCommandName(os.Args[1])
	// Parse flags from remaining args (skip the program and command name)
	config.ParseFlags(os.Args[2:])
	// Cluster name is optional argument (3rd position), fallback to default
	clusterName := defaultClusterName
	if len(flag.Args()) > 0 {
		clusterName = flag.Args()[0]
	}

	if os.Geteuid() != 0 {
		log.Println("This CLI tool operations requires root privileges (e.g., run with sudo)")
		return
	}

	switch cmdOption {
	case config.SETUP:
		setupPrerequisites()
	case config.CLUSTER:
		deployK3s(clusterName)
	case config.DEPLOY:
		deployK8sManifests(clusterName)
	case config.STATUS:
		showStatus(clusterName)
	case config.DESTROY:
		destroyK3s(clusterName)
	default:
		log.Println("Unknown command. Use 'help' command to see available commands.")
		os.Exit(1)
	}
}

func setupPrerequisites() {
	log.Println("Setting up system tools...")
	setup.InstallTools()
}

func deployK3s(clusterName string) {
	log.Printf("Deploying a local K3s cluster [%s]...\n", clusterName)
	cluster.K3dCluster(clusterName)
}

func deployK8sManifests(clusterName string) {
	if config.HelmChart != "" {
		log.Println("Deploying Helm chart...")
		deploy.HelmChart(clusterName, config.HelmChart, config.Namespace, config.ValuesFile, config.RepoURL)
	} else {
		log.Panicf("Helm chart URI is required for deployment. Use --helm flag to specify the chart.")
	}
}

func showStatus(clusterName string) {
	log.Printf("Checking the application status for cluster [%s]...\n", clusterName)
	status.LogPodStatuses(clusterName, defaultNamespace)
}

func destroyK3s(clusterName string) {
	log.Printf("Destroying local K3s cluster [%s]...\n", clusterName)
	cluster.DestroyK3dCluster(clusterName)
}

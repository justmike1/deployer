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
		log.Println("No command provided.\n")
		help()
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

	switch cmdOption {
	case config.HELP:
		help()
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
		log.Println("Unknown command.\n")
		help()
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

func help() {
	log.Println("Usage: deployer <command> [flags] [cluster-name]")
	log.Println("")
	log.Println("Commands:")
	log.Println("  help       Show this help message")
	log.Println("  setup      Install prerequisite tools (k3d, helm, etc.)")
	log.Println("  cluster    Create a local K3s cluster")
	log.Println("  deploy     Deploy Helm chart to the cluster")
	log.Println("  status     Show status of application pods")
	log.Println("  destroy    Tear down the local K3s cluster")
	log.Println("")
	log.Println("Flags:")
	log.Println("  --helm     Helm chart URI (OCI or repo)")
	log.Println("  -n         Kubernetes namespace (default: default)")
	log.Println("  -f         Helm values file")
	log.Println("  --repo     Helm repository URL")
	log.Println("")
	log.Println("Example:")
	log.Println("  deployer deploy --helm oci://myrepo/mychart -n dev -f values.yaml my-cluster")
}

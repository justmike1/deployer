package config

import "flag"

// CmdOptions represents the available command options as an enum
type CmdOptions int

const (
	CLUSTER CmdOptions = iota
	SETUP
	DEPLOY
	STATUS
	DESTROY
)

func (c CmdOptions) String() string {
	switch c {
	case SETUP:
		return "setup"
	case CLUSTER:
		return "cluster"
	case DEPLOY:
		return "deploy"
	case STATUS:
		return "status"
	case DESTROY:
		return "destroy"
	default:
		return "unknown"
	}
}

var (
	HelmChart  string
	Namespace  string
	ValuesFile string
	RepoURL    string
)

// ParseCommandName returns the enum command option from CLI arg
func ParseCommandName(option string) CmdOptions {
	switch option {
	case "setup":
		return SETUP
	case "cluster":
		return CLUSTER
	case "deploy":
		return DEPLOY
	case "status":
		return STATUS
	case "destroy":
		return DESTROY
	default:
		return -1
	}
}

// ParseFlags parses flags like --helm, -n, -f after command name
func ParseFlags(args []string) {
	flag.StringVar(&HelmChart, "helm", "", "Helm chart URI (OCI or repo)")
	flag.StringVar(&Namespace, "n", "default", "Kubernetes namespace")
	flag.StringVar(&ValuesFile, "f", "", "Helm values file")
	flag.StringVar(&RepoURL, "repo", "", "Helm repository URL")

	// only parse args after the command
	_ = flag.CommandLine.Parse(args)
}

package setup

import (
	"log"
	"os/exec"
	"runtime"
	"strings"
)

func InstallTools() {
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		log.Fatal("This setup script currently supports only Linux amd64.")
	}

	installDocker()
	installHelm()
	installK9s()
	installKubectl()
}

func installDocker() {
	if isInstalled("docker", "--version") {
		log.Println("Docker is already installed. Skipping.")
		return
	}

	distro := detectDistro()
	log.Printf("Detected Linux distribution: %s", distro)

	switch distro {
	case "amzn": // Amazon Linux
		log.Println("Installing Docker on Amazon Linux...")
		cmd := exec.Command("sh", "-c", `
if grep -q '2023' /etc/os-release; then
  yum install -y docker
else
  amazon-linux-extras install -y docker
fi
groupadd docker || true
usermod -aG docker $(whoami)
service docker start
systemctl enable docker || true
`)
		runOrExit(cmd, "Docker")

	case "ubuntu", "debian":
		log.Println("Installing Docker on Ubuntu/Debian...")
		cmd := exec.Command("sh", "-c", `
apt-get update
apt-get install -y docker.io
groupadd docker || true
usermod -aG docker $(whoami)
systemctl enable docker
systemctl start docker
`)
		runOrExit(cmd, "Docker")

	default:
		log.Fatalf("Unsupported distro: %s. Please install Docker manually.", distro)
	}

	log.Println("Docker installed successfully.")
}

func installHelm() {
	if isInstalled("helm", "version") {
		log.Println("Helm is already installed. Skipping.")
		return
	}

	log.Println("Installing Helm...")

	cmd := exec.Command("sh", "-c", `
curl -fsSL https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
`)
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to install Helm: %v", err)
	}

	log.Println("Helm installed successfully.")
}

func installK9s() {
	if isInstalled("k9s", "version") {
		log.Println("K9s is already installed. Skipping.")
		return
	}

	log.Println("Installing K9s...")

	cmd := exec.Command("sh", "-c", "curl -sS https://webinstall.dev/k9s | bash")
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to install K9s: %v", err)
	}

	log.Println("K9s installed successfully.")
}

func installKubectl() {
	if isInstalled("kubectl", "version --client=true") {
		log.Println("kubectl is already installed. Skipping.")
		return
	}

	log.Println("Installing kubectl...")

	// Download kubectl using correct substitution with -L for redirection
	cmd := exec.Command("sh", "-c", `
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x kubectl
mv kubectl /usr/local/bin/kubectl
`)
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to install kubectl: %v", err)
	}

	log.Println("kubectl installed successfully.")
}

func isInstalled(command string, versionArg string) bool {
	cmd := exec.Command("sh", "-c", command+" "+versionArg)
	err := cmd.Run()
	return err == nil
}

func detectDistro() string {
	out, err := exec.Command("sh", "-c", "source /etc/os-release && echo $ID").Output()
	if err != nil {
		log.Fatalf("Failed to detect Linux distribution: %v", err)
	}
	return strings.TrimSpace(string(out))
}

func runOrExit(cmd *exec.Cmd, toolName string) {
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to install %s: %v", toolName, err)
	}
}

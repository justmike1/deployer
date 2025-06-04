package status

import (
	"context"
	"github.com/justmike1/deployer/pkg/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

func LogPodStatuses(clusterName string, namespace string) {
	kubeconfigContent, err := config.GetKubeconfigContent(clusterName)
	if err != nil {
		log.Fatalf("Failed to get kubeconfig for cluster %s: %v", clusterName, err)
	}
	kubeconfigPath, err := config.CreateTempKubeconfigFile(kubeconfigContent)
	if err != nil {
		log.Fatalf("Failed to create temporary kubeconfig file: %v", err)
	}
	defer os.Remove(kubeconfigPath)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		log.Fatalf("Failed to load kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes clientset: %v", err)
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to list pods in %s namespace: %v", namespace, err)
	}

	for _, pod := range pods.Items {
		log.Printf("Pod: %s, Status: %s\n", pod.Name, pod.Status.Phase)
		for _, containerStatus := range pod.Status.ContainerStatuses {
			log.Printf("  Container: %s, Ready: %v, RestartCount: %d, State: %v\n",
				containerStatus.Name,
				containerStatus.Ready,
				containerStatus.RestartCount,
				containerStatus.State)
		}
	}
}

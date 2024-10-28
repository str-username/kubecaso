package kubecaso

import (
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client

// KubernetesClient structure
type KubernetesClient struct {
	This *kubernetes.Clientset
}

// NewKubernetesClient create kubernetes client
func NewKubernetesClient() (*KubernetesClient, error) {
	var config *rest.Config
	var err error

	// Check where app running
	if _, inCluster := os.LookupEnv("KUBERNETES_SERVICE_HOST"); inCluster {
		// If in-cluster, use in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		log.Info().Str("kubecaso", "NewKubernetesClient").Msg("in-cluster client")
	} else {
		// If local, use kubeconfig
		kubeConfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			return nil, err
		}
		log.Info().Str("kubecaso", "NewKubernetesClient").Msg("local client")
	}

	// Create kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	log.Info().Str("kubecaso", "NewKubernetesClient").Msg("create client")
	return &KubernetesClient{This: clientset}, nil
}

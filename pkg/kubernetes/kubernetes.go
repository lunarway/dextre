package kubernetes

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	clientset  *kubernetes.Clientset
	kubeConfig string
}

// NewKubernets creates a new container with kubeclient, and the generic flags
func NewClient(kubeConfig string) (*Client, error) {
	// Configure kubeconfig if not set.
	if kubeConfig == "" {
		kubeConfig = initKubeConfig()
	}

	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfig},
		&clientcmd.ConfigOverrides{})

	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("kubeconfig: %s", kubeConfig))
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("kubeconfig: %s", kubeConfig))
	}

	client := Client{
		kubeConfig: kubeConfig,
		clientset:  clientset,
	}

	return &client, nil
}

func initKubeConfig() string {
	if os.Getenv("KUBECONFIG") != "" {
		return os.Getenv("KUBECONFIG")
	}
	return clientcmd.RecommendedHomeFile
}

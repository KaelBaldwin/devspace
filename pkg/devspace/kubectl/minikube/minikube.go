package minikube

import (
	"github.com/devspace-cloud/devspace/pkg/devspace/config/versions/latest"
	"k8s.io/client-go/tools/clientcmd"
)

var isMinikubeVar *bool

// IsMinikube returns true if the Kubernetes cluster is a minikube
func IsMinikube(config *latest.Config) bool {
	if isMinikubeVar == nil {
		isMinikube := false
		if config == nil || config.Cluster == nil || config.Cluster.APIServer == nil {
			if config == nil || config.Cluster == nil || config.Cluster.KubeContext == nil {
				loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
				kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
				cfg, err := kubeConfig.RawConfig()
				if err != nil {
					return false
				}

				isMinikube = cfg.CurrentContext == "minikube"
			} else {
				isMinikube = *config.Cluster.KubeContext == "minikube"
			}
		}

		isMinikubeVar = &isMinikube
	}

	return *isMinikubeVar
}

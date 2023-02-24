package internal

import (
	"flag"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func LoadKubeConfig(kubeConfigPath string) (*kubernetes.Clientset, *dynamic.DynamicClient, error) {
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)

	if err != nil {
		return nil, nil, err
	}

	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	dd, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	return c, dd, nil
}

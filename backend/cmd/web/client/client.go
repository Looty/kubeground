package client

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type Client struct {
	ClientSet *kubernetes.Clientset
	Dynamic   *dynamic.DynamicClient
}

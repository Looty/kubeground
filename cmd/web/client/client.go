package client

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/Looty/kubeground/internal"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/homedir"
)

type Client struct {
	ClientSet *kubernetes.Clientset
	Dynamic   *dynamic.DynamicClient
}

func Init(c *Client) {
	inCluster := flag.Lookup("incluster").Value.(flag.Getter).Get().(bool)

	if inCluster {
		config, err := rest.InClusterConfig()
		if err != nil {
			log.Println("Failed to get InClusterConfig")
			panic(err)
		}

		clientSet, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Println("Failed to get InCluster ClientSet")
			panic(err)
		}

		dynamicClient, err := dynamic.NewForConfig(config)
		if err != nil {
			log.Println("Failed to get InCluster DynamicClient")
			panic(err)
		}

		*c = Client{
			ClientSet: clientSet,
			Dynamic:   dynamicClient,
		}
	} else {
		cc, dd, err := internal.LoadKubeConfig(filepath.Join(homedir.HomeDir(), ".kube", "config"))

		if err != nil {
			log.Printf("Please check if you have a local cluster running: %s", err)
			panic(err)
		}

		*c = Client{
			ClientSet: cc,
			Dynamic:   dd,
		}
	}
}

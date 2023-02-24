package client

import (
	"fmt"
	"path/filepath"

	"github.com/Looty/kubeground/internal"
	"k8s.io/client-go/util/homedir"
)

func Init(c *Client) {
	cc, dd, err := internal.LoadKubeConfig(filepath.Join(homedir.HomeDir(), ".kube", "config"))

	if err != nil {
		fmt.Println("Please check if you have a local cluster running")
		panic(err)
	}

	*c = Client{
		ClientSet: cc,
		Dynamic:   dd,
	}
}

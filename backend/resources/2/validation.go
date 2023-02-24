package levels

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// TODO: continue
func a() {

	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	pod, err := clientset.CoreV1().Pods("namespace").Get(
		context.Background(),
		"podName",
		v1.GetOptions{},
	)

	// if err != nil {
	// 	return "", err
	// }

	labelValue, ok := pod.ObjectMeta.Labels["labelKey"]
	if !ok {
		fmt.Printf("no label with key %s for pod %s/%s", "labelKey", "namespace", "podName")
	}

	fmt.Println(strings.ToUpper(labelValue))

	// return strings.ToUpper(labelValue), nil
}

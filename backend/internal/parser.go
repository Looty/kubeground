package internal

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
)

func ParseResourceFolder(folder string, dd *dynamic.DynamicClient, c *kubernetes.Clientset, create bool) {
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if !info.IsDir() { // is a file
			if create {
				ParseResourceFile(path, dd, c, true)
			} else {
				ParseResourceFile(path, dd, c, false)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}

func ParseResourceFile(file string, dd *dynamic.DynamicClient, c *kubernetes.Clientset, create bool) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	if create {
		CreateResource(f, dd, c)
	} else {
		DeleteResource(f, dd, c)
	}
}

func ReadResource(yamlFile []byte, dd *dynamic.DynamicClient, c *kubernetes.Clientset, dr *dynamic.ResourceInterface, unr *unstructured.Unstructured) {
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(yamlFile), 100)

	for {
		var rawObj runtime.RawExtension
		if err := decoder.Decode(&rawObj); err != nil {
			break
		}

		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			log.Fatal(err)
		}

		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

		gr, err := restmapper.GetAPIGroupResources(c.Discovery())
		if err != nil {
			log.Fatal(err)
		}

		mapper := restmapper.NewDiscoveryRESTMapper(gr)
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			log.Fatal(err)
		}

		var dri dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			if unstructuredObj.GetNamespace() == "" {
				unstructuredObj.SetNamespace("default")
			}
			dri = dd.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
		} else {
			dri = dd.Resource(mapping.Resource)
		}

		*dr = dri
		*unr = *unstructuredObj
	}
}

func CreateResource(yamlFile []byte, dd *dynamic.DynamicClient, c *kubernetes.Clientset) {
	var dri dynamic.ResourceInterface
	var unr unstructured.Unstructured

	ReadResource(yamlFile, dd, c, &dri, &unr)

	if _, err := dri.Create(context.Background(), &unr, metav1.CreateOptions{}); err != nil {
		log.Println("Couldn't create resource")
	}
}

func DeleteResource(yamlFile []byte, dd *dynamic.DynamicClient, c *kubernetes.Clientset) {
	var dri dynamic.ResourceInterface
	var unr unstructured.Unstructured

	ReadResource(yamlFile, dd, c, &dri, &unr)

	if err := dri.Delete(context.Background(), unr.GetName(), metav1.DeleteOptions{}); err != nil {
		log.Println("Couldn't delete resource")
	}
}

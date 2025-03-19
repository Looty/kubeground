package quest

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/Looty/kubeground/cmd/web/client"
	"github.com/Looty/kubeground/internal"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
)

func ApplyManifests(id string) string {
	var cl client.Client
	client.Init(&cl)

	var errorMessage string

	data, err := cl.ClientSet.RESTClient().
		Get().
		AbsPath("/apis/quest.looty.com/v1/quests").
		DoRaw(context.TODO())

	if err != nil {
		log.Printf("Failed to fetch quest in API: %s", data)
		log.Print(err)
		return err.Error()
	}

	var manifests string

	// Unmarshal the JSON data into the Response struct
	var response internal.QuestList
	if err := json.Unmarshal(data, &response); err != nil {
		log.Println("Failed to Unmarshal YAML manifests")
		log.Print(err)
		return ""
	}

	// Print the resource name
	for _, item := range response.Quests {
		if strconv.Itoa(item.Spec.Level) == id {
			manifests = item.Spec.Manifests
		}
	}

	if manifests == "" {
		errorMessage = "No quest resources were defined, nothing to apply.."
		log.Println(errorMessage)
		return errorMessage
	}

	log.Println("Attempting to apply quest resources..")

	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(manifests)), 100)
	for {
		var rawObj runtime.RawExtension
		err := decoder.Decode(&rawObj)
		if err != nil {
			log.Println("Failed to convert apply YAML manifests")
			log.Print(err)
			return ""
		}

		log.Println("Debug2:", rawObj)

		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if err != nil {
			log.Println("Failed to convert apply YAML manifests 2")
			log.Print(err)
			return ""
		}
		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			log.Println("Failed to convert apply YAML manifests 3")
			log.Print(err)
			return ""
		}

		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}
		log.Println("Debug3:", unstructuredMap)

		gr, err := restmapper.GetAPIGroupResources(cl.ClientSet.Discovery())
		if err != nil {
			log.Println("Failed to convert apply YAML manifests 4")
			log.Print(err)
			return ""
		}

		log.Println("Debug4:", gr)

		mapper := restmapper.NewDiscoveryRESTMapper(gr)
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			log.Println("Failed to convert apply YAML manifests 5")
			log.Print(err)
			return ""
		}

		log.Println("Debug5:", mapper)
		log.Println("Debug5:", mapping)

		var dri dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			if unstructuredObj.GetNamespace() == "" {
				unstructuredObj.SetNamespace("default")
			}
			dri = cl.Dynamic.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
		} else {
			dri = cl.Dynamic.Resource(mapping.Resource)
		}

		if _, err := dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{}); err != nil {
			log.Println("Failed to apply YAML manifests")
			log.Print(err)
			return err.Error()
		}
	}
}

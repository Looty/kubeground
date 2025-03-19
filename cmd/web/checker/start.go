package checker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/Looty/kubeground/cmd/web/client"
	"github.com/Looty/kubeground/internal"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/util/wait"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"

	"github.com/Masterminds/sprig/v3"
	"github.com/labstack/echo/v4"
)

const jobTemplate = `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: {{ .Name }}
subjects:
  - kind: ServiceAccount
    name: {{ .Name }}
    namespace: {{ .Namespace }}
roleRef:
  kind: ClusterRole
  name: kubeground-checker-role
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: batch/v1
kind: Job
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  ttlSecondsAfterFinished: {{ .CheckerTTL }}
  backoffLimit: 0
  completions: 1
  template:
    metadata:
      name: {{ .Name }}
    spec:
      restartPolicy: Never
      serviceAccountName: {{ .Name }}
      containers:
        - name: kubectl
          image: bitnami/kubectl
          resources:
            requests:
              memory: "32Mi"
              cpu: "50m"
            limits:
              memory: "64Mi"
              cpu: "100m"
          command:
          - bash
          - -c
          - |
{{ .Validation | indent 12 }}
`

func ValidateQuest(c echo.Context) error {
	id := c.Param("id")
	var alert string
	var userMessage string

	var jobData = internal.JobData{}
	var cl client.Client
	client.Init(&cl)

	// Start the long-running process in a separate goroutine
	go func() {

		data, err := cl.ClientSet.RESTClient().
			Get().
			AbsPath("/apis/checker.looty.com/v1/checkers").
			DoRaw(context.TODO())

		if err != nil {
			log.Printf("Validating quest failed: %s", data)
			log.Print(err)
			// return err
		}

		// Unmarshal the JSON data into the Response struct
		var response internal.CheckerList
		if err := json.Unmarshal(data, &response); err != nil {
			log.Println("Unmarshaling checker failed")
			log.Print(err)
			// return err
		}

		// Print the resource name
		// TODO: you can filter the resoure directly above, without needing to find it
		for _, item := range response.Checkers {
			if item.Spec.QuestRef == id {

				jobData = internal.JobData{
					Name:        "checker-" + item.Metadata.Name,
					Namespace:   item.Metadata.Namespace,
					CheckerTTL:  "15",
					CheckerName: item.Metadata.Name,
					Validation:  item.Spec.Validation,
				}

				log.Printf("Caught job!\nName: %s\nNamespace: %s\nValidation: %s",
					"checker-"+item.Metadata.Name,
					item.Metadata.Namespace,
					item.Spec.Validation,
				)
			}
		}

		alert = fmt.Sprintf("Validating quest %s.. Please wait a moment and refresh the page shortly...", id)
		internal.AddMessageToQueue("warning", alert, "", true)

		log.Println("Attempting to apply new resources..")

		tmpl, err := template.New("checkerjob").Funcs(sprig.TxtFuncMap()).Parse(jobTemplate)
		if err != nil {
			log.Printf("checker template error: %v", err)
		}

		var outputBuffer bytes.Buffer
		tmpl.Execute(&outputBuffer, jobData)

		log.Printf("parsed job resources: " + outputBuffer.String())

		log.Printf("deleting job resources before applying..")
		backgroundDeletion := metav1.DeletePropagationBackground

		cl.ClientSet.BatchV1().Jobs(jobData.Namespace).Delete(context.TODO(), jobData.Name, metav1.DeleteOptions{
			PropagationPolicy: &backgroundDeletion,
		})
		cl.ClientSet.RbacV1().ClusterRoles().Delete(context.TODO(), jobData.Name, metav1.DeleteOptions{})
		cl.ClientSet.RbacV1().ClusterRoleBindings().Delete(context.TODO(), jobData.Name, metav1.DeleteOptions{})
		cl.ClientSet.CoreV1().ServiceAccounts(jobData.Namespace).Delete(context.TODO(), jobData.Name, metav1.DeleteOptions{})

		decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(outputBuffer.String())), 100)
		for {
			var rawObj runtime.RawExtension
			err := decoder.Decode(&rawObj)
			if err != nil {
				log.Println("Failed to convert YAML to kubernetes objects in checker")
				log.Print(err)
				break
			}

			obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
			unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
			if err != nil {
				log.Println("Failed to convert YAML to kubernetes objects in checker 2")
				log.Print(err)
				// return err
			}

			unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

			gr, err := restmapper.GetAPIGroupResources(cl.ClientSet.Discovery())
			if err != nil {
				log.Println("Failed to convert YAML to kubernetes objects in checker 3")
				log.Print(err)
				// return err
			}

			mapper := restmapper.NewDiscoveryRESTMapper(gr)
			mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
			if err != nil {
				log.Println("Failed to convert YAML to kubernetes objects in checker 4")
				log.Print(err)
				// return err
			}

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
				log.Println("Failed to create checker resources")
				log.Print(err)
				// return err
			}
		}

		log.Println("Checking job result..")

		var exitCode int32
		_ = wait.PollImmediate(3*time.Second, 2*time.Minute, func() (bool, error) {
			job, err := cl.ClientSet.BatchV1().Jobs(jobData.Namespace).Get(context.TODO(), jobData.Name, metav1.GetOptions{})
			if err != nil {
				return false, err
			}
			if job.Status.Active == 0 && job.Status.Succeeded == 0 && job.Status.Failed == 0 {
				log.Println("Job is still not ready..")
				return false, nil
			}

			if job.Status.Active > 0 {
				log.Println("Job is still not active..")
				return false, nil
			}

			if job.Status.Succeeded > 0 {
				log.Println("Job Succeeded..")

				// Job has succeeded, get the exit code from the pod
				podList, err := cl.ClientSet.CoreV1().Pods(jobData.Namespace).List(context.TODO(), metav1.ListOptions{})
				if err != nil {
					return false, err
				}
				if len(podList.Items) == 0 {
					// No pods found yet
					return false, nil
				}

				// Get the exit code from the first container in the first pod
				for _, pod := range podList.Items {
					for _, containerStatus := range pod.Status.ContainerStatuses {
						if containerStatus.State.Terminated != nil {
							exitCode = containerStatus.State.Terminated.ExitCode
							return true, nil
						}
					}
				}
				return true, nil
			}

			if job.Status.Failed == 1 {
				log.Println("Job failed..")
				exitCode = 1
				return true, nil
			}

			return false, nil
		})

		userMessage = "Sorry, something is still missing.."
		if exitCode == 0 {
			userMessage = "Congratulations on completing this quest!"

			var questDataFound = internal.Quest{}
			questsData, err := cl.ClientSet.RESTClient().
				Get().
				AbsPath("/apis/quest.looty.com/v1/quests").
				DoRaw(context.TODO())

			if err != nil {
				log.Println("Couldn't find quests!")
				log.Print(err)
				// return err
			}

			// Unmarshal the JSON data into the Response struct
			var response internal.QuestList
			if err := json.Unmarshal(questsData, &response); err != nil {
				log.Println("Unmarshaling quests failed")
				log.Print(err)
				// return err
			}

			// Print the resource name
			// TODO: you can filter the resoure directly above, without needing to find it
			id_to_string, err := strconv.Atoi(id)
			if err != nil {
				log.Println("Unmarshaling level to string failed")
				log.Print(err)
				// return err
			}

			for _, item := range response.Quests {
				if item.Spec.Level == id_to_string {

					questDataFound = internal.Quest{
						Metadata: internal.QuestMetadata{
							Name:      item.Metadata.Name,
							Namespace: item.Metadata.Namespace,
						},
						Spec: internal.QuestSpec{
							Completed:    item.Spec.Completed,
							Hints:        item.Spec.Hints,
							Instructions: item.Spec.Instructions,
							Level:        item.Spec.Level,
							Manifests:    item.Spec.Manifests,
						},
					}
				}
			}

			//TODO: update quest.Completed here instead
			questData, err := cl.ClientSet.RESTClient().
				Get().
				AbsPath("/apis/quest.looty.com/v1").
				Namespace(questDataFound.Metadata.Namespace).
				Resource("quests").
				Name(questDataFound.Metadata.Name).
				DoRaw(context.TODO())

			if err != nil {
				log.Printf("Validating quest failed: %s", questData)
				log.Println(err)
				// return err
			}

			// Unmarshal the JSON data into a map
			var quest map[string]interface{}
			if err := json.Unmarshal(questData, &quest); err != nil {
				log.Printf("Error unmarshaling Quest: %v", err)
				// return err
			}

			// Modify the replicas field
			spec, ok := quest["spec"].(map[string]interface{})
			if !ok {
				log.Print("Error: unable to find spec in Quest object")
			}
			spec["completed"] = true

			// Marshal the modified Quest object back to JSON
			updatedQuestData, err := json.Marshal(quest)
			if err != nil {
				log.Printf("Error marshaling updated Quest: %v", err)
				// return err
			}

			// Update the Quest object with the new configuration
			_, err = cl.ClientSet.RESTClient().
				Put().
				AbsPath("/apis/quest.looty.com/v1").
				Namespace(questDataFound.Metadata.Namespace).
				Resource("quests").
				Name(questDataFound.Metadata.Name).
				Body(updatedQuestData).
				DoRaw(context.TODO())

			if err != nil {
				log.Printf("Error updating Quest: %s", data)
				log.Println(err)
				// return err
			}

			log.Println("Quest completed successfully.")
			alert = fmt.Sprintf("Your answer for quest %s is correct!", id)
			internal.AddMessageToQueue("success", alert, "", false)
		} else {
			alert = fmt.Sprintf("Your answer for quest %s is incorrect, please try again", id)
			internal.AddMessageToQueue("danger", alert, "", false)
		}
	}()

	log.Println(userMessage)
	c.Redirect(http.StatusFound, "/")

	return c.String(http.StatusOK, userMessage)
}

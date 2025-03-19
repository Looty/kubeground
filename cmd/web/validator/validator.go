package validator

type Validator struct {
	PodName   string `yaml:"podName"`
	Namespace string `yaml:"namespace"`
}

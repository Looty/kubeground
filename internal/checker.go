package internal

type CheckerMetadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type CheckerSpec struct {
	QuestRef   string `json:"questref,omitempty"`
	Validation string `json:"validation,omitempty"`
}

type Checker struct {
	Metadata CheckerMetadata `json:"metadata"`
	Spec     CheckerSpec     `json:"spec"`
}

type CheckerList struct {
	Checkers []Checker `json:"items"`
}

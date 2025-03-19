package internal

type QuestMetadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type QuestSpec struct {
	Completed    bool   `json:"completed,omitempty"`
	Hints        string `json:"hints,omitempty"`
	Instructions string `json:"instructions,omitempty"`
	Level        int    `json:"level,omitempty"`
	Manifests    string `json:"manifests,omitempty"`
}

type Quest struct {
	Metadata QuestMetadata `json:"metadata"`
	Spec     QuestSpec     `json:"spec"`
}

type QuestList struct {
	Quests []Quest `json:"items"`
}

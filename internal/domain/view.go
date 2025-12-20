package domain

type PhaseView struct {
	Phase Phase `json:"phase"`
	Description string `json:"description"`
	Hints []string `json:"hints,omitempty"`
	Actions []Action `json:"actions,omitempty"`
}

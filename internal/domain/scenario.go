package domain

type Scenario struct {
	Meta ScenarioMeta `json:"meta"`
	Roles []Role `json:"roles"`
	Phases []PhaseContent `json:"phases"`
}

type ScenarioMeta struct {
	Title string `json:"title"`
	DurationMinutes int `json:"durationMinutes"`
	PlayerCount int `json:"playerCount"`
}

type Role struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Secret []Secret `json:"secret,omitempty"`
}

type Secret struct {
	TargetRoleID string `json:"targetRoleID,omitempty"`
	Content string `json:"content"`
}

type PhaseContent struct {
	Phase Phase `json:"phase"`
	Public PhasePublicInfo `json:"public"`
	Private []PhasePrivateInfo `json:"private,omitempty"`
}

type PhasePublicInfo struct {
	Description string `json:"description"`
	Actions []Action `json:"actions,omitempty"`
}

type PhasePrivateInfo struct {
	RoleID string `json:"roleId"`
	Hints []string `json:"hints"`
}

type Action struct {
	Type string `json:"type"`
	Description string `json:"description"`
}

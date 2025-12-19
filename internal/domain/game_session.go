package domain

type Session struct {
	ID       string
	Phase    Phase
	Scenario []byte // JSON Schema validated scenario
	Players  map[string]*Player
}



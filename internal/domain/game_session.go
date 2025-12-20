package domain

type Session struct {
	ID       string
	Phase    Phase
	Scenario *Scenario
	Players  map[string]*Player
}



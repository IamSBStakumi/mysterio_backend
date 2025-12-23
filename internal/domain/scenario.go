package domain

// Scenario はマーダーミステリーのシナリオを表す
type Scenario struct {
	Background string
	Setting    string
	Characters []Character
	Truth      string
	Phases     []Phase
	Difficulty Difficulty
}

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

// GetCharacterByPlayerID は指定されたプレイヤーIDのキャラクターを取得する
func (s *Scenario) GetCharacterByPlayerID(playerID string) *Character {
	for i := range s.Characters {
		if s.Characters[i].PlayerID == playerID {
			return &s.Characters[i]
		}
	}
	return nil
}

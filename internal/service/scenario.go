package service

import "context"

type ScenarioService struct {}

func NewScenarioService() *ScenarioService {
	return &ScenarioService{}
}

// 今はダミー。後でAI生成＋Schema validationに差し替える
func (s *ScenarioService) Generate(
	ctx context.Context,
	playerCount int,
	difficulty string,
) ([]byte, error) {

	dummy := []byte(`{
		"dummy": true
	}`)

	return dummy, nil
}

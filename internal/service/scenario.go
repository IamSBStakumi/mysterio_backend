package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

type ScenarioService struct {
	Schema *jsonschema.Schema
}

const (
	SCHEMA_PATH = "internal/schema/scenario.schema.json"
)

func NewScenarioService() (*ScenarioService, error) {
	compiler := jsonschema.NewCompiler()

	schemaBytes, err := os.ReadFile(SCHEMA_PATH)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	var schemaJSON any
	if err := json.Unmarshal(schemaBytes, &schemaJSON); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema: %w", err)
	}

	if err := compiler.AddResource(SCHEMA_PATH, schemaJSON); err != nil {
		return nil, fmt.Errorf("failed to add schema resource: %w", err)
	}

	schema, err := compiler.Compile(SCHEMA_PATH)
	if err != nil {
		return nil, fmt.Errorf("failed to compile schema: %w", err)
	}

	return &ScenarioService{Schema: schema}, nil
}

func (s *ScenarioService) Generate(
	ctx context.Context,
	playerCount int,
	difficulty string,
) ([]byte, error) {

	// 仮の生成結果。後でAIに差し替える
	scenarioJSON := []byte(`{
		"meta": {
			"title": "Dummy Mystery",
			"durationMinutes": 90,
			"playerCount": 5
		},
		"roles": [],
		"phases": []
	}`)

	// JSON Schema Validation
	if err := s.Schema.Validate(bytes.NewReader(scenarioJSON)); err != nil {
		return nil, fmt.Errorf("failed to validate scenario: %w", err)
	}

	return scenarioJSON, nil
}

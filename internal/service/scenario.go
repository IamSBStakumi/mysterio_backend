package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/IamSBStakumi/mysterio_backend/internal/domain"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

type ScenarioService struct {
	Schema *jsonschema.Schema
}

const (
	// SCHEMA_PATH = "internal/schema/scenario.schema.json"
	SCHEMA_PATH = "internal/schema/scenario.mvp.json"
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
) (*domain.Scenario, error) {

	// 仮の生成結果。後でAIに差し替える
	scenarioJSON := []byte(`{
		"meta": {
			"title": "Dummy Mystery",
			"durationMinutes": 90,
			"playerCount": 5
		},
		"roles": [
			{
				"id": "p1",
				"name": "Detective",
				"description": "You are a detective."
			},
			{
				"id": "p2",
				"name": "Witness",
				"description": "You saw something important."
			},
			{
				"id": "p3",
				"name": "Suspect",
				"description": "You are hiding something."
			}
		],
		"phases": [
			{
				"phase": "intro",
				"public": {
					"description": "The Story begins."
				}
			}
		]
	}`)

	// 1. JSON Unmarshal
	var raw any
	if err := json.Unmarshal(scenarioJSON, &raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal scenario: %w", err)
	}

	// 2. Schema Validation
	if err := s.Schema.Validate(raw); err != nil {
		return nil, fmt.Errorf("failed to validate scenario: %w", err)
	}

	// 3. struct にパース
	var scenario domain.Scenario
	if err := json.Unmarshal(scenarioJSON, &scenario); err != nil {
		return nil, fmt.Errorf("failed to parse scenario: %w", err)
	}

	return &scenario, nil
}

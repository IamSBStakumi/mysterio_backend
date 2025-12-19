package domain

import (
	"encoding/json"
	"time"
)

type GameSession struct {
	ID string
	Phase string
	Scenario json.RawMessage
	CreatedAt time.Time
}

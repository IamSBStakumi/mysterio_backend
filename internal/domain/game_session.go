package domain

import "time"

// SessionStatus はセッションの状態を表す
type SessionStatus string

const (
	StatusActive   SessionStatus = "active"
	StatusFinished SessionStatus = "finished"
)

// Difficulty はゲームの難易度を表す
type Difficulty string

const (
	DifficultyEasy Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard Difficulty = "hard"
)

// Session はゲームセッションを表す
type Session struct {
	ID            string
	OwnerPlayerID string
	PlayerIDs     []string
	Scenario      *Scenario
	CurrentPhase  int
	Votes         map[string]string // playerID -> targetPlayerID
	Status        SessionStatus
	CreatedAt     time.Time
}



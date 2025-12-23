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

// IsOwner は指定されたプレイヤーIDがセッションオーナーかどうかを判定する
func (s *Session) IsOwner(playerID string) bool {
	return s.OwnerPlayerID == playerID
}

// HasVoted は指定されたプレイヤーが既に投票済みかどうかを判定する
func (s *Session) HasVoted(playerID string) bool {
	_, exists := s.Votes[playerID]
	return exists
}

// IsVotingPhase は現在のフェーズが投票フェーズかどうかを判定する
func (s *Session) IsVotingPhase() bool {
	if s.CurrentPhase >= len(s.Scenario.Phases) {
		return false
	}
	return s.Scenario.Phases[s.CurrentPhase].Type == PhaseTypeVoting
}

// IsLastPhase は現在のフェーズが最終フェーズかどうかを判定する
func (s *Session) IsLastPhase() bool {
	return s.CurrentPhase >= len(s.Scenario.Phases)-1
}

// GetWinner は投票結果から犯人として最も票を集めたプレイヤーIDを返す
func (s *Session) GetWinner() string {
	voteCount := make(map[string]int)

	for _, targetID := range s.Votes {
		voteCount[targetID]++
	}

	maxVotes := 0
	winner := ""

	for playerID, count := range voteCount {
		if count > maxVotes {
			maxVotes = count
			winner = playerID
		}
	}

	return winner
}

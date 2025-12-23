package domain

// PhaseType はフェーズの種類を表す
type PhaseType string

const (
	PhaseTypeIntroduction  PhaseType = "introduction"
	PhaseInvestigation1 	PhaseType = "investigation1"
	PhaseInvestigation2 	PhaseType = "investigation2"
	PhaseInvestigation3 	PhaseType = "investigation3"
	PhaseInvestigation4 	PhaseType = "investigation4"
	PhaseTypeDiscussion    PhaseType = "discussion"
	PhaseTypeVoting        PhaseType = "voting"
	PhaseTypeReveal        PhaseType = "reveal"
)

// Phase はゲームの進行フェーズを表す
type Phase struct {
	Number      int
	Type        PhaseType
	Description string
	PublicText  string
	Duration    int // 分単位
}

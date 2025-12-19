package domain

type Phase string

const (
	PhaseIntro Phase = "intro"
	PhaseInvestigation1 Phase = "investigation1"
	PhaseInvestigation2 Phase = "investigation2"
	PhaseDiscussion Phase = "discussion"
	PhaseVoting Phase = "voting"
	PhaseEnding Phase = "ending"
)
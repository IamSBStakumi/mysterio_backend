package service

import (
	"errors"
	"log"
	"sync"

	"fmt"

	"github.com/IamSBStakumi/mysterio_backend/internal/domain"
)

type SessionService struct {
	mu        sync.Mutex
	sessions  map[string]*domain.Session
	scenarioS *ScenarioService
}

func NewSessionService(scenarioS *ScenarioService) *SessionService {
	return &SessionService{
		sessions:  make(map[string]*domain.Session),
		scenarioS: scenarioS,
	}
}

func (s *SessionService) CreateSession(
	playerCount int,
	difficulty string,
) (*domain.Session, error) {

	scenario, err := s.scenarioS.Generate(nil, playerCount, difficulty)
	if err != nil {
		return nil, err
	}

	session := &domain.Session{
		ID:       "session_1", // TODO: UUID
		Phase:    domain.PhaseIntro,
		Scenario: scenario,
		Players:  make(map[string]*domain.Player),
	}

	for i := 0; i < scenario.Meta.PlayerCount; i++ {
		playerID := fmt.Sprintf("player_%d", i+1) // TODO: UUID


		session.Players[playerID] = &domain.Player{
			ID:playerID,
			RoleID:fmt.Sprintf("p%d", i+1),
		}
	}

	if len(session.Players) == 0 {
	return nil, errors.New("no players created from scenario")
}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[session.ID] = session

	log.Printf("scenario title=%s phaseCount=%d",
		session.Scenario.Meta.Title,
		len(session.Scenario.Phases),
	)

	return session, nil
}

func (s *SessionService) JoinPlayer(
	sessionID string,
	playerName string,
) (*domain.Player, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[sessionID]
	if !ok {
		return nil, errors.New("session not found")
	}

	playerID := "player_" + playerName // TODO: UUID
	roleID := "p" + fmt.Sprint(len(session.Players)+1+'0')

	player := &domain.Player{
		ID:     playerID,
		RoleID: roleID,
	}

	session.Players[playerID] = player
	return player, nil
}

func (s *SessionService) GetPhase(
	sessionID string,
	playerID string,
) (domain.Phase, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[sessionID]
	if !ok {
		return "", errors.New("session not found")
	}

	if _, ok := session.Players[playerID]; !ok {
		return "", errors.New("player not found")
	}

	return session.Phase, nil
}

func (s *SessionService) AdvancePhase(sessionID string) (domain.Phase, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[sessionID]
	if !ok {
		return "", errors.New("session not found")
	}

	for i, p := range domain.PhaseOrder {
		if p == session.Phase && i+1 < len(domain.PhaseOrder) {
			session.Phase = domain.PhaseOrder[i+1]
			return session.Phase, nil
		}
	}

	return session.Phase, nil
}

func (s *SessionService) GetCurrentPhaseView(
	sessionID string,
	playerID string,
)(*domain.PhaseView, error){

	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[sessionID]
	if !ok {
		return nil, errors.New("session not found")
	}

	for _, playerLog := range session.Players {
		log.Printf(
			"[session=%s] playerID=%s",
			session.ID,
			playerLog.ID,
		)
	}

	player, ok := session.Players[playerID]
	if !ok {
		return nil, errors.New("player not found")
	}

	// 現在フェーズを取得
	currentPhase := session.Phase

	// フェーズの内容を取得
	var phaseContent *domain.PhaseContent
	for _, p := range session.Scenario.Phases {
		if p.Phase == currentPhase {
			phaseContent = &p
			break
		}
	}

	if phaseContent == nil {
		return nil, errors.New("phase content not found")
	}

	// 公開情報は全員に渡す
	view := &domain.PhaseView{
		Phase: currentPhase,
		Description: phaseContent.Public.Description,
		Actions: phaseContent.Public.Actions,
	}

	// 非公開情報はロールが一致するものだけ渡す
	for _, priv := range phaseContent.Private {
		if priv.RoleID == player.RoleID {
			view.Hints = append(view.Hints, priv.Hints...)
		}
	}

	return view, nil
}

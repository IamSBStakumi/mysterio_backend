package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/IamSBStakumi/mysterio_backend/internal/ai"
	"github.com/IamSBStakumi/mysterio_backend/internal/domain"
	"github.com/IamSBStakumi/mysterio_backend/internal/repository"
)

var (
	ErrSessionNotFound    = errors.New("session not found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidPhase       = errors.New("invalid phase")
	ErrAlreadyVoted       = errors.New("already voted")
	ErrGameNotFinished    = errors.New("game not finished")
	ErrInvalidPlayerCount = errors.New("invalid player count")
)

// SessionService はセッション管理のビジネスロジックを担当する
type SessionService struct {
	repo              repository.SessionRepository
	scenarioGenerator ai.ScenarioGenerator
}

// NewSessionService は新しいSessionServiceを作成する
func NewSessionService(
	repo repository.SessionRepository,
	scenarioGenerator ai.ScenarioGenerator,
) *SessionService {
	return &SessionService{
		repo:              repo,
		scenarioGenerator: scenarioGenerator,
	}
}

// CreateSession は新しいセッションを作成する
func (s *SessionService) CreateSession(
	ctx context.Context,
	playerCount int,
	difficulty string,
) (*domain.Session, error) {
	// バリデーション
	if playerCount < 2 || playerCount > 5 {
		return nil, ErrInvalidPlayerCount
	}

	diff := domain.Difficulty(difficulty)
	if diff != domain.DifficultyEasy &&
		diff != domain.DifficultyMedium &&
		diff != domain.DifficultyHard {
		return nil, errors.New("invalid difficulty")
	}

	// シナリオ生成
	scenario, err := s.scenarioGenerator.Generate(ctx, playerCount, diff)
	if err != nil {
		return nil, err
	}

	// プレイヤーIDを生成してキャラクターに割り当て
	playerIDs := make([]string, playerCount)
	for i := 0; i < playerCount; i++ {
		playerID := uuid.New().String()
		playerIDs[i] = playerID
		if i < len(scenario.Characters) {
			scenario.Characters[i].PlayerID = playerID
		}
	}

	// セッション作成
	session := &domain.Session{
		ID:            uuid.New().String(),
		OwnerPlayerID: playerIDs[0], // 最初のプレイヤーをオーナーに
		PlayerIDs:     playerIDs,
		Scenario:      scenario,
		CurrentPhase:  0,
		Votes:         make(map[string]string),
		Status:        domain.StatusActive,
		CreatedAt:     time.Now(),
	}

	// 保存
	if err := s.repo.Save(session); err != nil {
		return nil, err
	}

	return session, nil
}

// GetSession はセッションを取得する
func (s *SessionService) GetSession(ctx context.Context, sessionID string) (*domain.Session, error) {
	session, err := s.repo.FindByID(sessionID)
	if err != nil {
		if errors.Is(err, repository.ErrSessionNotFound) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}
	return session, nil
}

// GetPhaseInfo は現在のフェーズ情報を取得する
func (s *SessionService) GetPhaseInfo(
	ctx context.Context,
	sessionID string,
	playerID string,
) (*PhaseInfo, error) {
	session, err := s.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if session.CurrentPhase >= len(session.Scenario.Phases) {
		return nil, ErrInvalidPhase
	}

	phase := session.Scenario.Phases[session.CurrentPhase]
	character := session.Scenario.GetCharacterByPlayerID(playerID)

	privateText := ""
	if character != nil {
		if session.CurrentPhase == 0 {
			// 最初のフェーズでは役割と秘密情報を表示
			privateText = "【あなたの役割】\n" +
				"名前: " + character.Name + "\n" +
				"役割: " + character.Role + "\n\n" +
				"【あなただけが知る情報】\n" +
				character.SecretInfo
		} else {
			// その後のフェーズでは簡潔な情報のみ
			privateText = "あなたは " + character.Name + " です。"
		}
	}

	availableActions := []string{}
	if phase.Type == domain.PhaseTypeVoting && !session.HasVoted(playerID) {
		availableActions = append(availableActions, "vote")
	}
	if session.IsOwner(playerID) && !session.IsLastPhase() {
		availableActions = append(availableActions, "advance_phase")
	}

	return &PhaseInfo{
		PhaseNumber:      phase.Number,
		PhaseType:        string(phase.Type),
		Description:      phase.Description,
		PublicText:       phase.PublicText,
		PrivateText:      privateText,
		AvailableActions: availableActions,
	}, nil
}

// AdvancePhase は次のフェーズに進む
func (s *SessionService) AdvancePhase(
	ctx context.Context,
	sessionID string,
	playerID string,
) error {
	session, err := s.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	// オーナー権限チェック
	if !session.IsOwner(playerID) {
		return ErrUnauthorized
	}

	// 最終フェーズチェック
	if session.IsLastPhase() {
		session.Status = domain.StatusFinished
		return s.repo.Update(session)
	}

	// フェーズを進める
	session.CurrentPhase++

	return s.repo.Update(session)
}

// Vote は投票を記録する
func (s *SessionService) Vote(
	ctx context.Context,
	sessionID string,
	playerID string,
	targetPlayerID string,
) error {
	session, err := s.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	// 投票フェーズチェック
	if !session.IsVotingPhase() {
		return ErrInvalidPhase
	}

	// 既に投票済みチェック
	if session.HasVoted(playerID) {
		return ErrAlreadyVoted
	}

	// 投票対象が有効なプレイヤーかチェック
	validTarget := false
	for _, pid := range session.PlayerIDs {
		if pid == targetPlayerID {
			validTarget = true
			break
		}
	}
	if !validTarget {
		return errors.New("invalid target player")
	}

	// 投票を記録
	session.Votes[playerID] = targetPlayerID

	return s.repo.Update(session)
}

// GetResult はゲーム結果を取得する
func (s *SessionService) GetResult(
	ctx context.Context,
	sessionID string,
) (*GameResult, error) {
	session, err := s.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// 最終フェーズ以降のみ結果を公開
	if !session.IsLastPhase() && session.CurrentPhase < len(session.Scenario.Phases)-1 {
		return nil, ErrGameNotFinished
	}

	winner := session.GetWinner()

	// 投票結果を集計
	votingResult := make(map[string]int)
	for _, targetID := range session.Votes {
		votingResult[targetID]++
	}

	// 各プレイヤーの名前を取得
	playerNames := make(map[string]string)
	for _, char := range session.Scenario.Characters {
		playerNames[char.PlayerID] = char.Name
	}

	return &GameResult{
		Truth:         session.Scenario.Truth,
		WinnerID:      winner,
		WinnerName:    playerNames[winner],
		VotingResult:  &votingResult,
		PlayerNames:   &playerNames,
	}, nil
}

// PhaseInfo はフェーズ情報を表す
type PhaseInfo struct {
	PhaseNumber      int      `json:"phaseNumber"`
	PhaseType        string   `json:"phaseType"`
	Description      string   `json:"description"`
	PublicText       string   `json:"publicText"`
	PrivateText      string   `json:"privateText"`
	AvailableActions []string `json:"availableActions"`
}

// GameResult はゲーム結果を表す
type GameResult struct {
	Truth        string         `json:"truth"`
	WinnerID     string         `json:"winnerId"`
	WinnerName   string         `json:"winnerName"`
	VotingResult *map[string]int `json:"votingResult"`
	PlayerNames  *map[string]string `json:"playerNames"`
}

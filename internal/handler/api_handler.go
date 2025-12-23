package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/IamSBStakumi/mysterio_backend/internal/api"
	"github.com/IamSBStakumi/mysterio_backend/internal/domain"
	"github.com/IamSBStakumi/mysterio_backend/internal/service"
)

// APIHandler はOpenAPI仕様に基づくハンドラ
type APIHandler struct {
	sessionService *service.SessionService
}

// NewAPIHandler は新しいAPIHandlerを作成する
func NewAPIHandler(sessionService *service.SessionService) *APIHandler {
	return &APIHandler{
		sessionService: sessionService,
	}
}

// CreateSession はセッションを作成する (POST /sessions)
func (h *APIHandler) CreateSession(ctx echo.Context) error {
	var req api.CreateSessionRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Message: "invalid request body",
		})
	}

	session, err := h.sessionService.CreateSession(
		ctx.Request().Context(),
		req.PlayerCount,
		string(req.Difficulty),
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidPlayerCount) {
			return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
				Message: err.Error(),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Message: "failed to create session",
		})
	}

	// 初期フェーズ情報を構築
	phase := session.Scenario.Phases[0]
	initialPhase := api.PhaseInfo{
		PhaseNumber: phase.Number,
		PhaseType:   convertPhaseType(phase.Type),
		Description: phase.Description,
		PublicText:  phase.PublicText,
		Duration:    phase.Duration,
	}

	response := api.CreateSessionResponse{
		SessionId:     session.ID,
		OwnerPlayerId: session.OwnerPlayerID,
		PlayerIds:     session.PlayerIDs,
		InitialPhase:  initialPhase,
	}

	return ctx.JSON(http.StatusCreated, response)
}

// GetPhase は現在のフェーズ情報を取得する (GET /sessions/{sessionId}/phase)
func (h *APIHandler) GetPhase(ctx echo.Context, sessionId string, params api.GetPhaseParams) error {
	playerID := params.XPlayerId

	phaseInfo, err := h.sessionService.GetPhaseInfo(
		ctx.Request().Context(),
		sessionId,
		playerID,
	)
	if err != nil {
		return handleServiceError(ctx, err)
	}

	response := api.PhaseResponse{
		PhaseNumber:      phaseInfo.PhaseNumber,
		PhaseType:        convertPhaseTypeToResponse(phaseInfo.PhaseType),
		Description:      phaseInfo.Description,
		PublicText:       phaseInfo.PublicText,
		PrivateText:      phaseInfo.PrivateText,
		AvailableActions: phaseInfo.AvailableActions,
	}

	return ctx.JSON(http.StatusOK, response)
}

// AdvancePhase は次のフェーズに進める (POST /sessions/{sessionId}/phase/advance)
func (h *APIHandler) AdvancePhase(ctx echo.Context, sessionId string, params api.AdvancePhaseParams) error {
	playerID := params.XPlayerId

	err := h.sessionService.AdvancePhase(
		ctx.Request().Context(),
		sessionId,
		playerID,
	)
	if err != nil {
		return handleServiceError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
		Message: "phase advanced successfully",
	})
}

// Vote は投票を記録する (POST /sessions/{sessionId}/vote)
func (h *APIHandler) Vote(ctx echo.Context, sessionId string, params api.VoteParams) error {
	playerID := params.XPlayerId

	var req api.VoteRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Message: "invalid request body",
		})
	}

	err := h.sessionService.Vote(
		ctx.Request().Context(),
		sessionId,
		playerID,
		req.TargetPlayerId,
	)
	if err != nil {
		return handleServiceError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
		Message: "vote recorded successfully",
	})
}

// GetResult はゲーム結果を取得する (GET /sessions/{sessionId}/result)
func (h *APIHandler) GetResult(ctx echo.Context, sessionId string) error {
	result, err := h.sessionService.GetResult(
		ctx.Request().Context(),
		sessionId,
	)
	if err != nil {
		return handleServiceError(ctx, err)
	}

	response := api.ResultResponse{
		Truth:        result.Truth,
		WinnerId:     result.WinnerID,
		WinnerName:   result.WinnerName,
		VotingResult: *result.VotingResult,
		PlayerNames:  *result.PlayerNames,
	}

	return ctx.JSON(http.StatusOK, response)
}

// handleServiceError はサービス層のエラーを適切なHTTPレスポンスに変換する
func handleServiceError(ctx echo.Context, err error) error {
	switch {
	case errors.Is(err, service.ErrSessionNotFound):
		return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
			Message: "session not found",
		})
	case errors.Is(err, service.ErrUnauthorized):
		return ctx.JSON(http.StatusForbidden, api.ErrorResponse{
			Message: "unauthorized",
		})
	case errors.Is(err, service.ErrInvalidPhase):
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Message: "invalid phase",
		})
	case errors.Is(err, service.ErrAlreadyVoted):
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Message: "already voted",
		})
	case errors.Is(err, service.ErrGameNotFinished):
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Message: "game not finished yet",
		})
	case errors.Is(err, service.ErrInvalidPlayerCount):
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Message: "invalid player count",
		})
	default:
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Message: "internal server error",
		})
	}
}

// convertPhaseType はdomain.PhaseTypeをapi.PhaseInfoPhaseTypeに変換する
func convertPhaseType(phaseType domain.PhaseType) api.PhaseInfoPhaseType {
	switch phaseType {
	case domain.PhaseTypeIntroduction:
		return api.PhaseInfoPhaseTypeIntroduction
	case domain.PhaseTypeDiscussion:
		return api.PhaseInfoPhaseTypeDiscussion
	case domain.PhaseInvestigation1:
		return api.PhaseInfoPhaseTypeInvestigation
	case domain.PhaseTypeVoting:
		return api.PhaseInfoPhaseTypeVoting
	case domain.PhaseTypeReveal:
		return api.PhaseInfoPhaseTypeReveal
	default:
		return api.PhaseInfoPhaseTypeIntroduction
	}
}

// convertPhaseTypeToResponse はstring型のphaseTypeをapi.PhaseResponsePhaseTypeに変換する
func convertPhaseTypeToResponse(phaseType string) api.PhaseResponsePhaseType {
	switch phaseType {
	case "introduction":
		return api.PhaseResponsePhaseTypeIntroduction
	case "discussion":
		return api.PhaseResponsePhaseTypeDiscussion
	case "investigation":
		return api.PhaseResponsePhaseTypeInvestigation
	case "voting":
		return api.PhaseResponsePhaseTypeVoting
	case "reveal":
		return api.PhaseResponsePhaseTypeReveal
	default:
		return api.PhaseResponsePhaseTypeIntroduction
	}
}

// Ensure APIHandler implements ServerInterface
var _ api.ServerInterface = (*APIHandler)(nil)

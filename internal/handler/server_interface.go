package handler

import (
	"github.com/IamSBStakumi/mysterio_backend/internal/api"
	"github.com/IamSBStakumi/mysterio_backend/internal/service"
	"github.com/labstack/echo/v4"
)

type Server struct{
	SessionS *service.SessionService
}

type ServerInterface interface {
	// Create a new game session
	// (POST /sessions)
	PostSessions(ctx echo.Context) error
	// Advance to the next phase
	// (POST /sessions/{sessionId}/advance)
	PostSessionAdvance(ctx echo.Context, sessionId string) error
	// Get current phase info
	// (GET /sessions/{sessionId}/phase)
	GetSessionPhase(ctx echo.Context, sessionId string, params api.GetSessionPhaseParams) error
	// Join Session
	// (POST /sessions/{sessionId}/players)
	PostSessionPlayers(ctx echo.Context, sessionId string) error
}

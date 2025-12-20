package handler

import (
	"net/http"

	"github.com/IamSBStakumi/mysterio_backend/internal/api"
	"github.com/labstack/echo/v4"
)

// GET /sessions/{sessionId}/phase
func (s *Server) GetSessionPhase(c echo.Context, sessionId string, params api.GetSessionPhaseParams) error {

	// 1. HTTP レイヤーの値を取得
	playerID := c.Request().Header.Get("X-Player-Id")
	if playerID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "playerId is required"})
	}

	// 2. サービス層の値を取得
	view, err := s.SessionS.GetCurrentPhaseView(sessionId, playerID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, view)
}

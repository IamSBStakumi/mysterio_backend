package handler

import (
	"net/http"

	"github.com/IamSBStakumi/mysterio_backend/internal/api"
	"github.com/labstack/echo/v4"
)

// GET /sessions/{sessionId}/phase
func (s *Server) GetSessionPhase(c echo.Context, sessionId string, params api.GetSessionPhaseParams) error {
	resp := api.PhaseResponse{
		Phase:  api.PhaseResponsePhaseIntro,
		GmText: "ゲームを開始します。",
	}

	return c.JSON(http.StatusOK, resp)
}

package handler

import (
	"net/http"

	"github.com/IamSBStakumi/mysterio_backend/internal/api"
	"github.com/labstack/echo/v4"
)

// POST /sessions/{sessionId}/advance
func (s *Server) PostSessionAdvance(c echo.Context, sessionId string) error {
	resp := api.AdvancePhaseResponse{
		Phase: api.AdvancePhaseResponsePhaseIntro,
	}

	return c.JSON(http.StatusOK, resp)
}

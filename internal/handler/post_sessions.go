package handler

import (
	"net/http"

	"github.com/IamSBStakumi/mysterio_backend/internal/api"
	"github.com/labstack/echo/v4"
)

func(s *Server) PostSessions(c echo.Context) error {
	var req api.CreateSessionRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	sessionID := "uuid"

	return c.JSON(http.StatusOK, api.CreateSessionResponse{
		SessionId: sessionID,
	})
}

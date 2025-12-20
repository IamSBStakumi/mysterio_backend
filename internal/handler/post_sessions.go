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

	session, err := s.SessionS.CreateSession(
		int(req.PlayerCount),
		string(req.Difficulty),
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, api.CreateSessionResponse{
		SessionId: session.ID,
	})
}

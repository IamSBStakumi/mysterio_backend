package handler

import (
	"net/http"

	"github.com/IamSBStakumi/mysterio_backend/internal/api"
	"github.com/labstack/echo/v4"
)

func(s *Server) PostSessionPlayers(c echo.Context, sessionId string) error {
	var req api.JoinPlayerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	playerId := "player_dummy"
	roleId := "p1"

	resp := api.JoinPlayerResponse{
		PlayerId: playerId,
		RoleId:   roleId,
	}

	return c.JSON(http.StatusOK, resp)
}

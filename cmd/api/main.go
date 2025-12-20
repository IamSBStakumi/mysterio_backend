package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/IamSBStakumi/mysterio_backend/internal/api"
	"github.com/IamSBStakumi/mysterio_backend/internal/handler"
	"github.com/IamSBStakumi/mysterio_backend/internal/service"
)

func main(){
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	scenarioS, err := service.NewScenarioService()
	if err != nil {
		log.Fatal(err)
	}

	sessionS := service.NewSessionService(scenarioS)

	server := &handler.Server{
		SessionS: sessionS,
	}
	api.RegisterHandlers(e, server)

	e.Logger.Fatal(e.Start(":8080"))
}

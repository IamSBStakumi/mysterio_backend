package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/IamSBStakumi/mysterio_backend/internal/api"
	"github.com/IamSBStakumi/mysterio_backend/internal/handler"
)

func main(){
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	server := &handler.Server{}
	api.RegisterHandlers(e, server)

	e.Logger.Fatal(e.Start(":8080"))
}

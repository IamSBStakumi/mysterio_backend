package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/IamSBStakumi/mysterio_backend/internal/ai"
	"github.com/IamSBStakumi/mysterio_backend/internal/api"
	"github.com/IamSBStakumi/mysterio_backend/internal/handler"
	"github.com/IamSBStakumi/mysterio_backend/internal/repository"
	"github.com/IamSBStakumi/mysterio_backend/internal/service"
)

func main() {
	// Echo インスタンスを作成
	e := echo.New()

	// ミドルウェアを設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// 依存関係を初期化
	sessionRepo := repository.NewInMemorySessionRepository()
	scenarioGenerator := ai.NewMockScenarioGenerator()
	sessionService := service.NewSessionService(sessionRepo, scenarioGenerator)

	// OpenAPI準拠のハンドラを作成
	apiHandler := handler.NewAPIHandler(sessionService)

	// ヘルスチェック
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
		})
	})

	// OpenAPI仕様に基づくルーティングを登録
	apiGroup := e.Group("/api/v1")
	api.RegisterHandlers(apiGroup, apiHandler)

	// サーバーを起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

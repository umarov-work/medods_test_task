// @title           Medods Test Task
// @version         1.0
// @description     Тестовое задание Medods.
// @host            localhost:8080
// @BasePath        /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	_ "medods_test_task/docs"
	"medods_test_task/internal/config"
	db "medods_test_task/internal/db/impl"
	"medods_test_task/internal/handler"
	"medods_test_task/internal/model"
	repo "medods_test_task/internal/repository/impl"
	service "medods_test_task/internal/service/impl"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	cfg := config.Load()
	database := db.NewPostgresDB()
	if err := database.Connect(cfg.DbDsn); err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Fatalf("failed to close DB: %v", err)
		}
	}()

	if err := database.Migrate(&model.RefreshToken{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	refreshTokenRepository := repo.NewRefreshTokenRepository(database.DB())
	authService := service.NewAuthService(refreshTokenRepository)
	authHandler := handler.NewAuthHandler(authService)

	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := router.Group("/api")

	authHandler.RegisterAuthHandlers(api)
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

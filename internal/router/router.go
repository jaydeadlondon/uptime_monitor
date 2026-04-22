package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jaydeadlondon/project_na_go/internal/config"
	"github.com/jaydeadlondon/project_na_go/internal/handlers"
	"github.com/jaydeadlondon/project_na_go/internal/middleware"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	authHandler := handlers.NewAuthHandler(db, cfg)

	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	protected := api.Group("", middleware.Protected(cfg.JWT.Secret))
	protected.Get("/auth/me", authHandler.Me)
}
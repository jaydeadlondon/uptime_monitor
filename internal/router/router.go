package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jaydeadlondon/project_na_go/internal/config"
	"github.com/jaydeadlondon/project_na_go/internal/handlers"
	"github.com/jaydeadlondon/project_na_go/internal/middleware"
	"github.com/jaydeadlondon/project_na_go/internal/scheduler"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB, cfg *config.Config, s *scheduler.Scheduler) {
	authHandler := handlers.NewAuthHandler(db, cfg)
	monitorHandler := handlers.NewMonitorHandler(db, s)

	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	protected := api.Group("", middleware.Protected(cfg.JWT.Secret))
	protected.Get("/auth/me", authHandler.Me)

	monitors := protected.Group("/monitors")
	monitors.Post("/", monitorHandler.Create)
	monitors.Get("/", monitorHandler.GetAll)
	monitors.Get("/:id", monitorHandler.GetOne)
	monitors.Put("/:id", monitorHandler.Update)
	monitors.Delete("/:id", monitorHandler.Delete)
	monitors.Get("/:id/checks", monitorHandler.GetChecks)
	monitors.Get("/:id/stats", monitorHandler.GetStats)
}
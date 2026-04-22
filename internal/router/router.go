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
	telegramHandler := handlers.NewTelegramHandler(db, cfg.Telegram.BotName)
	analyticsHandler := handlers.NewAnalyticsHandler(db)
	statusPageHandler := handlers.NewStatusPageHandler(db)

	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	app.Get("/status/:userID", statusPageHandler.GetStatusPage)

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

	analytics := protected.Group("/analytics")
	analytics.Get("/overview", analyticsHandler.Overview)
	analytics.Get("/monitors/:id/chart", analyticsHandler.MonitorChart)
	analytics.Get("/monitors/:id/incidents", analyticsHandler.MonitorIncidents)

	tg := protected.Group("/telegram")
	tg.Post("/link", telegramHandler.GenerateLink)
	tg.Delete("/unlink", telegramHandler.Unlink)
	tg.Get("/status", telegramHandler.Status)
}
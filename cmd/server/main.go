package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jaydeadlondon/project_na_go/internal/config"
	"github.com/jaydeadlondon/project_na_go/internal/database"
	"github.com/jaydeadlondon/project_na_go/internal/router"
	"github.com/jaydeadlondon/project_na_go/internal/scheduler"
	"github.com/jaydeadlondon/project_na_go/internal/telegram"
)

func main() {
	cfg := config.Load()

	db := database.NewPostgresDB(cfg)
	database.Migrate(db)

	rdb := database.NewRedisClient(cfg)
	_ = rdb

	var tgBot *telegram.Bot
	if cfg.Telegram.BotToken != "" {
		var err error
		tgBot, err = telegram.NewBot(cfg.Telegram.BotToken, db)
		if err != nil {
			log.Printf("⚠️  Telegram bot failed to start: %v", err)
		} else {
			tgBot.Start()
			defer tgBot.Stop()
		}
	} else {
		log.Println("⚠️  Telegram bot token not set, notifications disabled")
	}

	s := scheduler.NewScheduler(db, tgBot)
	s.Start()
	defer s.Stop()

	app := fiber.New(fiber.Config{
		AppName: "Uptime Monitor API",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, PATCH",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"app":    "Uptime Monitor",
		})
	})

	router.Setup(app, db, cfg, s)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", cfg.App.Port)))
}
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
)

func main() {
	cfg := config.Load()

	db := database.NewPostgresDB(cfg)
	database.Migrate(db)

	rdb := database.NewRedisClient(cfg)
	_ = rdb

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


	log.Fatal(app.Listen(fmt.Sprintf(":%s", cfg.App.Port)))
}
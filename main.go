package main

import (
	"context"
	"log"

	"gocsms/config"
	"gocsms/db"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/mutoulbj/gocsms/config"
)

func main() {
	_ = godotenv.Load()
	config.Init()
	db.Init()

	app := fiber.New()

	app.Get("/hello", func(c *fiber.Ctx) error {
		ctx := context.Background()
		err := db.BunDB.PingContext(ctx)
		if err != nil {
			return c.Status(500).SendString("Database not reachable")
		}
		return c.SendString("world")
	})

	log.Fatal(app.Listen(":" + config.Cfg.Port))
}

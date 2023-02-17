package main

import (
	"github.com/cityos-dev/Qiao-Lin/database"
	"github.com/gofiber/fiber/v2"
)

func main() {
	database.ConnectToPostgresDb()

	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // this is the default limit of 10MB
	})
	setupRoutes(app)

	app.Listen("0.0.0.0:8080")
}

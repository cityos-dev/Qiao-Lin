package main

import (
	"time"

	"fmt"

	"github.com/cityos-dev/Qiao-Lin/database"
	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("testing server ")
	database.ConnectToPostgresDb()

	app := fiber.New(fiber.Config{
		BodyLimit:    10 * 1024 * 1024, // this is the default limit of 10MB
		IdleTimeout:  20 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})
	setupRoutes(app)
	fmt.Println("server started ")
	app.Listen("0.0.0.0:8080")
}

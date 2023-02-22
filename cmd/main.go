package main

import (
	"time"

	"fmt"

	"github.com/cityos-dev/Qiao-Lin/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	fmt.Println("testing server ")
	database.ConnectToPostgresDb()

	app := fiber.New(fiber.Config{
		BodyLimit:    20 * 1024 * 1024, // this is the default limit of 10MB
		IdleTimeout:  20 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	setupRoutes(app)
	fmt.Println("server started ")
	app.Listen("0.0.0.0:8080")
}

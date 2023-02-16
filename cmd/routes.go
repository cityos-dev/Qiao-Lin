package main

import (
	"github.com/cityos-dev/Qiao-Lin/handlers"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {

	app.Post("/v1/files", handlers.UploadFile)

	app.Get("/v1/files", handlers.ListUploadedFiles)

	app.Delete("/v1/:fileid", handlers.DeleteOneFile)

	app.Get("/v1/files/:fileid", handlers.GetOneFile)

	app.Get("/v1/health", handlers.HealthCheck)
}

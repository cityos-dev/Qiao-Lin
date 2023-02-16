package main

import (
	"github.com/cityos-dev/Qiao-Lin/handlers"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {

	app.Post("/files", handlers.UploadFile)

	app.Get("/files", handlers.ListUploadedFiles)

	app.Delete("/:fileid", handlers.DeleteOneFile)

	app.Get("/files/:fileid", handlers.GetOneFile)
}

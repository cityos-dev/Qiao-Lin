package handlers

import (
	"fmt"
	"mime"
	"net"
	"os"
	"time"

	"errors"
	"log"

	"github.com/cityos-dev/Qiao-Lin/database"
	"github.com/cityos-dev/Qiao-Lin/models"
	"github.com/gofiber/fiber/v2"
)

func HealthCheck(c *fiber.Ctx) error {
	timeout := 1 * time.Second
	_, err := net.DialTimeout("http", "localhost:8080/v1", timeout)
	if err != nil {
		fmt.Println("Server unreachable, error: ", err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":        200,
		"description": "server is healthy",
	})
}

func UploadFile(c *fiber.Ctx) error {
	queryValue := c.Query("files")
	fmt.Println(queryValue)
	file, err := c.FormFile("data")

	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"code":        400,
			"description": "bad request",
		})
	}

	fileName := file.Filename
	contentType := file.Header["Content-Type"][0]
	fileSize := file.Size

	//check if media type if supported
	if contentType != "video/mpeg" && contentType != "video/mp4" {
		fmt.Println("Uploaded file media type is not supported!!")
		c.Response().Header.Set("Location", "")
		return c.Status(415).JSON(fiber.Map{
			"code":        415,
			"description": "media not supported",
		})
	}
	//check if file exists
	fmt.Println("check if uploaded file exists in database")
	var fileExist models.File
	database.DB.Db.Where("name = ?", fileName).Find(&fileExist)

	if fileExist.Name != "" {
		fmt.Println("File exists" + fileExist.Name)
		return c.Status(409).JSON(fiber.Map{
			"code":        409,
			"description": "file with same name exists",
		})
	}
	fmt.Println("File name does not exist in database")

	//save the file locally
	//create uploads folder if not exist
	createDirectoryIfNotExist()

	fmt.Println("Saving uploaded file")
	c.SaveFile(file, fmt.Sprintf("./uploads/%s", fileName))

	//save the file in a docker based postgres db

	newFile := &models.File{
		FileId:     fileName,
		Name:       fileName,
		Size:       fileSize,
		Created_At: time.Now(),
	}

	database.DB.Db.Create(&newFile)
	fileLocation := "./uploads/" + fileName
	c.Response().Header.Set("Location", fileLocation)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":        201,
		"description": "File uploaded",
	})
}

func ListUploadedFiles(c *fiber.Ctx) error {
	var apiFiles []models.File
	err := database.DB.Db.Model(&models.File{}).Find(&apiFiles)
	if err.Error != nil {
		return c.SendStatus(fiber.StatusNoContent)
	}
	return c.Status(fiber.StatusOK).JSON(apiFiles)
}

func DeleteOneFile(c *fiber.Ctx) error {
	name := c.Params("fileId")
	var fileExist models.File
	var file models.File
	database.DB.Db.Where("name = ?", name).Find(&fileExist)
	fmt.Println(fileExist.Name)
	if fileExist.Name == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":        404,
			"description": "File not found",
		})
	}
	fmt.Println("Now deleting")
	database.DB.Db.Where("name = ?", name).Delete(&file)
	path := "./uploads/" + fileExist.Name
	err := os.Remove(path)

	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":        500,
			"description": "server error",
		})
	}
	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"code":        204,
		"description": "File removed",
	})
}

func GetOneFile(c *fiber.Ctx) error {
	//check requested file exist or not
	name := c.Params("fileId")
	var fileExist models.File
	database.DB.Db.Where("name = ?", name).Find(&fileExist)
	if fileExist.Name == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":        404,
			"description": "File not found",
		})
	}
	fmt.Println(fileExist.Name)

	fileLocation := "./uploads/" + fileExist.Name
	c.Response().Header.Set("Content-Disposition", fmt.Sprintf("form-data; name='data'; filename=%s", name))
	c.Response().Header.Set("Content-Type", mime.TypeByExtension(fileLocation))
	return c.Status(fiber.StatusOK).SendFile(fileLocation)
}

func createDirectoryIfNotExist() {
	path := "./uploads"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("path %s created", path)
	}
}

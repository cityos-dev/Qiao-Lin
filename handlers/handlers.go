package handlers

import (
	"fmt"
	"net"
	"os"
	"time"

	"errors"
	"log"

	"github.com/cityos-dev/Qiao-Lin/database"
	"github.com/cityos-dev/Qiao-Lin/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func HealthCheck(c *fiber.Ctx) error {
	timeout := 1 * time.Second
	_, err := net.DialTimeout("http", "localhost:8080/v1", timeout)
	if err != nil {
		fmt.Println("Server unreachable, error: ", err)
	}
	return c.Status(200).JSON(fiber.Map{
		"code":        200,
		"description": "server is healthy",
	})
}

func UploadFile(c *fiber.Ctx) error {
	queryValue := c.Query("files")
	fmt.Println(queryValue)
	file, err := c.FormFile("files")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	buffer, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	defer buffer.Close()

	fileName := file.Filename
	fileBuffer := buffer
	contentType := file.Header["Content-Type"][0]
	fileSize := file.Size
	fileValue := c.Body()

	fmt.Println(fileName)
	fmt.Println(fileBuffer)
	fmt.Println(contentType)
	fmt.Println(fileSize)

	//check if media type if supported
	if contentType != "video/mpeg" && contentType != "video/mp4" {
		fmt.Println("Uploaded file media type is not supported!!")
		return c.Status(415).JSON("Unsupported Media Type")
	}
	//check if file exists
	fmt.Println("check if uploaded file exists in database")
	var fileExist models.File
	database.DB.Db.Where("name = ?", fileName).Find(&fileExist)

	if fileExist.Name != "" {
		fmt.Println("File exists" + fileExist.Name)
		return c.Status(409).JSON("File exists")
	}
	fmt.Println("File name does not exist in database")

	//save the file locally
	//create uploads folder if not exist
	createDirectoryIfNotExist()

	fmt.Println("Saving uploaded file")
	c.SaveFile(file, fmt.Sprintf("./uploads/%s", fileName))

	//save the file in a docker based postgres db

	newFile := &models.File{
		FileId:     uuid.New(),
		Name:       fileName,
		Size:       fileSize,
		Created_at: time.Now(),
		Content:    fileValue,
	}

	database.DB.Db.Create(&newFile)
	return c.Status(201).JSON("File uploaded")
}

func ListUploadedFiles(c *fiber.Ctx) error {
	type APIFiles struct {
		FileId     string
		Name       string
		Size       int64
		Created_at time.Time
	}

	var apiFiles []APIFiles
	err := database.DB.Db.Model(&models.File{}).Find(&apiFiles)
	if err.Error != nil {
		return c.SendStatus(fiber.StatusNoContent)
	}
	return c.Status(200).JSON(apiFiles)
}

func DeleteOneFile(c *fiber.Ctx) error {
	id := c.Params("fileid")
	var fileExist models.File
	var file models.File
	database.DB.Db.Where("file_id = ?", id).Find(&fileExist)
	fmt.Println(fileExist.Name)
	if fileExist.Name == "" {
		return c.Status(404).JSON("File not found")
	}
	fmt.Println("Now deleting")
	database.DB.Db.Where("file_id = ?", id).Delete(&file)
	path := "./uploads/" + fileExist.Name
	err := os.Remove(path)

	if err != nil {
		fmt.Println(err)
		return c.Status(500).JSON("Server error")
	}
	return c.Status(204).JSON("File was successfully removed")
}

func GetOneFile(c *fiber.Ctx) error {
	//check requested file exist or not
	id := c.Params("fileid")
	var fileExist models.File
	database.DB.Db.Where("file_id = ?", id).Find(&fileExist)
	if fileExist.Name == "" {
		return c.Status(404).JSON("File not found")
	}
	fmt.Println(fileExist.Name)

	fileLocation := "./uploads/" + fileExist.Name

	return c.Status(200).Download(fileLocation)
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

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
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":        200,
		"description": "server is healthy",
	})
}

func UploadFile(c *fiber.Ctx) error {
	c.Accepts("video/mp4")
	queryValue := c.Query("files")
	fmt.Println(queryValue)
	file, _ := c.FormFile("files")
	testFile, _ := c.Context().Request.MultipartForm()

	// var testFile models.File
	// if err := c.BodyParser(&testFile); err != nil {
	// 	fmt.Println("error = ", err)
	// 	return c.SendStatus(fiber.StatusInternalServerError)
	// }
	fmt.Printf("Get testFile info %+v", testFile)
	fmt.Println(testFile.Value)

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
		FileId: uuid.New(),
		Name:   fileName,
		Size:   fileSize,
	}

	database.DB.Db.Create(&newFile)
	fileLocation := "./uploads/" + fileExist.Name
	c.Response().Header.Set("Location", fileLocation)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":        201,
		"description": "File uploaded",
	})
}

func ListUploadedFiles(c *fiber.Ctx) error {
	// type APIFiles struct {
	// 	FileId     string
	// 	Name       string
	// 	Size       int64
	// 	Created_at time.Time
	// }

	var apiFiles []models.File
	err := database.DB.Db.Model(&models.File{}).Find(&apiFiles)
	if err.Error != nil {
		return c.SendStatus(fiber.StatusNoContent)
	}
	return c.Status(fiber.StatusOK).JSON(apiFiles)
}

func DeleteOneFile(c *fiber.Ctx) error {
	id := c.Params("fileid")
	var fileExist models.File
	var file models.File
	database.DB.Db.Where("file_id = ?", id).Find(&fileExist)
	fmt.Println(fileExist.Name)
	if fileExist.Name == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":        404,
			"description": "File not found",
		})
	}
	fmt.Println("Now deleting")
	database.DB.Db.Where("file_id = ?", id).Delete(&file)
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
	id := c.Params("fileid")
	var fileExist models.File
	database.DB.Db.Where("file_id = ?", id).Find(&fileExist)
	if fileExist.Name == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":        404,
			"description": "File not found",
		})
	}
	fmt.Println(fileExist.Name)

	fileLocation := "./uploads/" + fileExist.Name

	return c.Status(fiber.StatusOK).Download(fileLocation)
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

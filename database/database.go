package database

import (
	"fmt"
	"log"
	"os"

	"github.com/cityos-dev/Qiao-Lin/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var DB DbInstance

func ConnectToPostgresDb() {
	fmt.Println(os.Getenv("DB_USER"))
	fmt.Println(os.Getenv("DB_PASSWORD"))
	fmt.Println(os.Getenv("DB_NAME"))
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		"postgres",
		"5432",
		os.Getenv("DB_NAME"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to postgres-db. \n", err)
		os.Exit(2)
	}

	log.Println("postgres-db connected")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("running migrations")

	db.AutoMigrate(&models.File{})

	DB = DbInstance{
		Db: db,
	}
}

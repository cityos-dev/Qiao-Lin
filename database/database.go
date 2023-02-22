package database

import (
	"fmt"
	"log"
	"os"

	"github.com/cityos-dev/Qiao-Lin/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DbInstance struct {
	Db *gorm.DB
}

var DB DbInstance

func ConnectToPostgresDb() {
	fmt.Println(os.Getenv("DB_USERNAME"))
	fmt.Println(os.Getenv("DB_PASSWORD"))
	fmt.Println(os.Getenv("DB_DATABASE_NAME"))

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"postgres", 5432, os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DATABASE_NAME"))

	db, err := gorm.Open("postgres", dsn)

	if err != nil {
		log.Fatal("Failed to connect to postgres db. \n", err)
		panic(err)
	}

	if err := db.DB().Ping(); err != nil {
		log.Fatal("Failed to connect to postgres db. \n", err)
		panic(err)
	}

	log.Println("postgres db connected")

	log.Println("running migrations")

	db.AutoMigrate(&models.File{})

	DB = DbInstance{
		Db: db,
	}
}

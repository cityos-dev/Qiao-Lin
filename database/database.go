package database

import (
	"fmt"
	"log"

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
	dbconn := fmt.Sprintf(
		"host=db user=qilin password=qilin dbname=qilin port=5432 sslmode=disable TimeZone=Asia/Shanghai",
	)

	db, err := gorm.Open(postgres.Open(dbconn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to postgres-db. \n", err)
		// os.Exit(2)
	}

	log.Println("postgres-db connected")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("running migrations")

	db.AutoMigrate(&models.File{})

	DB = DbInstance{
		Db: db,
	}
}

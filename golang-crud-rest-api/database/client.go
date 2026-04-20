package database

import (
	"golang-crud-rest-api/entities"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Instance *gorm.DB

func Connect(connectionString string) {
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatalf("cannot connect to DB: %v", err)
	}
	Instance = db
}

func Migrate() {
	if err := Instance.AutoMigrate(&entities.Product{}); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}
	log.Println("Database migration completed")
}

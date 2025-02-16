package config

import (
	"fmt"
	"log"
	"os"
	"gorm.io/driver/postgres"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var DB *gorm.DB;
func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using system env vars")
	}
	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Africa/Accra",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{});
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	DB = db
	fmt.Println("Database connected")
}
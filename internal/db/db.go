package db

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var DB *sql.DB

func Init() error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s  sslmode=disable",
		os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
	)

	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Println("Error opening database connection:", err)
		return err
	}
	err = DB.Ping()
	if err != nil {
		log.Println("Error pinging database connection:", err)
		return err
	}

	return nil
}

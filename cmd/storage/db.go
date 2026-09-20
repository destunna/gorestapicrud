package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")

	// sql.Open - функция подключения к DB
	db, err = sql.Open("postgres", fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbName, dbPort))

	if err != nil {
		panic(err.Error())
	}

	// db.Ping() - функция проверки успешности подключения к DB
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Successfully connected to database")
}

func GetDB() *sql.DB {
	return db
}

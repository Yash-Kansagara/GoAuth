package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // mysql driver import
)

var sqlDB *sql.DB

func ConnectDB() (*sql.DB, error) {

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	} else {
		log.Println("connected to DB")
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	} else {
		log.Println("Ping successful")
	}
	sqlDB = db
	return db, nil
}

func GetDB() *sql.DB {
	return sqlDB
}

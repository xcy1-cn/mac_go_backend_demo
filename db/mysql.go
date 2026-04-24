package db

import (
	"database/sql"
	"fmt"
	"log"

	"demo/day6-9/config"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitMySQL() {
	user := config.GetEnv("DB_USER")
	password := config.GetEnv("DB_PASSWORD")
	host := config.GetEnv("DB_HOST")
	port := config.GetEnv("DB_PORT")
	dbName := config.GetEnv("DB_NAME")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName,
	)

	database, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("failed to connect mysql:", err)
	}

	if err = database.Ping(); err != nil {
		log.Fatal("failed to ping mysql:", err)
	}

	DB = database
	log.Println("mysql connected")
}

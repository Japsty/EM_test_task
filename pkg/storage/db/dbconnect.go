package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"
)

// ConnectToDb - функция создающее подключение с бд и предоставляющее его наружу
func ConnectToDb() (*sql.DB, error) {
	DBDriver := os.Getenv("DB_DRIVER")
	DBSource := os.Getenv("DB_SOURCE")

	db, err := sql.Open(DBDriver, DBSource)

	slog.Debug("Db connection opened")
	if err != nil {
		slog.Error("Failed to connect to db")
		fmt.Println("Error connecting to the database:", err)
		return nil, err
	}
	time.Sleep(2 * time.Second)

	err = db.Ping()
	if err != nil {
		slog.Error("Failed to ping db")
		fmt.Println("Error pinging the database:", err)
		return nil, err
	}
	return db, nil
}

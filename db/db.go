package db

import (
	"database/sql"
	"os"

	"github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {
	cfg := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASS"),
		Net:    "tcp",
		Addr:   os.Getenv("DB_ADDRESS"),
		DBName: os.Getenv("DB_NAME"),
	}

	return sql.Open("mysql", cfg.FormatDSN())
}

func GetAllData(db *sql.DB, limit int) (*sql.Rows, error) {

}

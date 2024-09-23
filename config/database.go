package config

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sync"
	"time"
)

var (
	DBS  *sql.DB
	once sync.Once
)

// GetMySQLDB returns a singleton instance of the database connection.
func GetMySQLDB() (*sql.DB, error) {
	var err error
	once.Do(func() {
		dsn := DSN
		var db *sql.DB
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Error opening database: %v", err)
		}

		db.SetMaxOpenConns(100)
		db.SetMaxIdleConns(50)
		db.SetConnMaxLifetime(5 * time.Minute)
		// Test the connection
		err = db.Ping()
		if err != nil {
			log.Fatalf("Error connecting to the database: %v", err)
		}

		DBS = db
	})

	return DBS, err
}

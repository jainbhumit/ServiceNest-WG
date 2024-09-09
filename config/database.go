package config

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sync"
)

var (
	DBS  *sql.DB
	once sync.Once
)

// GetMySQLDB returns a singleton instance of the database connection.
func GetMySQLDB() *sql.DB {
	once.Do(func() {
		dsn := "root:Asdfghjkl@0987@tcp(localhost:3306)/servicenest"
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Error opening database: %v", err)
		}

		// Test the connection
		err = db.Ping()
		if err != nil {
			log.Fatalf("Error connecting to the database: %v", err)
		}

		DBS = db
	})

	return DBS
}

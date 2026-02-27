package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"p2p-management-service/config"
	"p2p-management-service/scripts"
)

var Conn *sql.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.AppConfig.DB.User,
		config.AppConfig.DB.Password,
		config.AppConfig.DB.Host,
		config.AppConfig.DB.Port,
		config.AppConfig.DB.Name,
		config.AppConfig.DB.SSLMode,
	)

	var err error
	Conn, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err = Conn.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// RunMigrations(dsn, "./db/migrations")

	log.Println("Database connection established!")

	// Run Migrations
	scripts.RunMigrations(dsn)
}

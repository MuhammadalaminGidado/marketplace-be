// cmd/migrate/main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	db "example/api/internal/database"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func createDBIfNotExists(maintenanceDSN, dbname string) error {
	sqlDB, err := sql.Open("postgres", maintenanceDSN)
	if err != nil {
		return fmt.Errorf("connect to maintenance db: %w", err)
	}
	defer sqlDB.Close()

	var exists bool
	if err := sqlDB.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbname).Scan(&exists); err != nil {
		return fmt.Errorf("check db existence: %w", err)
	}

	if !exists {
		if _, err := sqlDB.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, dbname)); err != nil {
			return fmt.Errorf("create db: %w", err)
		}
		log.Printf("database %q created", dbname)
	}

	return nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found")
	}

	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <up|down|force> [version]")
	}

	if err := createDBIfNotExists(os.Getenv("MAINTENANCE_DSN"), os.Getenv("DB_NAME")); err != nil {
		log.Fatalf("ensure db: %v", err)
	}

	database, err := db.Init(db.Config{
		DSN:             os.Getenv("DATABASE_URL"),
		Env:             os.Getenv("APP_ENV"),
		MaxIdleConns:    10,
		MaxOpenConns:    25,
		ConnMaxLifetime: 30 * time.Minute,
	})
	if err != nil {
		log.Fatalf("db: %v", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		log.Fatalf("unwrap sql.DB: %v", err)
	}

	if err := db.RunMigrations(sqlDB, "file://db/migrations", os.Args[1], os.Args[2:]...); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	log.Printf("migration %s successful", os.Args[1])
}

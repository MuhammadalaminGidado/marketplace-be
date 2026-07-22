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

	// Parse command line arguments
	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <up|down|force> [version]")
	}

	// Set environment variables to force migrations
	os.Setenv("RUN_MIGRATIONS", "true")

	// Set direction from command line
	if len(os.Args) >= 2 {
		os.Setenv("MIGRATION_DIRECTION", os.Args[1])
	}

	// Set force version if provided
	if len(os.Args) >= 3 && os.Args[1] == "force" {
		os.Setenv("MIGRATION_FORCE_VERSION", os.Args[2])
	}

	// Set steps for down migration if provided
	if len(os.Args) >= 3 && os.Args[1] == "down" {
		os.Setenv("MIGRATION_STEPS", os.Args[2])
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

	// Run migrations - this will use the environment variables we set above
	// Pass empty string for direction to use environment variable
	if err := db.RunMigrations(sqlDB, "file://db/migrations", ""); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	log.Printf("migration %s successful", os.Args[1])
}

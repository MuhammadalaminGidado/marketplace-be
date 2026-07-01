package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DSN             string
	Env             string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

func Init(cfg Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logLevelFor(cfg.Env)),
	})
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("unwrap sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return db, nil
}

func logLevelFor(env string) logger.LogLevel {
	if env == "production" {
		return logger.Warn
	}
	return logger.Info
}

func createDBIfNotExists(maintenanceDSN, dbname string) error {
	sqlDB, err := sql.Open("postgres", maintenanceDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to maintenance database: %w", err)
	}
	defer sqlDB.Close()

	var exists bool
	if err := sqlDB.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbname).Scan(&exists); err != nil {
		return fmt.Errorf("failed to check database existence: %w", err)
	}

	if !exists {
		if _, err := sqlDB.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, dbname)); err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
	}

	return nil
}

package database

import (
	"errors"
	"fmt"
	"go-rest-api/config"
	"go-rest-api/internal/infra/logger"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(cfg config.Configuration) error {
	if cfg.MigrateToVersion == "" {
		logger.Logger.Error("Empty migration version")
		return nil
	}

	migrationsPath := cfg.MigrationLocation

	_, err := os.Stat(migrationsPath)
	if err != nil {
		logger.Logger.Error(err)
		return err
	}

	urlString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseName,
	)

	m, err := migrate.New(
		"file://"+migrationsPath,
		urlString)
	if err != nil {
		logger.Logger.Error(err)
		return err
	}

	dbVersion, err := strconv.Atoi(cfg.MigrateToVersion)
	if err == nil {
		logger.Logger.Info("Migrate: starting migration to version %v\n", dbVersion)
		err = m.Migrate(uint(dbVersion))
		if err != nil {
			logger.Logger.Error("Migrate: failed migration to version %v\n", dbVersion)
			logger.Logger.Error("Migration table will be forcing to version %v\n You should clean your data base from wrong tables and then start server mith 'MIGRATE=latest' enviroment variable!", dbVersion)
			err = m.Force(dbVersion)
		}
	} else {
		logger.Logger.Info("Migrate: starting migration to the latest version")
		err = m.Up()
	}
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Logger.Info("Migrate: no changes found")
			return nil
		}
		logger.Logger.Info("file://" + migrationsPath)
		logger.Logger.Error(err)
		return err
	}
	logger.Logger.Info("Migrate: migrations are done successfully")
	return nil
}

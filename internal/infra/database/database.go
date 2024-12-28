package database

import (
	"database/sql"
	"fmt"
	"go-rest-api/config"
	"go-rest-api/internal/infra/logger"

	_ "github.com/lib/pq"
)

type databaseManager struct {
	cfg config.Configuration
}

func New(cfg config.Configuration) *sql.DB {
	dbManager := databaseManager{
		cfg: cfg,
	}

	db, err := dbManager.newDatabase()
	if err != nil {
		logger.Logger.Panic(err)
		panic(err)
	}
	return db
}

func (dm databaseManager) newDatabase() (*sql.DB, error) {
	db, err := sql.Open("postgres", dm.getConnectionString())
	if err != nil {
		logger.Logger.Error(err)
		return nil, err
	}
	if err = db.Ping(); err != nil {
		logger.Logger.Error(err)
		return nil, err
	}
	return db, nil
}

func (dm databaseManager) getConnectionString() string {
	connectingStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dm.cfg.DatabaseHost, dm.cfg.DatabasePort, dm.cfg.DatabaseUser, dm.cfg.DatabasePassword, dm.cfg.DatabaseName)
	return connectingStr
}

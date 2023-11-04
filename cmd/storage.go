package main

import (
	"database/sql"
	"fmt"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"

	"go.uber.org/zap"
)

func initPostgres(logger *zap.Logger, dbCfg dbConfig) *sql.DB {
	database, err := sql.Open(
		"postgres",
		fmt.Sprintf("user=%s password=%s sslmode=disable", dbCfg.User, dbCfg.Password),
	)
	if err != nil {
		logger.Fatal("Internal PostgreSQL", zap.Error(err))
	}
	return database
}

func initRepository(logger *zap.Logger, appCfg appConfig) service.UrlRepository {
	switch appCfg.StorageType {
	case "postgres":
		logger.Info("PostgreSQL storage initializing")
		return repository.NewPgDatabase(initPostgres(logger, appCfg.Db))
	case "in-memory":
		logger.Info("In-memory storage initializing")
		return repository.NewInMemoryDatabase()
	}
	logger.Fatal("Unknown storage type")
	return nil
}

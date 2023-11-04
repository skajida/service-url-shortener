package main

import (
	"github.com/caarlos0/env/v9"
	"go.uber.org/zap"
)

type dbConfig struct {
	User     string `env:"PG_USER" envDefault:"postgres"`
	Password string `env:"PG_PASSWORD"`
}

type appConfig struct {
	Db          dbConfig
	Port        uint16 `env:"SERVICE_PORT,notEmpty"`
	GrpcPort    uint16 `env:"GRPC_PORT,notEmpty"`
	StorageType string `env:"STORAGE_TYPE,notEmpty"`
}

func initConfig(logger *zap.Logger) (appCfg appConfig) {
	if err := env.Parse(&appCfg); err != nil {
		logger.Fatal("Internal application config", zap.Error(err))
	}
	return
}

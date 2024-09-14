package config

import (
	"crypto/rand"
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

/*
адрес и порт запуска сервиса: переменная окружения ОС RUN_ADDRESS или флаг -a;
адрес подключения к базе данных: переменная окружения ОС DATABASE_URI или флаг -d;
адрес системы расчёта начислений: переменная окружения ОС ACCRUAL_SYSTEM_ADDRESS или флаг -r.
*/

type ServerConfig struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseDSN          string `env:"DATABASE_URI"`
	Secret               []byte
}

func NewConfig() (*ServerConfig, error) {
	var params ServerConfig
	err := env.Parse(&params)
	if err != nil {
		return nil, err
	}

	var commandLineParams ServerConfig

	flag.StringVar(&commandLineParams.RunAddress, "a", "localhost:8080", "Base address to listen on")
	flag.StringVar(&commandLineParams.DatabaseDSN, "d", "postgres://postgres@localhost:5432/postgres?sslmode=disable", "Database DSN")
	flag.Parse()

	if params.RunAddress == "" {
		params.RunAddress = commandLineParams.RunAddress
	}
	if params.DatabaseDSN == "" {
		params.DatabaseDSN = commandLineParams.DatabaseDSN
	}

	secret := make([]byte, 10)
	_, err = rand.Read(secret)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}
	params.Secret = secret

	return &params, nil
}

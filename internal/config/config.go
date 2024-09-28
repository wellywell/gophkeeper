// Package config отвечает за конфигурацию клиента и сервера
package config

import (
	"crypto/rand"
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

/*
адрес и порт запуска сервера: переменная окружения ОС RUN_ADDRESS или флаг -a;
адрес подключения к базе данных: переменная окружения ОС DATABASE_URI или флаг -d;
путь к серверному сертификату и ключу SSL_CERT_PATH и SSL_KEY_PATH или флаги -с -k

для запуска клиента:
адрес сервера env SERVER_ADDRESS или флаг -s
путь к файлу ключа сертификата CA_KEY или флаг -ssl
*/

// ServerConfig структура с параметрами для сервера
type ServerConfig struct {
	RunAddress  string `env:"RUN_ADDRESS"`
	DatabaseDSN string `env:"DATABASE_URI"`
	Secret      []byte
	SSLCert     string `env:"SSL_CERT_PATH"`
	SSLKey      string `env:"SSL_KEY_PATH"`
}

// ClientConfig структура с параметрами для клиента
type ClientConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	SSLKey        string `env:"CA_KEY"`
}

// NewServerConfig конструктор для создания конфига сервера
func NewServerConfig() (*ServerConfig, error) {
	var params ServerConfig
	err := env.Parse(&params)
	if err != nil {
		return nil, err
	}

	var commandLineParams ServerConfig

	flag.StringVar(&commandLineParams.RunAddress, "a", "localhost:8080", "Base address to listen on")
	flag.StringVar(&commandLineParams.DatabaseDSN, "d", "postgres://postgres@localhost:5432/postgres?sslmode=disable", "Database DSN")
	flag.StringVar(&commandLineParams.SSLCert, "c", "../../.ssl/server.crt", "Path to certificate")
	flag.StringVar(&commandLineParams.SSLKey, "k", "../../.ssl/server.key", "Path to certificate key")
	flag.Parse()

	if params.RunAddress == "" {
		params.RunAddress = commandLineParams.RunAddress
	}
	if params.DatabaseDSN == "" {
		params.DatabaseDSN = commandLineParams.DatabaseDSN
	}
	if params.SSLCert == "" {
		params.SSLCert = commandLineParams.SSLCert
	}
	if params.SSLKey == "" {
		params.SSLKey = commandLineParams.SSLKey
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

// NewClientConfig структура для создания клиента конфига
func NewClientConfig() (*ClientConfig, error) {
	var params ClientConfig
	err := env.Parse(&params)
	if err != nil {
		return nil, err
	}

	var commandLineParams ClientConfig

	flag.StringVar(&commandLineParams.ServerAddress, "s", "localhost:8080", "Server address")
	flag.StringVar(&commandLineParams.SSLKey, "ssl", "../../.ssl/ca.key", "Path to certificate key")
	flag.Parse()

	if params.ServerAddress == "" {
		params.ServerAddress = commandLineParams.ServerAddress
	}
	if params.SSLKey == "" {
		params.SSLKey = commandLineParams.SSLKey
	}
	return &params, nil
}

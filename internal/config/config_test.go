// Package config отвечает за конфигурацию клиента и сервера
package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServerConfig(t *testing.T) {
	got, err := NewServerConfig()
	assert.NoError(t, err)

	assert.Equal(t, "localhost:8080", got.RunAddress)
	assert.Equal(t, "postgres://postgres@localhost:5432/postgres?sslmode=disable", got.DatabaseDSN)
	assert.Equal(t, "../../.ssl/server.key", got.SSLKey)
	assert.Equal(t, "../../.ssl/server.crt", got.SSLCert)

}

func TestNewClientConfig(t *testing.T) {

	got, err := NewClientConfig()
	assert.NoError(t, err)

	assert.Equal(t, "https://localhost:8080", got.ServerAddress)
	assert.Equal(t, "../../.ssl/ca.key", got.SSLKey)
}

// Package config отвечает за конфигурацию клиента и сервера
package config

import (
	"reflect"
	"testing"
)

func TestNewServerConfig(t *testing.T) {
	tests := []struct {
		name    string
		want    *ServerConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewServerConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServerConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServerConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClientConfig(t *testing.T) {
	tests := []struct {
		name    string
		want    *ClientConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClientConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClientConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClientConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

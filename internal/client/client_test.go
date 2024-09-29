// Package client содержит методы для создания http-запросов от клиента на сервер
package client

import (
	"reflect"
	"testing"

	"github.com/wellywell/gophkeeper/internal/config"
)

func TestNewClient(t *testing.T) {
	type args struct {
		conf *config.ClientConfig
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

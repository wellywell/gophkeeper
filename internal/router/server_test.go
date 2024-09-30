// Package router отвечает за инициализацию объекта роутер, матчинг набора маршрутов с нужными хендлерами-обработчикамиб
// а также подключает требуемые middleware
package router

import (
	"context"
	"testing"
	"time"

	"github.com/wellywell/gophkeeper/internal/config"
	"github.com/wellywell/gophkeeper/internal/handlers"
	"gotest.tools/assert"
)

func TestNewServer(t *testing.T) {
	type args struct {
		conf        config.ServerConfig
		h           handlers.HandlerSet
		middlewares []Middleware
	}
	conf := config.ServerConfig{}
	tests := []struct {
		name string
		args args
	}{
		{"create server", args{conf, handlers.HandlerSet{}, []Middleware{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewServer(tt.args.conf, tt.args.h, tt.args.middlewares...)
			assert.Equal(t, got.server.Addr, conf.RunAddress)
		})
	}
}

func TestServer_ListenAndServe_Shutdown(t *testing.T) {

	s := NewServer(config.ServerConfig{}, handlers.HandlerSet{})

	go func() {
		err := s.ListenAndServe()
		assert.ErrorContains(t, err, "http: Server closed")
	}()
	time.Sleep(1000)

	err := s.Shutdown(context.Background())
	assert.NilError(t, err)
}

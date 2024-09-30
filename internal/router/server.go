// Package router отвечает за инициализацию объекта роутер, матчинг набора маршрутов с нужными хендлерами-обработчикамиб
// а также подключает требуемые middleware
package router

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/wellywell/gophkeeper/internal/auth"
	"github.com/wellywell/gophkeeper/internal/config"
	"github.com/wellywell/gophkeeper/internal/handlers"
)

// Middleware - интерфейс, которому должны соответствовать используемые Middleware
type Middleware interface {
	Handle(h http.Handler) http.Handler
}

// Router - объект роутера
type Server struct {
	server http.Server
	config config.ServerConfig
}

// NewRouter инициализирует Router, прописывает пути, на которых сервер будет слушать
func NewServer(conf config.ServerConfig, h handlers.HandlerSet, middlewares ...Middleware) *Server {

	r := chi.NewRouter()

	for _, m := range middlewares {
		r.Use(m.Handle)
	}
	r.Use(middleware.Logger)

	r.Post("/api/user/register", h.HandleRegisterUser)
	r.Post("/api/user/login", h.HandleLogin)

	authMiddleware := &auth.AuthenticateMiddleware{Secret: conf.Secret}

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Handle)
		r.Post("/api/item/login_password", h.HandleStoreLoginAndPassword)
		r.Put("/api/item/login_password", h.HandleUpdateLoginAndPassword)
		r.Post("/api/item/credit_card", h.HandleStoreCreditCard)
		r.Put("/api/item/credit_card", h.HandleUpdateCreditCard)
		r.Post("/api/item/binary", h.HandleStoreBinaryItem)
		r.Put("/api/item/binary", h.HandleUpdateBinaryItem)
		r.Get("/api/item/{key}", h.HandleGetItem)
		r.Delete("/api/item/{key}", h.HandleDeleteItem)
		r.Post("/api/item/text", h.HandleStoreText)
		r.Put("/api/item/text", h.HandleUpdateText)
		r.Get("/api/item/binary/{key}/download", h.HandleDownloadBinaryItem)
		r.Get("/api/item/list", h.HandleItemList)
	})

	return &Server{server: http.Server{Addr: conf.RunAddress, Handler: r}, config: conf}
}

// ListenAndServe - метод для запуска сервера
func (s *Server) ListenAndServe() error {
	err := s.server.ListenAndServeTLS(s.config.SSLCert, s.config.SSLKey)
	return err
}

// Shutdown gracefull shutddown
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

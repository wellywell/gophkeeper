package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/wellywell/gophkeeper/internal/auth"
	"github.com/wellywell/gophkeeper/internal/db"
)

//go:generate mockery --name Database
type Database interface {
	GetUserHashedPassword(context.Context, string) (string, error)
	CreateUser(context.Context, string, string) error
	GetUserID(context.Context, string) (int, error)
}

type HandlerSet struct {
	secret   []byte
	database Database
}

var (
	ErrCouldNotParseBody = errors.New("could not parse body")
	ErrAuthDataEmpty     = errors.New("login or password cannot be empty")
)

func NewHandlerSet(secret []byte, database Database) *HandlerSet {
	return &HandlerSet{
		secret:   secret,
		database: database,
	}
}

func (h *HandlerSet) parseAuthData(body []byte) (username string, password string, err error) {

	var data struct {
		Username string `json:"login"`
		Password string `json:"password"`
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", "", ErrCouldNotParseBody
	}

	if data.Username == "" || data.Password == "" {
		return "", "", ErrAuthDataEmpty
	}

	return data.Username, data.Password, nil

}

func (h *HandlerSet) handleAuthErrors(err error, w http.ResponseWriter) {

	if errors.Is(err, ErrCouldNotParseBody) {
		http.Error(w, "Could not parse body",
			http.StatusBadRequest)
	} else if errors.Is(err, ErrAuthDataEmpty) {
		http.Error(w, "Login and password cannot be empty",
			http.StatusBadRequest)
	} else {
		http.Error(w, "Unknown error", http.StatusInternalServerError)
	}
}

func (h *HandlerSet) HandleLogin(w http.ResponseWriter, req *http.Request) {

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Something went wrong",
			http.StatusInternalServerError)
		return
	}

	username, password, err := h.parseAuthData(body)

	if err != nil {
		h.handleAuthErrors(err, w)
		return
	}

	passwordInDB, err := h.database.GetUserHashedPassword(req.Context(), username)
	if err != nil {
		var userNotFound *db.UserNotFoundError
		if errors.As(err, &userNotFound) {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !auth.CheckPasswordHash(password, passwordInDB) {
		http.Error(w, "Wrong password", http.StatusUnauthorized)
		return
	}

	err = auth.SetToken(username, w, h.secret)
	if err != nil {
		http.Error(w, "Something went wrong",
			http.StatusInternalServerError)
	}

	w.Header().Set("content-type", "text/plain")

	_, err = w.Write([]byte("success"))
	if err != nil {
		http.Error(w, "Something went wrong",
			http.StatusInternalServerError)
	}
}

func (h *HandlerSet) HandleRegisterUser(w http.ResponseWriter, req *http.Request) {

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Something went wrong",
			http.StatusInternalServerError)
		return
	}

	username, password, err := h.parseAuthData(body)

	if err != nil {
		h.handleAuthErrors(err, w)
		return
	}

	hashed, err := auth.HashPassword(password)
	if err != nil {
		http.Error(w, "Something went wrong",
			http.StatusInternalServerError)
		return
	}

	err = h.database.CreateUser(req.Context(), username, hashed)
	if err != nil {
		var userExists *db.UserExistsError
		if errors.As(err, &userExists) {
			http.Error(w, "User exists", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = auth.SetToken(username, w, h.secret)
	if err != nil {
		http.Error(w, "Something went wrong",
			http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "text/plain")

	_, err = w.Write([]byte("success"))
	if err != nil {
		http.Error(w, "Something went wrong",
			http.StatusInternalServerError)
	}
}

func (h *HandlerSet) handleAuthorizeUser(w http.ResponseWriter, req *http.Request) (int, error) {
	username, ok := auth.GetAuthenticatedUser(req)
	if !ok {
		http.Error(w, "Something went wrong",
			http.StatusInternalServerError)
		return 0, fmt.Errorf("authentication error")
	}

	userID, err := h.database.GetUserID(req.Context(), username)
	if err != nil {
		http.Error(w, "User not found",
			http.StatusUnauthorized)
		return 0, err
	}
	return userID, nil

}

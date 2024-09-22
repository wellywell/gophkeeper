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
	"github.com/wellywell/gophkeeper/internal/types"
)

//go:generate mockery --name Database
type Database interface {
	GetUserHashedPassword(context.Context, string) (string, error)
	CreateUser(context.Context, string, string) error
	GetUserID(context.Context, string) (int, error)
	InsertLogoPass(context.Context, int, types.LoginPasswordItem) error
	InsertCreditCard(context.Context, int, types.CreditCardItem) error
	GetItem(context.Context, int, string) (*types.Item, error)
	GetLogoPass(context.Context, int) (*types.LoginPassword, error)
	GetCreditCard(context.Context, int) (*types.CreditCardData, error)
	DeleteItem(context.Context, int, string) error
	UpdateLogoPass(context.Context, int, types.LoginPasswordItem) error
	UpdateCreditCard(context.Context, int, types.CreditCardItem) error
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

func (h *HandlerSet) prepareLoginAndPasswordItem(w http.ResponseWriter, req *http.Request) (*types.LoginPasswordItem, error) {

	var logopass *types.LoginPasswordItem

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Something went wrong",
			http.StatusUnauthorized)
		return nil, err
	}
	err = json.Unmarshal(body, &logopass)

	if err != nil {
		http.Error(w, "Could not unmarshal body",
			http.StatusBadRequest)
		return nil, err
	}
	return logopass, nil
}

func (h *HandlerSet) prepareCreditCardItem(w http.ResponseWriter, req *http.Request) (*types.CreditCardItem, error) {

	var card *types.CreditCardItem

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Something went wrong",
			http.StatusUnauthorized)
		return nil, err
	}
	err = json.Unmarshal(body, &card)

	if err != nil {
		http.Error(w, "Could not unmarshal body",
			http.StatusBadRequest)
		return nil, err
	}
	return card, nil
}

func (h *HandlerSet) HandleUpdateLoginAndPassword(w http.ResponseWriter, req *http.Request) {
	userID, err := h.handleAuthorizeUser(w, req)
	if err != nil {
		http.Error(w, "Error authenticating",
			http.StatusUnauthorized)
	}
	logopass, err := h.prepareLoginAndPasswordItem(w, req)
	if err != nil {
		return
	}
	err = h.database.UpdateLogoPass(req.Context(), userID, *logopass)

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Could not update",
			http.StatusInternalServerError)
		return
	}
}

func (h *HandlerSet) HandleUpdateCreditCard(w http.ResponseWriter, req *http.Request) {
	userID, err := h.handleAuthorizeUser(w, req)
	if err != nil {
		http.Error(w, "Error authenticating",
			http.StatusUnauthorized)
	}
	card, err := h.prepareCreditCardItem(w, req)
	if err != nil {
		return
	}
	err = h.database.UpdateCreditCard(req.Context(), userID, *card)

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Could not update",
			http.StatusInternalServerError)
		return
	}
}

func (h *HandlerSet) HandleStoreLoginAndPassword(w http.ResponseWriter, req *http.Request) {

	userID, err := h.handleAuthorizeUser(w, req)

	if err != nil {
		http.Error(w, "Something went wrong",
			http.StatusUnauthorized)
	}

	logopass, err := h.prepareLoginAndPasswordItem(w, req)
	if err != nil {
		return
	}
	err = h.database.InsertLogoPass(req.Context(), userID, *logopass)

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Something went wrong",
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *HandlerSet) HandleStoreCreditCard(w http.ResponseWriter, req *http.Request) {

	userID, err := h.handleAuthorizeUser(w, req)

	if err != nil {
		http.Error(w, "Something went wrong",
			http.StatusUnauthorized)
	}

	card, err := h.prepareCreditCardItem(w, req)
	if err != nil {
		return
	}
	err = h.database.InsertCreditCard(req.Context(), userID, *card)

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Something went wrong",
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *HandlerSet) HandleGetItem(w http.ResponseWriter, req *http.Request) {

	userID, err := h.handleAuthorizeUser(w, req)

	if err != nil {
		http.Error(w, "Something went wrong",
			http.StatusUnauthorized)
	}

	idString := req.PathValue("key")

	if idString == "" {
		http.Error(w, "Key not passed", http.StatusBadRequest)
		return
	}
	item, err := h.database.GetItem(req.Context(), userID, idString)

	if err != nil {
		var keyNotFound *db.KeyNotFoundError
		if errors.As(err, &keyNotFound) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Something went wrong",
			http.StatusInternalServerError)
		return
	}

	var data []byte

	switch item.Type {
	case types.TypeLogoPass:
		logopass, err := h.database.GetLogoPass(req.Context(), item.Id)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		result := types.LoginPasswordItem{Item: *item, Data: logopass}
		data, err = json.Marshal(result)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	case types.TypeCreditCard:
		card, err := h.database.GetCreditCard(req.Context(), item.Id)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		result := types.CreditCardItem{Item: *item, Data: card}
		data, err = json.Marshal(result)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("content-type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "Something went wrong",
			http.StatusInternalServerError)
	}
}

func (h *HandlerSet) HandleDeleteItem(w http.ResponseWriter, req *http.Request) {

	userID, err := h.handleAuthorizeUser(w, req)

	if err != nil {
		http.Error(w, "Something went wrong",
			http.StatusUnauthorized)
	}

	idString := req.PathValue("key")

	if idString == "" {
		http.Error(w, "Key not passed", http.StatusBadRequest)
		return
	}
	err = h.database.DeleteItem(req.Context(), userID, idString)

	if err != nil {
		http.Error(w, "Something went wrong",
			http.StatusInternalServerError)
		return
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

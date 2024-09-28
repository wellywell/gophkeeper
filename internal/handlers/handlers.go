// Package handlers определяет обработчики http-запросов для сервера
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/wellywell/gophkeeper/internal/auth"
	"github.com/wellywell/gophkeeper/internal/db"
	"github.com/wellywell/gophkeeper/internal/types"
)

// Database интерфейс определяет методы для работы с БД
//
//go:generate mockery --name Database
type Database interface {
	GetUserHashedPassword(context.Context, string) (string, error)
	CreateUser(context.Context, string, string) error
	GetUserID(context.Context, string) (int, error)
	InsertLogoPass(context.Context, int, types.LoginPasswordItem) error
	InsertCreditCard(context.Context, int, types.CreditCardItem) error
	InsertText(context.Context, int, types.TextItem) error
	InsertBinaryData(context.Context, int, types.BinaryItem) error
	UpdateBinaryData(context.Context, int, types.BinaryItem) error
	GetItem(context.Context, int, string) (*types.Item, error)
	GetItems(context.Context, int, int, int) ([]types.Item, error)
	GetBinaryData(context.Context, int, string) ([]byte, error)
	GetLogoPass(context.Context, int) (*types.LoginPassword, error)
	GetCreditCard(context.Context, int) (*types.CreditCardData, error)
	GetText(context.Context, int) (*types.TextData, error)
	DeleteItem(context.Context, int, string) error
	UpdateLogoPass(context.Context, int, types.LoginPasswordItem) error
	UpdateCreditCard(context.Context, int, types.CreditCardItem) error
	UpdateText(context.Context, int, types.TextItem) error
}

// HandlerSet структура для работы с хендлерами
type HandlerSet struct {
	secret   []byte
	database Database
}

var (
	ErrCouldNotParseBody = errors.New("could not parse body")
	ErrAuthDataEmpty     = errors.New("login or password cannot be empty")
)

// NewHandlerSet инициализирует набор хендлеров
func NewHandlerSet(secret []byte, database Database) *HandlerSet {
	return &HandlerSet{
		secret:   secret,
		database: database,
	}
}

// HandleLogin авторизация пользователя
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

// HandleRegisterUser регистрация пользователя
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

// HandleStoreLoginAndPassword хендлер, обрабатывающий запрос на сохранение на сервере данных типа "логин и пароль"
func (h *HandlerSet) HandleStoreLoginAndPassword(w http.ResponseWriter, req *http.Request) {

	userID, err := h.handleAuthorizeUser(w, req)

	if err != nil {
		return
	}

	logopass, err := h.prepareLoginAndPasswordItem(w, req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = h.database.InsertLogoPass(req.Context(), userID, *logopass)
	if err != nil {
		var keyExistsError *db.KeyExistsError
		if errors.As(err, &keyExistsError) {
			http.Error(w, "Key exists", http.StatusConflict)
			return
		}
		fmt.Println(err.Error())
		http.Error(w, "Something went wrong",
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// HandleUpdateLoginAndPassword хендлер, обрабатывающий запрос на изменени данных, хранимых не сервере типа "логин-пароль"
func (h *HandlerSet) HandleUpdateLoginAndPassword(w http.ResponseWriter, req *http.Request) {
	userID, err := h.handleAuthorizeUser(w, req)
	if err != nil {
		http.Error(w, "Error authenticating",
			http.StatusUnauthorized)
		return
	}
	logopass, err := h.prepareLoginAndPasswordItem(w, req)
	if err != nil {
		return
	}
	err = h.database.UpdateLogoPass(req.Context(), userID, *logopass)

	if err != nil {
		var keyNotFound *db.KeyNotFoundError
		if errors.As(err, &keyNotFound) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		fmt.Println(err.Error())
		http.Error(w, "Could not update",
			http.StatusInternalServerError)
		return
	}
}

// HandleUpdateCreditCard хендлер, обрабатывающий запрос на изменени данных, хранимых не сервере типа "кредитная карта"
func (h *HandlerSet) HandleUpdateCreditCard(w http.ResponseWriter, req *http.Request) {
	userID, err := h.handleAuthorizeUser(w, req)
	if err != nil {
		return
	}
	card, err := h.prepareCreditCardItem(w, req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = h.database.UpdateCreditCard(req.Context(), userID, *card)

	if err != nil {
		var keyNotFound *db.KeyNotFoundError
		if errors.As(err, &keyNotFound) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		fmt.Println(err.Error())
		http.Error(w, "Could not update",
			http.StatusInternalServerError)
		return
	}
}

// HandleUpdateText обрабатывает запрос на обновление текстовых данных, хранимых на сервере
func (h *HandlerSet) HandleUpdateText(w http.ResponseWriter, req *http.Request) {

	userID, err := h.handleAuthorizeUser(w, req)

	if err != nil {
		http.Error(w, "Something went wrong",
			http.StatusUnauthorized)
		return
	}

	text, err := h.prepareTextItem(w, req)
	if err != nil {
		return
	}
	err = h.database.UpdateText(req.Context(), userID, *text)

	if err != nil {
		var keyNotFound *db.KeyNotFoundError
		if errors.As(err, &keyNotFound) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		fmt.Println(err.Error())
		http.Error(w, "Something went wrong",
			http.StatusInternalServerError)
		return
	}
}

// HandleStoreText обрабатывает запрос на создание текстовых данных для хранения на сервере
func (h *HandlerSet) HandleStoreText(w http.ResponseWriter, req *http.Request) {

	userID, err := h.handleAuthorizeUser(w, req)

	if err != nil {
		http.Error(w, "Something went wrong",
			http.StatusUnauthorized)
		return
	}

	text, err := h.prepareTextItem(w, req)
	if err != nil {
		return
	}
	err = h.database.InsertText(req.Context(), userID, *text)

	if err != nil {
		var keyExistsError *db.KeyExistsError
		if errors.As(err, &keyExistsError) {
			http.Error(w, "Key exists", http.StatusConflict)
			return
		}
		fmt.Println(err.Error())
		http.Error(w, "Something went wrong",
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// HandleStoreCreditCard обрабатывает запрос на создание записи с данными кредитной карты на сервере
func (h *HandlerSet) HandleStoreCreditCard(w http.ResponseWriter, req *http.Request) {

	userID, err := h.handleAuthorizeUser(w, req)

	if err != nil {
		http.Error(w, "Something went wrong",
			http.StatusUnauthorized)
		return
	}

	card, err := h.prepareCreditCardItem(w, req)
	if err != nil {
		return
	}
	err = h.database.InsertCreditCard(req.Context(), userID, *card)

	if err != nil {
		var keyExistsError *db.KeyExistsError
		if errors.As(err, &keyExistsError) {
			http.Error(w, "Key exists", http.StatusConflict)
			return
		}
		fmt.Println(err.Error())
		http.Error(w, "Something went wrong",
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// HandleStoreBinaryItem обрабатывает запрос на сохранение бинарных данных на сервере
func (h *HandlerSet) HandleStoreBinaryItem(w http.ResponseWriter, req *http.Request) {

	userID, err := h.handleAuthorizeUser(w, req)

	if err != nil {
		http.Error(w, "Something went wrong",
			http.StatusUnauthorized)
		return
	}

	item, err := h.prepareBinaryItem(w, req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = h.database.InsertBinaryData(req.Context(), userID, *item)
	if err != nil {
		var keyExistsError *db.KeyExistsError
		if errors.As(err, &keyExistsError) {
			http.Error(w, "Key exists", http.StatusConflict)
			return
		}
		fmt.Println(err.Error())
		http.Error(w, "Something went wrong",
			http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// HandleUpdateBinaryItem обрабатывает запрос на обновление бинарных данных на сервере
func (h *HandlerSet) HandleUpdateBinaryItem(w http.ResponseWriter, req *http.Request) {

	userID, err := h.handleAuthorizeUser(w, req)

	if err != nil {
		http.Error(w, "Something went wrong",
			http.StatusUnauthorized)
		return
	}

	item, err := h.prepareBinaryItem(w, req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = h.database.UpdateBinaryData(req.Context(), userID, *item)

	if err != nil {
		var keyNotFound *db.KeyNotFoundError
		if errors.As(err, &keyNotFound) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		fmt.Println(err.Error())
		http.Error(w, "Something went wrong",
			http.StatusInternalServerError)
		return
	}
}

// HandleDownloadBinaryItem обрабатывает запрос на скачивание бинарных данных
func (h *HandlerSet) HandleDownloadBinaryItem(w http.ResponseWriter, req *http.Request) {
	userID, err := h.handleAuthorizeUser(w, req)

	if err != nil {
		return
	}

	idString := req.PathValue("key")

	if idString == "" {
		http.Error(w, "Key not passed", http.StatusBadRequest)
		return
	}
	data, err := h.database.GetBinaryData(req.Context(), userID, idString)

	if err != nil {
		var keyNotFound *db.KeyNotFoundError
		if errors.As(err, &keyNotFound) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Something went wrong",
			http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	w.Header().Set("content-type", "application/octet-stream")
	_, err = w.Write(data)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Something went wrong",
			http.StatusInternalServerError)
	}
}

// HandleItemList обрабатывает запрос на получение списка метаданных о записях, хранимых на сервере
func (h *HandlerSet) HandleItemList(w http.ResponseWriter, req *http.Request) {

	userID, err := h.handleAuthorizeUser(w, req)

	if err != nil {
		return
	}
	s := req.URL.Query().Get("page")

	var page int
	var limit int

	if s == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(s)
		if err != nil {
			http.Error(w, "Error parsing page", http.StatusBadRequest)
			return
		}
	}
	s = req.URL.Query().Get("limit")
	if s == "" {
		limit = 10
	} else {
		limit, err = strconv.Atoi(s)
		if err != nil {
			http.Error(w, "Error parsing limit", http.StatusBadRequest)
			return
		}
	}

	offset := (page - 1) * limit
	items, err := h.database.GetItems(req.Context(), userID, limit, offset)

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Something went wrong",
			http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(items)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "Something went wrong",
			http.StatusInternalServerError)
	}
}

// HandleGetItem возвращает запись, хранимую на сервере произвольного типа (из числа поддерживаемых)
func (h *HandlerSet) HandleGetItem(w http.ResponseWriter, req *http.Request) {

	userID, err := h.handleAuthorizeUser(w, req)

	if err != nil {
		http.Error(w, "Something went wrong",
			http.StatusUnauthorized)
		return
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
	case types.TypeText:
		text, err := h.database.GetText(req.Context(), item.Id)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		result := types.TextItem{Item: *item, Data: *text}
		data, err = json.Marshal(result)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	case types.TypeBinary:
		// only metadata in this handler
		item, err := h.database.GetItem(req.Context(), userID, item.Key)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		result := types.BinaryItem{Item: *item, Data: []byte{}}
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

// HandleDeleteItem удаляет данные с сервера
func (h *HandlerSet) HandleDeleteItem(w http.ResponseWriter, req *http.Request) {

	userID, err := h.handleAuthorizeUser(w, req)

	if err != nil {
		http.Error(w, "Something went wrong",
			http.StatusUnauthorized)
		return
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

func (h *HandlerSet) prepareTextItem(w http.ResponseWriter, req *http.Request) (*types.TextItem, error) {

	var text *types.TextItem

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Something went wrong",
			http.StatusUnauthorized)
		return nil, err
	}
	err = json.Unmarshal(body, &text)

	if err != nil {
		http.Error(w, "Could not unmarshal body",
			http.StatusBadRequest)
		return nil, err
	}
	return text, nil
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

func (h *HandlerSet) prepareBinaryItem(w http.ResponseWriter, r *http.Request) (*types.BinaryItem, error) {

	contentType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || !strings.HasPrefix(contentType, "multipart/") {
		http.Error(w, "expecting a multipart message", http.StatusBadRequest)
		return nil, err
	}
	multipartReader := multipart.NewReader(r.Body, params["boundary"])
	defer r.Body.Close()

	item := types.BinaryItem{}

	for {
		part, err := multipartReader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "unexpected error when retrieving a part of the message", http.StatusBadRequest)
			return nil, err
		}
		defer part.Close()

		fileBytes, err := io.ReadAll(part)
		if err != nil {
			http.Error(w, "failed to read content of the part", http.StatusBadRequest)
			return nil, err
		}
		switch part.Header.Get("Content-Type") {
		case "application/json":
			err = json.Unmarshal(fileBytes, &item.Item)
			if err != nil {
				http.Error(w, "failed to read metadata", http.StatusBadRequest)
				return nil, err
			}
		case "application/octet-stream":
			item.Data = fileBytes
		}
	}
	return &item, nil
}

func (h *HandlerSet) handleAuthorizeUser(w http.ResponseWriter, req *http.Request) (int, error) {
	username, ok := auth.GetAuthenticatedUser(req)
	if !ok {
		http.Error(w, "Something went wrong",
			http.StatusUnauthorized)
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

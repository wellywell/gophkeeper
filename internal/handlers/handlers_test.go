package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mock "github.com/stretchr/testify/mock"
	"github.com/wellywell/gophkeeper/internal/auth"
	"github.com/wellywell/gophkeeper/internal/db"
	"github.com/wellywell/gophkeeper/internal/types"
	"gotest.tools/assert"
)

var (
	binaryData = []byte(`--56a7182d3cfb7e97d66458cd95b042fbd704c16ef53b72fa871890d92f15
Content-Type: application/json

{"key":"111","type":"binary"}
--56a7182d3cfb7e97d66458cd95b042fbd704c16ef53b72fa871890d92f15
Content-Type: application/octet-stream

test
--56a7182d3cfb7e97d66458cd95b042fbd704c16ef53b72fa871890d92f15--`)
	binaryItem     = types.BinaryItem{Item: types.Item{Key: "111", Type: types.TypeBinary}, Data: []byte("test")}
	token          = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6InVzZXIifQ.2MU1vpPXOY4ypHv1bPVpTqOG0zQmk4-ZD8Qoze4xVKg"
	textBody       = []byte(`{"item": {"type": "text", "key": "111"}, "data": "text"}`)
	logopassBody   = []byte(`{"item": {"type": "logopass", "key": "111"}, "data": {"login": "112", "password": "222"}}`)
	logopassItem   = types.LoginPasswordItem{Item: types.Item{Key: "111", Type: types.TypeLogoPass}, Data: &types.LoginPassword{Login: "112", Password: "222"}}
	creditCardBody = []byte(`{"item": {"type": "credit_card", "key": "111"}, "data": {"cvc": "1", "number":"1", "name":"1", "valid_month": "1", "valid_year": "2000"}}`)
	creditCardItem = types.CreditCardItem{Item: types.Item{Key: "111", Type: types.TypeCreditCard}, Data: &types.CreditCardData{Number: "1", CVC: "1", Name: "1", ValidMonth: "1", ValidYear: "2000"}}
	textItem       = types.TextItem{Item: types.Item{Key: "111", Type: types.TypeText}, Data: types.TextData("text")}
)

func TestHandlerSet_HandleLogin(t *testing.T) {
	type fields struct {
		secret []byte
	}

	db := &MockDatabase{}
	tests := []struct {
		name         string
		fields       fields
		login        string
		password     string
		body         []byte
		userExists   bool
		expextedCode int
	}{
		{"userOK", fields{[]byte("secret")}, "user", "pass", []byte(`{"login": "user", "password": "pass"}`), true, http.StatusOK},
		{"wrongPassword", fields{[]byte("secret")}, "user", "pass", []byte(`{"login": "user", "password": "wrong"}`), true, http.StatusUnauthorized},
		{"userNotExists", fields{[]byte("secret")}, "user", "pass", []byte(`{"login": "user", "password": "wrong"}`), false, http.StatusUnauthorized},
		{"wrongRequest", fields{[]byte("secret")}, "user", "pass", []byte(`{}`), true, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HandlerSet{
				secret:   tt.fields.secret,
				database: db,
			}

			req, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(tt.body))
			w := httptest.NewRecorder()

			hash, _ := auth.HashPassword(tt.password)

			if tt.userExists {
				db.EXPECT().GetUserHashedPassword(req.Context(), tt.login).Return(hash, nil)
				db.EXPECT().GetUserID(req.Context(), tt.login).Return(1, nil)
			} else {
				db.EXPECT().GetUserID(req.Context(), tt.login).Return(0, fmt.Errorf("smth"))
			}
			h.HandleLogin(w, req)
			assert.Equal(t, tt.expextedCode, w.Code)

			if tt.expextedCode == http.StatusOK {
				assert.Equal(t, token, w.Header().Get("X-Auth-Token"))
			}
		})
	}
}

func TestHandlerSet_HandleRegisterUser(t *testing.T) {
	type fields struct {
		secret []byte
	}

	tests := []struct {
		name         string
		fields       fields
		login        string
		password     string
		body         []byte
		userExists   bool
		expextedCode int
	}{
		{"userOK", fields{[]byte("secret")}, "user", "pass", []byte(`{"login": "user", "password": "pass"}`), false, http.StatusOK},
		{"wrongRequest", fields{[]byte("secret")}, "user", "pass", []byte(`{}`), false, http.StatusBadRequest},
		{"userExists", fields{[]byte("secret")}, "user", "pass", []byte(`{"login": "user", "password": "pass"}`), true, http.StatusConflict},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mdb := &MockDatabase{}
			h := &HandlerSet{
				secret:   tt.fields.secret,
				database: mdb,
			}
			req, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(tt.body))
			w := httptest.NewRecorder()

			if tt.userExists {
				err := fmt.Errorf("%w", &db.UserExistsError{Username: tt.login})
				mdb.EXPECT().CreateUser(req.Context(), tt.login, mock.Anything).Return(err)
			} else {
				mdb.EXPECT().CreateUser(req.Context(), tt.login, mock.Anything).Return(nil)
			}
			h.HandleRegisterUser(w, req)
			assert.Equal(t, tt.expextedCode, w.Code)

			if tt.expextedCode == http.StatusOK {
				assert.Equal(t, token, w.Header().Get("X-Auth-Token"))
			}
		})
	}
}

func TestHandlerSet_HandleStoreLoginAndPassword(t *testing.T) {

	tests := []struct {
		name               string
		isAuthorized       bool
		keyExists          bool
		userExists         bool
		body               []byte
		item               types.LoginPasswordItem
		expectedStatusCode int
	}{
		{"ok", true, false, true, logopassBody, logopassItem, http.StatusCreated},
		{"notAuthorized", false, false, false, logopassBody, logopassItem, http.StatusUnauthorized},
		{"userNotExists", true, false, false, logopassBody, logopassItem, http.StatusUnauthorized},
		{"keyExists", true, true, true, logopassBody, logopassItem, http.StatusConflict},
		{"badData", true, false, true, []byte(`"wrong"`), logopassItem, http.StatusBadRequest},
	}
	for _, tt := range tests {
		mdb := &MockDatabase{}

		t.Run(tt.name, func(t *testing.T) {

			h := &HandlerSet{
				secret:   []byte("secret"),
				database: mdb,
			}

			req, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(tt.body))
			if tt.isAuthorized {
				const contextKey auth.UserKey = "username"
				ctx := context.WithValue(req.Context(), contextKey, "user")
				req = req.WithContext(ctx)
			}
			if tt.userExists {
				mdb.EXPECT().GetUserID(req.Context(), "user").Return(1, nil)
			} else {
				mdb.EXPECT().GetUserID(req.Context(), "user").Return(0, &db.UserNotFoundError{Username: "user"})
			}

			if tt.keyExists {
				mdb.EXPECT().InsertLogoPass(req.Context(), 1, tt.item).Return(&db.KeyExistsError{Key: "111"})
			} else {
				mdb.EXPECT().InsertLogoPass(req.Context(), 1, tt.item).Return(nil)
			}
			w := httptest.NewRecorder()
			h.HandleStoreLoginAndPassword(w, req)
			assert.Equal(t, w.Code, tt.expectedStatusCode)
		})
	}
}

func TestHandlerSet_HandleUpdateLoginAndPassword(t *testing.T) {
	tests := []struct {
		name               string
		isAuthorized       bool
		keyExists          bool
		userExists         bool
		body               []byte
		item               types.LoginPasswordItem
		expectedStatusCode int
	}{
		{"ok", true, true, true, logopassBody, logopassItem, http.StatusOK},
		{"notAuthorized", false, true, false, logopassBody, logopassItem, http.StatusUnauthorized},
		{"userNotExists", true, true, false, logopassBody, logopassItem, http.StatusUnauthorized},
		{"keyNotExists", true, false, true, logopassBody, logopassItem, http.StatusNotFound},
		{"badData", true, true, true, []byte(`"wrong"`), logopassItem, http.StatusBadRequest},
	}
	for _, tt := range tests {
		mdb := &MockDatabase{}

		t.Run(tt.name, func(t *testing.T) {

			h := &HandlerSet{
				secret:   []byte("secret"),
				database: mdb,
			}

			req, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(tt.body))
			if tt.isAuthorized {
				const contextKey auth.UserKey = "username"
				ctx := context.WithValue(req.Context(), contextKey, "user")
				req = req.WithContext(ctx)
			}
			if tt.userExists {
				mdb.EXPECT().GetUserID(req.Context(), "user").Return(1, nil)
			} else {
				mdb.EXPECT().GetUserID(req.Context(), "user").Return(0, &db.UserNotFoundError{Username: "user"})
			}

			if tt.keyExists {
				mdb.EXPECT().UpdateLogoPass(req.Context(), 1, tt.item).Return(nil)
			} else {
				mdb.EXPECT().UpdateLogoPass(req.Context(), 1, tt.item).Return(&db.KeyNotFoundError{Key: "111"})
			}
			w := httptest.NewRecorder()
			h.HandleUpdateLoginAndPassword(w, req)
			assert.Equal(t, w.Code, tt.expectedStatusCode)
		})
	}
}

func TestHandlerSet_HandleUpdateCreditCard(t *testing.T) {
	tests := []struct {
		name               string
		isAuthorized       bool
		keyExists          bool
		userExists         bool
		body               []byte
		item               types.CreditCardItem
		expectedStatusCode int
	}{
		{"ok", true, true, true, creditCardBody, creditCardItem, http.StatusOK},
		{"notAuthorized", false, true, false, creditCardBody, creditCardItem, http.StatusUnauthorized},
		{"userNotExists", true, true, false, creditCardBody, creditCardItem, http.StatusUnauthorized},
		{"keyNotExists", true, false, true, creditCardBody, creditCardItem, http.StatusNotFound},
		{"badData", true, true, true, []byte(`"wrong"`), creditCardItem, http.StatusBadRequest},
	}
	for _, tt := range tests {
		mdb := &MockDatabase{}

		t.Run(tt.name, func(t *testing.T) {

			h := &HandlerSet{
				secret:   []byte("secret"),
				database: mdb,
			}

			req, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(tt.body))
			if tt.isAuthorized {
				const contextKey auth.UserKey = "username"
				ctx := context.WithValue(req.Context(), contextKey, "user")
				req = req.WithContext(ctx)
			}
			if tt.userExists {
				mdb.EXPECT().GetUserID(req.Context(), "user").Return(1, nil)
			} else {
				mdb.EXPECT().GetUserID(req.Context(), "user").Return(0, &db.UserNotFoundError{Username: "user"})
			}

			if tt.keyExists {
				mdb.EXPECT().UpdateCreditCard(req.Context(), 1, tt.item).Return(nil)
			} else {
				mdb.EXPECT().UpdateCreditCard(req.Context(), 1, tt.item).Return(&db.KeyNotFoundError{Key: "111"})
			}
			w := httptest.NewRecorder()
			h.HandleUpdateCreditCard(w, req)
			assert.Equal(t, w.Code, tt.expectedStatusCode)
		})
	}
}

func TestHandlerSet_HandleUpdateText(t *testing.T) {
	tests := []struct {
		name               string
		isAuthorized       bool
		keyExists          bool
		userExists         bool
		body               []byte
		item               types.TextItem
		expectedStatusCode int
	}{
		{"ok", true, true, true, textBody, textItem, http.StatusOK},
		{"notAuthorized", false, true, false, textBody, textItem, http.StatusUnauthorized},
		{"userNotExists", true, true, false, textBody, textItem, http.StatusUnauthorized},
		{"keyNotExists", true, false, true, textBody, textItem, http.StatusNotFound},
		{"badData", true, true, true, []byte(`"wrong"`), textItem, http.StatusBadRequest},
	}
	for _, tt := range tests {
		mdb := &MockDatabase{}

		t.Run(tt.name, func(t *testing.T) {

			h := &HandlerSet{
				secret:   []byte("secret"),
				database: mdb,
			}

			req, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(tt.body))
			if tt.isAuthorized {
				const contextKey auth.UserKey = "username"
				ctx := context.WithValue(req.Context(), contextKey, "user")
				req = req.WithContext(ctx)
			}
			if tt.userExists {
				mdb.EXPECT().GetUserID(req.Context(), "user").Return(1, nil)
			} else {
				mdb.EXPECT().GetUserID(req.Context(), "user").Return(0, &db.UserNotFoundError{Username: "user"})
			}

			if tt.keyExists {
				mdb.EXPECT().UpdateText(req.Context(), 1, tt.item).Return(nil)
			} else {
				mdb.EXPECT().UpdateText(req.Context(), 1, tt.item).Return(&db.KeyNotFoundError{Key: "111"})
			}
			w := httptest.NewRecorder()
			h.HandleUpdateText(w, req)
			assert.Equal(t, w.Code, tt.expectedStatusCode)
		})
	}
}

func TestHandlerSet_HandleStoreText(t *testing.T) {
	tests := []struct {
		name               string
		isAuthorized       bool
		keyExists          bool
		userExists         bool
		body               []byte
		item               types.TextItem
		expectedStatusCode int
	}{
		{"ok", true, false, true, textBody, textItem, http.StatusCreated},
		{"notAuthorized", false, false, false, textBody, textItem, http.StatusUnauthorized},
		{"userNotExists", true, false, false, textBody, textItem, http.StatusUnauthorized},
		{"keyExists", true, true, true, textBody, textItem, http.StatusConflict},
		{"badData", true, false, true, []byte(`"wrong"`), textItem, http.StatusBadRequest},
	}
	for _, tt := range tests {
		mdb := &MockDatabase{}

		t.Run(tt.name, func(t *testing.T) {

			h := &HandlerSet{
				secret:   []byte("secret"),
				database: mdb,
			}

			req, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(tt.body))
			if tt.isAuthorized {
				const contextKey auth.UserKey = "username"
				ctx := context.WithValue(req.Context(), contextKey, "user")
				req = req.WithContext(ctx)
			}
			if tt.userExists {
				mdb.EXPECT().GetUserID(req.Context(), "user").Return(1, nil)
			} else {
				mdb.EXPECT().GetUserID(req.Context(), "user").Return(0, &db.UserNotFoundError{Username: "user"})
			}

			if tt.keyExists {
				mdb.EXPECT().InsertText(req.Context(), 1, tt.item).Return(&db.KeyExistsError{Key: "111"})
			} else {
				mdb.EXPECT().InsertText(req.Context(), 1, tt.item).Return(nil)
			}
			w := httptest.NewRecorder()
			h.HandleStoreText(w, req)
			assert.Equal(t, w.Code, tt.expectedStatusCode)
		})
	}
}

func TestHandlerSet_HandleStoreCreditCard(t *testing.T) {
	tests := []struct {
		name               string
		isAuthorized       bool
		keyExists          bool
		userExists         bool
		body               []byte
		item               types.CreditCardItem
		expectedStatusCode int
	}{
		{"ok", true, false, true, creditCardBody, creditCardItem, http.StatusCreated},
		{"notAuthorized", false, false, false, creditCardBody, creditCardItem, http.StatusUnauthorized},
		{"userNotExists", true, false, false, creditCardBody, creditCardItem, http.StatusUnauthorized},
		{"keyExists", true, true, true, creditCardBody, creditCardItem, http.StatusConflict},
		{"badData", true, false, true, []byte(`"wrong"`), creditCardItem, http.StatusBadRequest},
	}
	for _, tt := range tests {
		mdb := &MockDatabase{}

		t.Run(tt.name, func(t *testing.T) {

			h := &HandlerSet{
				secret:   []byte("secret"),
				database: mdb,
			}

			req, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(tt.body))
			if tt.isAuthorized {
				const contextKey auth.UserKey = "username"
				ctx := context.WithValue(req.Context(), contextKey, "user")
				req = req.WithContext(ctx)
			}
			if tt.userExists {
				mdb.EXPECT().GetUserID(req.Context(), "user").Return(1, nil)
			} else {
				mdb.EXPECT().GetUserID(req.Context(), "user").Return(0, &db.UserNotFoundError{Username: "user"})
			}

			if tt.keyExists {
				mdb.EXPECT().InsertCreditCard(req.Context(), 1, tt.item).Return(&db.KeyExistsError{Key: "111"})
			} else {
				mdb.EXPECT().InsertCreditCard(req.Context(), 1, tt.item).Return(nil)
			}
			w := httptest.NewRecorder()
			h.HandleStoreCreditCard(w, req)
			assert.Equal(t, w.Code, tt.expectedStatusCode)
		})
	}
}

func TestHandlerSet_HandleStoreBinaryItem(t *testing.T) {

	tests := []struct {
		name               string
		isAuthorized       bool
		keyExists          bool
		userExists         bool
		body               []byte
		item               types.BinaryItem
		expectedStatusCode int
	}{
		{"ok", true, false, true, binaryData, binaryItem, http.StatusCreated},
		{"notAuthorized", false, false, false, binaryData, binaryItem, http.StatusUnauthorized},
		{"userNotExists", true, false, false, binaryData, binaryItem, http.StatusUnauthorized},
		{"keyExists", true, true, true, binaryData, binaryItem, http.StatusConflict},
		{"badData", true, false, true, []byte("bad"), binaryItem, http.StatusBadRequest},
	}
	for _, tt := range tests {
		mdb := &MockDatabase{}

		t.Run(tt.name, func(t *testing.T) {

			h := &HandlerSet{
				secret:   []byte("secret"),
				database: mdb,
			}

			req, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(tt.body))
			if tt.isAuthorized {
				const contextKey auth.UserKey = "username"
				ctx := context.WithValue(req.Context(), contextKey, "user")
				req = req.WithContext(ctx)
			}

			req.Header.Set("Content-Type", "multipart/related; boundary=56a7182d3cfb7e97d66458cd95b042fbd704c16ef53b72fa871890d92f15")
			if tt.userExists {
				mdb.EXPECT().GetUserID(req.Context(), "user").Return(1, nil)
			} else {
				mdb.EXPECT().GetUserID(req.Context(), "user").Return(0, &db.UserNotFoundError{Username: "user"})
			}

			if tt.keyExists {
				mdb.EXPECT().InsertBinaryData(req.Context(), 1, tt.item).Return(&db.KeyExistsError{Key: "111"})
			} else {
				mdb.EXPECT().InsertBinaryData(req.Context(), 1, tt.item).Return(nil)
			}
			w := httptest.NewRecorder()
			h.HandleStoreBinaryItem(w, req)
			assert.Equal(t, w.Code, tt.expectedStatusCode)
		})
	}
}

func TestHandlerSet_HandleUpdateBinaryItem(t *testing.T) {

	tests := []struct {
		name               string
		isAuthorized       bool
		keyExists          bool
		userExists         bool
		body               []byte
		item               types.BinaryItem
		expectedStatusCode int
	}{
		{"ok", true, true, true, binaryData, binaryItem, http.StatusOK},
		{"notAuthorized", false, true, false, binaryData, binaryItem, http.StatusUnauthorized},
		{"userNotExists", true, true, false, binaryData, binaryItem, http.StatusUnauthorized},
		{"keyNotExists", true, false, true, binaryData, binaryItem, http.StatusNotFound},
		{"badData", true, true, true, []byte("bad"), binaryItem, http.StatusBadRequest},
	}
	for _, tt := range tests {
		mdb := &MockDatabase{}

		t.Run(tt.name, func(t *testing.T) {

			h := &HandlerSet{
				secret:   []byte("secret"),
				database: mdb,
			}

			req, _ := http.NewRequest(http.MethodPut, "", bytes.NewBuffer(tt.body))
			if tt.isAuthorized {
				const contextKey auth.UserKey = "username"
				ctx := context.WithValue(req.Context(), contextKey, "user")
				req = req.WithContext(ctx)
			}

			req.Header.Set("Content-Type", "multipart/related; boundary=56a7182d3cfb7e97d66458cd95b042fbd704c16ef53b72fa871890d92f15")
			if tt.userExists {
				mdb.EXPECT().GetUserID(req.Context(), "user").Return(1, nil)
			} else {
				mdb.EXPECT().GetUserID(req.Context(), "user").Return(0, &db.UserNotFoundError{Username: "user"})
			}

			if tt.keyExists {
				mdb.EXPECT().UpdateBinaryData(req.Context(), 1, tt.item).Return(nil)
			} else {
				mdb.EXPECT().UpdateBinaryData(req.Context(), 1, tt.item).Return(&db.KeyNotFoundError{Key: "111"})
			}
			w := httptest.NewRecorder()
			h.HandleUpdateBinaryItem(w, req)
			assert.Equal(t, w.Code, tt.expectedStatusCode)
		})
	}
}

func TestHandlerSet_HandleDownloadBinaryItem(t *testing.T) {

	tests := []struct {
		name               string
		isAuthorized       bool
		keyExists          bool
		userExists         bool
		expectedStatusCode int
		expectedBody       []byte
	}{
		{"ok", true, true, true, http.StatusOK, []byte("test")},
		{"notAuthorized", false, true, true, http.StatusUnauthorized, []byte("Something went wrong\n")},
		{"userNotExists", true, true, false, http.StatusUnauthorized, []byte("User not found\n")},
		{"notExists", true, false, true, http.StatusNotFound, []byte("Not found\n")},
	}
	for _, tt := range tests {
		mdb := &MockDatabase{}
		t.Run(tt.name, func(t *testing.T) {
			h := &HandlerSet{
				secret:   []byte("secret"),
				database: mdb,
			}
			req, _ := http.NewRequest(http.MethodGet, "/api/item/binary/111/download", nil)
			if tt.isAuthorized {
				const contextKey auth.UserKey = "username"
				ctx := context.WithValue(req.Context(), contextKey, "user")
				req = req.WithContext(ctx)
			}
			req.SetPathValue("key", "111")

			if tt.userExists {
				mdb.EXPECT().GetUserID(req.Context(), "user").Return(1, nil)
			} else {
				mdb.EXPECT().GetUserID(req.Context(), "user").Return(0, &db.UserNotFoundError{Username: "user"})
			}

			if tt.keyExists {
				mdb.EXPECT().GetBinaryData(req.Context(), 1, "111").Return([]byte("test"), nil)
			} else {
				mdb.EXPECT().GetBinaryData(req.Context(), 1, "111").Return([]byte{}, &db.KeyNotFoundError{Key: "111"})
			}
			w := httptest.NewRecorder()
			h.HandleDownloadBinaryItem(w, req)
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, string(tt.expectedBody), w.Body.String())
		})
	}
}

func TestHandlerSet_HandleItemList(t *testing.T) {
	tests := []struct {
		name           string
		isAuthorized   bool
		userExists     bool
		page           int
		limit          int
		expectedCode   int
		expectedResult []byte
	}{
		{"ok", true, true, 1, 2, http.StatusOK, []byte(`[{"Id":0,"key":"1","info":"","type":"binary"},{"Id":0,"key":"2","info":"","type":"text"}]`)},
		{"empty", true, true, 1, 2, http.StatusOK, []byte(`[]`)},
		{"notAuthorized", false, true, 1, 2, http.StatusUnauthorized, []byte("Something went wrong\n")},
		{"userNotExists", true, false, 1, 2, http.StatusUnauthorized, []byte("User not found\n")},
	}
	for _, tt := range tests {
		mdb := &MockDatabase{}
		t.Run(tt.name, func(t *testing.T) {
			h := &HandlerSet{
				secret:   []byte("secret"),
				database: mdb,
			}

			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/item/list?page=%d&limit=%d", tt.page, tt.limit), nil)
			if tt.isAuthorized {
				const contextKey auth.UserKey = "username"
				ctx := context.WithValue(req.Context(), contextKey, "user")
				req = req.WithContext(ctx)
			}
			req.SetPathValue("key", "111")

			if tt.userExists {
				mdb.EXPECT().GetUserID(req.Context(), "user").Return(1, nil)
			} else {
				mdb.EXPECT().GetUserID(req.Context(), "user").Return(0, &db.UserNotFoundError{Username: "user"})
			}
			if tt.name != "empty" {
				mdb.EXPECT().GetItems(req.Context(), 1, tt.limit, (tt.page-1)*tt.limit).Return([]types.Item{types.Item{Type: "binary", Key: "1"}, types.Item{Type: "text", Key: "2"}}, nil)
			} else {
				mdb.EXPECT().GetItems(req.Context(), 1, tt.limit, (tt.page-1)*tt.limit).Return([]types.Item{}, nil)
			}
			w := httptest.NewRecorder()
			h.HandleItemList(w, req)
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, string(tt.expectedResult), w.Body.String())
		})
	}
}

func TestHandlerSet_HandleGetItem(t *testing.T) {

	tests := []struct {
		name           string
		itemType       types.ItemType
		itemMeta       types.Item
		expextedStatus int
		expectedBody   string
	}{
		{"logopass", types.TypeLogoPass, logopassItem.Item, http.StatusOK, `{"item":{"Id":0,"key":"111","info":"","type":"logopass"},"data":{"login":"112","password":"222"}}`},
		{"text", types.TypeText, textItem.Item, http.StatusOK, `{"item":{"Id":0,"key":"111","info":"","type":"text"},"data":"text"}`},
		{"credit card", types.TypeCreditCard, creditCardItem.Item, http.StatusOK, `{"item":{"Id":0,"key":"111","info":"","type":"credit_card"},"data":{"number":"1","valid_month":"1","valid_year":"2000","name":"1","cvc":"1","ValidDate":"0001-01-01T00:00:00Z"}}`},
		{"binary", types.TypeBinary, types.Item{Key: "111", Type: types.TypeBinary}, http.StatusOK, `{"item":{"Id":0,"key":"111","info":"","type":"binary"},"data":""}`},
		{"not exists", "", types.Item{}, http.StatusNotFound, "Not found\n"},
	}
	for _, tt := range tests {
		mdb := &MockDatabase{}
		t.Run(tt.name, func(t *testing.T) {
			h := &HandlerSet{
				secret:   []byte("secret"),
				database: mdb,
			}
			req, _ := http.NewRequest(http.MethodGet, "/api/item/", nil)

			const contextKey auth.UserKey = "username"
			ctx := context.WithValue(req.Context(), contextKey, "user")
			req = req.WithContext(ctx)

			req.SetPathValue("key", "111")
			w := httptest.NewRecorder()

			mdb.EXPECT().GetUserID(req.Context(), "user").Return(1, nil)

			if tt.name != "not exists" {
				mdb.EXPECT().GetItem(req.Context(), 1, "111").Return(&tt.itemMeta, nil)
			} else {
				mdb.EXPECT().GetItem(req.Context(), 1, "111").Return(nil, &db.KeyNotFoundError{Key: "111"})
			}
			switch tt.itemType {
			case types.TypeText:
				mdb.EXPECT().GetText(req.Context(), tt.itemMeta.Id).Return(&textItem.Data, nil)
			case types.TypeCreditCard:
				mdb.EXPECT().GetCreditCard(req.Context(), tt.itemMeta.Id).Return(creditCardItem.Data, nil)
			case types.TypeLogoPass:
				mdb.EXPECT().GetLogoPass(req.Context(), tt.itemMeta.Id).Return(logopassItem.Data, nil)
			}

			h.HandleGetItem(w, req)
			assert.Equal(t, tt.expextedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestHandlerSet_HandleDeleteItem(t *testing.T) {

	tests := []struct {
		name           string
		exists         bool
		expectedStatus int
	}{
		{"exists", true, http.StatusOK},
	}
	for _, tt := range tests {
		mdb := &MockDatabase{}
		t.Run(tt.name, func(t *testing.T) {
			h := &HandlerSet{
				secret:   []byte("secret"),
				database: mdb,
			}
			req, _ := http.NewRequest(http.MethodDelete, "/api/item/", nil)

			const contextKey auth.UserKey = "username"
			ctx := context.WithValue(req.Context(), contextKey, "user")
			req = req.WithContext(ctx)

			req.SetPathValue("key", "111")
			mdb.EXPECT().GetUserID(req.Context(), "user").Return(1, nil)

			mdb.EXPECT().DeleteItem(req.Context(), 1, "111").Return(nil)

			w := httptest.NewRecorder()
			h.HandleDeleteItem(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

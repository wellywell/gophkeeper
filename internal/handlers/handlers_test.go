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
				assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6InVzZXIifQ.2MU1vpPXOY4ypHv1bPVpTqOG0zQmk4-ZD8Qoze4xVKg", w.Header().Get("X-Auth-Token"))
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
				assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6InVzZXIifQ.2MU1vpPXOY4ypHv1bPVpTqOG0zQmk4-ZD8Qoze4xVKg", w.Header().Get("X-Auth-Token"))
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
		{"ok", true, false, true, []byte(`{"item": {"type": "logopass", "key": "111"}, "data": {"login": "112", "password": "222"}}`), types.LoginPasswordItem{Item: types.Item{Key: "111", Type: types.TypeLogoPass}, Data: &types.LoginPassword{Login: "112", Password: "222"}}, http.StatusCreated},
		{"notAuthorized", false, false, false, []byte(`{"item": {"type": "logopass", "key": "111"}, "data": {"login": "112", "password": "222"}}`), types.LoginPasswordItem{Item: types.Item{Key: "111", Type: types.TypeLogoPass}, Data: &types.LoginPassword{Login: "112", Password: "222"}}, http.StatusUnauthorized},
		{"userNotExists", true, false, false, []byte(`{"item": {"type": "logopass", "key": "111"}, "data": {"login": "112", "password": "222"}}`), types.LoginPasswordItem{Item: types.Item{Key: "111", Type: types.TypeLogoPass}, Data: &types.LoginPassword{Login: "112", Password: "222"}}, http.StatusUnauthorized},
		{"keyExists", true, true, true, []byte(`{"item": {"type": "logopass", "key": "111"}, "data": {"login": "112", "password": "222"}}`), types.LoginPasswordItem{Item: types.Item{Key: "111", Type: types.TypeLogoPass}, Data: &types.LoginPassword{Login: "112", Password: "222"}}, http.StatusConflict},
		{"badData", true, false, true, []byte(`"wrong"`), types.LoginPasswordItem{Item: types.Item{Key: "111", Type: types.TypeLogoPass}, Data: &types.LoginPassword{Login: "112", Password: "222"}}, http.StatusBadRequest},
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
		{"ok", true, true, true, []byte(`{"item": {"type": "logopass", "key": "111"}, "data": {"login": "112", "password": "222"}}`), types.LoginPasswordItem{Item: types.Item{Key: "111", Type: types.TypeLogoPass}, Data: &types.LoginPassword{Login: "112", Password: "222"}}, http.StatusOK},
		{"notAuthorized", false, true, false, []byte(`{"item": {"type": "logopass", "key": "111"}, "data": {"login": "112", "password": "222"}}`), types.LoginPasswordItem{Item: types.Item{Key: "111", Type: types.TypeLogoPass}, Data: &types.LoginPassword{Login: "112", Password: "222"}}, http.StatusUnauthorized},
		{"userNotExists", true, true, false, []byte(`{"item": {"type": "logopass", "key": "111"}, "data": {"login": "112", "password": "222"}}`), types.LoginPasswordItem{Item: types.Item{Key: "111", Type: types.TypeLogoPass}, Data: &types.LoginPassword{Login: "112", Password: "222"}}, http.StatusUnauthorized},
		{"keyNotExists", true, false, true, []byte(`{"item": {"type": "logopass", "key": "111"}, "data": {"login": "112", "password": "222"}}`), types.LoginPasswordItem{Item: types.Item{Key: "111", Type: types.TypeLogoPass}, Data: &types.LoginPassword{Login: "112", Password: "222"}}, http.StatusNotFound},
		{"badData", true, true, true, []byte(`"wrong"`), types.LoginPasswordItem{Item: types.Item{Key: "111", Type: types.TypeLogoPass}, Data: &types.LoginPassword{Login: "112", Password: "222"}}, http.StatusBadRequest},
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
		{"ok", true, true, true, []byte(`{"item": {"type": "credit_card", "key": "111"}, "data": {"cvc": "1", "number":"1", "name":"1", "valid_month": "1", "valid_year": "2000"}}`), types.CreditCardItem{Item: types.Item{Key: "111", Type: types.TypeCreditCard}, Data: &types.CreditCardData{Number: "1", CVC: "1", Name: "1", ValidMonth: "1", ValidYear: "2000"}}, http.StatusOK},
		{"notAuthorized", false, true, false, []byte(`{"item": {"type": "credit_card", "key": "111"}, "data": {"cvc": "1", "number":"1", "name":"1",  "valid_month": "1", "valid_year": "2000"}}`), types.CreditCardItem{Item: types.Item{Key: "111", Type: types.TypeCreditCard}, Data: &types.CreditCardData{Number: "1", CVC: "1", Name: "1", ValidMonth: "1", ValidYear: "2000"}}, http.StatusUnauthorized},
		{"userNotExists", true, true, false, []byte(`{"item": {"type": "credit_card", "key": "111"}, "data": {"cvc": "1", "number":"1", "name":"1",  "valid_month": "1", "valid_year": "2000"}}`), types.CreditCardItem{Item: types.Item{Key: "111", Type: types.TypeCreditCard}, Data: &types.CreditCardData{Number: "1", CVC: "1", Name: "1", ValidMonth: "1", ValidYear: "2000"}}, http.StatusUnauthorized},
		{"keyNotExists", true, false, true, []byte(`{"item": {"type": "credit_card", "key": "111"}, "data": {"cvc": "1", "number":"1", "name":"1",  "valid_month": "1", "valid_year": "2000"}}`), types.CreditCardItem{Item: types.Item{Key: "111", Type: types.TypeCreditCard}, Data: &types.CreditCardData{Number: "1", CVC: "1", Name: "1", ValidMonth: "1", ValidYear: "2000"}}, http.StatusNotFound},
		{"badData", true, true, true, []byte(`"wrong"`), types.CreditCardItem{Item: types.Item{Key: "111", Type: types.TypeCreditCard}, Data: &types.CreditCardData{}}, http.StatusBadRequest},
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
		{"ok", true, true, true, []byte(`{"item": {"type": "text", "key": "111"}, "data": "text"}`), types.TextItem{Item: types.Item{Key: "111", Type: types.TypeText}, Data: types.TextData("text")}, http.StatusOK},
		{"notAuthorized", false, true, false, []byte(`{"item": {"type": "text", "key": "111"}, "data": "text"}`), types.TextItem{Item: types.Item{Key: "111", Type: types.TypeText}, Data: types.TextData("text")}, http.StatusUnauthorized},
		{"userNotExists", true, true, false, []byte(`{"item": {"type": "text", "key": "111"}, "data": "text"}`), types.TextItem{Item: types.Item{Key: "111", Type: types.TypeText}, Data: types.TextData("text")}, http.StatusUnauthorized},
		{"keyNotExists", true, false, true, []byte(`{"item": {"type": "text", "key": "111"}, "data": "text"}`), types.TextItem{Item: types.Item{Key: "111", Type: types.TypeText}, Data: types.TextData("text")}, http.StatusNotFound},
		{"badData", true, true, true, []byte(`"wrong"`), types.TextItem{Item: types.Item{Key: "111", Type: types.TypeLogoPass}, Data: types.TextData("text")}, http.StatusBadRequest},
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
		{"ok", true, false, true, []byte(`{"item": {"type": "text", "key": "111"}, "data": "text"}`), types.TextItem{Item: types.Item{Key: "111", Type: types.TypeText}, Data: types.TextData("text")}, http.StatusCreated},
		{"notAuthorized", false, false, false, []byte(`{"item": {"type": "text", "key": "111"}, "data": "text"}`), types.TextItem{Item: types.Item{Key: "111", Type: types.TypeText}, Data: types.TextData("text")}, http.StatusUnauthorized},
		{"userNotExists", true, false, false, []byte(`{"item": {"type": "text", "key": "111"}, "data": "text"}`), types.TextItem{Item: types.Item{Key: "111", Type: types.TypeText}, Data: types.TextData("text")}, http.StatusUnauthorized},
		{"keyExists", true, true, true, []byte(`{"item": {"type": "text", "key": "111"}, "data": "text"}`), types.TextItem{Item: types.Item{Key: "111", Type: types.TypeText}, Data: types.TextData("text")}, http.StatusConflict},
		{"badData", true, false, true, []byte(`"wrong"`), types.TextItem{Item: types.Item{Key: "111", Type: types.TypeLogoPass}, Data: types.TextData("text")}, http.StatusBadRequest},
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
	type fields struct {
		secret   []byte
		database Database
	}
	type args struct {
		w   http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HandlerSet{
				secret:   tt.fields.secret,
				database: tt.fields.database,
			}
			h.HandleStoreCreditCard(tt.args.w, tt.args.req)
		})
	}
}

func TestHandlerSet_HandleStoreBinaryItem(t *testing.T) {
	type fields struct {
		secret   []byte
		database Database
	}
	type args struct {
		w   http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HandlerSet{
				secret:   tt.fields.secret,
				database: tt.fields.database,
			}
			h.HandleStoreBinaryItem(tt.args.w, tt.args.req)
		})
	}
}

func TestHandlerSet_HandleUpdateBinaryItem(t *testing.T) {
	type fields struct {
		secret   []byte
		database Database
	}
	type args struct {
		w   http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HandlerSet{
				secret:   tt.fields.secret,
				database: tt.fields.database,
			}
			h.HandleUpdateBinaryItem(tt.args.w, tt.args.req)
		})
	}
}

func TestHandlerSet_HandleDownloadBinaryItem(t *testing.T) {
	type fields struct {
		secret   []byte
		database Database
	}
	type args struct {
		w   http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HandlerSet{
				secret:   tt.fields.secret,
				database: tt.fields.database,
			}
			h.HandleDownloadBinaryItem(tt.args.w, tt.args.req)
		})
	}
}

func TestHandlerSet_HandleItemList(t *testing.T) {
	type fields struct {
		secret   []byte
		database Database
	}
	type args struct {
		w   http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HandlerSet{
				secret:   tt.fields.secret,
				database: tt.fields.database,
			}
			h.HandleItemList(tt.args.w, tt.args.req)
		})
	}
}

func TestHandlerSet_HandleGetItem(t *testing.T) {
	type fields struct {
		secret   []byte
		database Database
	}
	type args struct {
		w   http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HandlerSet{
				secret:   tt.fields.secret,
				database: tt.fields.database,
			}
			h.HandleGetItem(tt.args.w, tt.args.req)
		})
	}
}

func TestHandlerSet_HandleDeleteItem(t *testing.T) {
	type fields struct {
		secret   []byte
		database Database
	}

	type args struct {
		w   http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HandlerSet{
				secret:   tt.fields.secret,
				database: tt.fields.database,
			}
			h.HandleDeleteItem(tt.args.w, tt.args.req)
		})
	}
}

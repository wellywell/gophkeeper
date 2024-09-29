// Package client содержит методы для создания http-запросов от клиента на сервер
package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wellywell/gophkeeper/internal/config"
	"github.com/wellywell/gophkeeper/internal/types"
)

var conf *config.ClientConfig

func TestMain(m *testing.M) {
	code, err := runMain(m)

	if err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}

func runMain(m *testing.M) (int, error) {

	conf, _ = config.NewClientConfig()
	exitCode := m.Run()

	return exitCode, nil

}

func TestNewClient(t *testing.T) {
	_, err := NewClient(conf)
	assert.NoError(t, err)
}

func TestClient_Login(t *testing.T) {

	type args struct {
		login    string
		password string
	}
	tests := []struct {
		name       string
		args       args
		serverCode int
		want       string
		wantErr    bool
	}{
		{"ok", args{"user", "pass"}, http.StatusOK, "token", false},
		{"notOk", args{"user", "pass"}, http.StatusUnauthorized, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.want != "" {
					w.Header().Set("X-Auth-Token", tt.want)
				}
				w.WriteHeader(tt.serverCode)
				fmt.Fprintln(w, "oooo")
			}))
			defer svr.Close()

			c, _ := NewClient(conf)
			c.address = svr.URL

			got, err := c.Login(tt.args.login, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Register(t *testing.T) {
	type args struct {
		login    string
		password string
	}
	tests := []struct {
		name       string
		args       args
		serverCode int
		want       string
		wantErr    bool
	}{
		{"ok", args{"user", "pass"}, http.StatusOK, "token", false},
		{"notOk", args{"user", "pass"}, http.StatusUnauthorized, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.want != "" {
					w.Header().Set("X-Auth-Token", tt.want)
				}
				w.WriteHeader(tt.serverCode)
			}))
			defer svr.Close()

			c, _ := NewClient(conf)
			c.address = svr.URL
			got, err := c.Register(tt.args.login, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.Register() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_CreateBinaryItem(t *testing.T) {

	type args struct {
		data    []byte
		headers map[string]string
	}
	tests := []struct {
		name        string
		args        args
		respone     string
		headerName  string
		headerValue string
	}{
		{"base", args{[]byte("123"), map[string]string{"my-header": "header_content"}}, "OK", "my-header", "header_content"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				data, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.Equal(t, tt.args.data, data)
				assert.Equal(t, tt.headerValue, r.Header.Get(tt.headerName))
				fmt.Fprint(w, tt.respone)
			}))
			defer svr.Close()

			c, _ := NewClient(conf)
			c.address = svr.URL
			got, err := c.CreateBinaryItem(tt.args.data, tt.args.headers)
			assert.NoError(t, err)

			data, err := io.ReadAll(got.Body)
			assert.NoError(t, err)
			defer got.Body.Close()

			if !reflect.DeepEqual(string(data), string(tt.respone)) {
				t.Errorf("Client.CreateBinaryItem() = %s, want %s", data, tt.respone)
			}
		})
	}
}

func TestClient_UpdateBinaryItem(t *testing.T) {
	type args struct {
		data    []byte
		headers map[string]string
	}
	tests := []struct {
		name        string
		args        args
		respone     string
		headerName  string
		headerValue string
	}{
		{"base", args{[]byte("123"), map[string]string{"my-header": "header_content"}}, "OK", "my-header", "header_content"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				data, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.Equal(t, tt.args.data, data)
				assert.Equal(t, tt.headerValue, r.Header.Get(tt.headerName))
				fmt.Fprint(w, tt.respone)
			}))
			defer svr.Close()

			c, _ := NewClient(conf)
			c.address = svr.URL
			got, err := c.UpdateBinaryItem(tt.args.data, tt.args.headers)
			assert.NoError(t, err)

			data, err := io.ReadAll(got.Body)
			assert.NoError(t, err)
			defer got.Body.Close()
			if !reflect.DeepEqual(string(data), string(tt.respone)) {
				t.Errorf("Client.UpdateBinaryItem() = %v, want %v", got, tt.respone)
			}
		})
	}
}

func TestClient_CreateLoginPasswordItem(t *testing.T) {
	type args struct {
		data    []byte
		headers map[string]string
	}
	tests := []struct {
		name        string
		args        args
		respone     string
		headerName  string
		headerValue string
	}{
		{"base", args{[]byte("123"), map[string]string{"my-header": "header_content"}}, "OK", "my-header", "header_content"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				data, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.Equal(t, tt.args.data, data)
				assert.Equal(t, tt.headerValue, r.Header.Get(tt.headerName))
				fmt.Fprint(w, tt.respone)
			}))
			defer svr.Close()

			c, _ := NewClient(conf)
			c.address = svr.URL
			got, err := c.CreateLoginPasswordItem(tt.args.data, tt.args.headers)
			assert.NoError(t, err)

			data, err := io.ReadAll(got.Body)
			assert.NoError(t, err)
			defer got.Body.Close()
			if !reflect.DeepEqual(string(data), string(tt.respone)) {
				t.Errorf("Client.CreateLoginPasswordItem() = %v, want %v", got, tt.respone)
			}
		})
	}
}

func TestClient_CreateCreditCardItem(t *testing.T) {
	type args struct {
		data    []byte
		headers map[string]string
	}
	tests := []struct {
		name        string
		args        args
		respone     string
		headerName  string
		headerValue string
	}{
		{"base", args{[]byte("123"), map[string]string{"my-header": "header_content"}}, "OK", "my-header", "header_content"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				data, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.Equal(t, tt.args.data, data)
				assert.Equal(t, tt.headerValue, r.Header.Get(tt.headerName))
				fmt.Fprint(w, tt.respone)
			}))
			defer svr.Close()

			c, _ := NewClient(conf)
			c.address = svr.URL
			got, err := c.CreateCreditCardItem(tt.args.data, tt.args.headers)
			assert.NoError(t, err)

			data, err := io.ReadAll(got.Body)
			assert.NoError(t, err)
			defer got.Body.Close()
			if !reflect.DeepEqual(string(data), string(tt.respone)) {
				t.Errorf("Client.CreateCreditCardItem() = %v, want %v", got, tt.respone)
			}
		})
	}
}

func TestClient_CreateTextItem(t *testing.T) {
	type args struct {
		data    []byte
		headers map[string]string
	}
	tests := []struct {
		name        string
		args        args
		respone     string
		headerName  string
		headerValue string
	}{
		{"base", args{[]byte("123"), map[string]string{"my-header": "header_content"}}, "OK", "my-header", "header_content"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				data, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.Equal(t, tt.args.data, data)
				assert.Equal(t, tt.headerValue, r.Header.Get(tt.headerName))
				fmt.Fprint(w, tt.respone)
			}))
			defer svr.Close()

			c, _ := NewClient(conf)
			c.address = svr.URL
			got, err := c.CreateTextItem(tt.args.data, tt.args.headers)
			assert.NoError(t, err)

			data, err := io.ReadAll(got.Body)
			assert.NoError(t, err)
			defer got.Body.Close()
			if !reflect.DeepEqual(string(data), string(tt.respone)) {
				t.Errorf("Client.CreateTextItem() = %v, want %v", got, tt.respone)
			}
		})
	}
}

func TestClient_GetItem(t *testing.T) {
	type args struct {
		token string
		key   string
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		respCode int
		respBody string
	}{
		{"ok", args{"token", "111"}, false, http.StatusOK, "something"},
		{"notOK", args{"token", "111"}, true, http.StatusInternalServerError, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/api/item/111", r.URL.Path)
				assert.Equal(t, r.Header.Get("X-Auth-Token"), tt.args.token)
				w.WriteHeader(tt.respCode)
				fmt.Fprint(w, tt.respBody)
			}))
			defer svr.Close()

			c, _ := NewClient(conf)
			c.address = svr.URL
			gotData, err := c.GetItem(tt.args.token, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(string(gotData), string(tt.respBody)) {
				t.Errorf("Client.GetItem() = %v, want %v", gotData, tt.respBody)
			}
		})
	}
}

func TestClient_SeeRecords(t *testing.T) {
	type args struct {
		token    string
		pass     string
		page     int
		pageSize int
	}
	tests := []struct {
		name         string
		args         args
		want         []types.Item
		wantErr      bool
		responseCode int
		responseBody string
	}{
		{"ok", args{"token", "pass", 1, 2}, []types.Item{{Key: "111", Type: "text"}, {Key: "222", Type: "binary"}}, false, http.StatusOK, `[{"key": "111", "type": "text", "info":""}, {"key": "222", "type": "binary", "info":""}]`},
		{"notOk", args{"token", "pass", 1, 2}, []types.Item{}, true, http.StatusInternalServerError, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, fmt.Sprint(tt.args.page), r.URL.Query().Get("page"))
				assert.Equal(t, fmt.Sprint(tt.args.pageSize), r.URL.Query().Get("limit"))
				assert.Equal(t, r.Header.Get("X-Auth-Token"), tt.args.token)
				w.WriteHeader(tt.responseCode)
				fmt.Fprint(w, tt.responseBody)
			}))
			defer svr.Close()

			c, _ := NewClient(conf)
			c.address = svr.URL
			got, err := c.SeeRecords(tt.args.token, tt.args.pass, tt.args.page, tt.args.pageSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SeeRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.SeeRecords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_DownloadBinaryData(t *testing.T) {

	bs := types.BinaryData([]byte("some text some text some text some text some text some text some text some text some text some text some text"))

	expectEncrypted := []byte{
		125, 249, 121, 247, 12, 220, 4, 198, 224, 186, 123, 29, 124, 254, 76, 113, 138, 167, 125, 67, 84, 18, 124, 206, 229,
		3, 93, 63, 235, 126, 7, 61, 148, 82, 244, 239, 85, 76, 248, 149, 181, 114, 60, 77, 213, 94, 48, 133, 244, 164, 232,
		86, 148, 10, 220, 55, 247, 112, 80, 180, 113, 164, 147, 124, 74, 153, 194, 80, 214, 101, 121, 165, 101, 224, 241, 232,
		181, 5, 43, 23, 88, 199, 221, 189, 246, 211, 89, 156, 242, 118, 93, 182, 216, 94, 48, 142, 195, 190, 55, 178, 93, 218,
		112, 73, 183, 186, 122, 239, 140}

	type args struct {
		token string
		pass  string
		key   string
	}
	tests := []struct {
		name     string
		args     args
		wantData []byte
		wantErr  bool
		respBody []byte
		respCode int
	}{
		{"ok", args{"token", "secret", "111"}, bs, false, expectEncrypted, http.StatusOK},
		{"ok", args{"token", "secret", "111"}, bs, true, expectEncrypted, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/api/item/binary/111/download", r.URL.Path)
				assert.Equal(t, r.Header.Get("X-Auth-Token"), tt.args.token)
				w.WriteHeader(tt.respCode)
				_, _ = w.Write(tt.respBody)
			}))
			defer svr.Close()

			c, _ := NewClient(conf)
			c.address = svr.URL
			gotData, err := c.DownloadBinaryData(tt.args.token, tt.args.pass, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.DownloadBinaryData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual([]byte(gotData), []byte(tt.wantData)) {
				t.Errorf("Client.DownloadBinaryData() = %v, \nwant %v", gotData, tt.wantData)
			}
		})
	}
}

func TestClient_UpdateLogoPassData(t *testing.T) {
	type args struct {
		data    []byte
		headers map[string]string
	}
	tests := []struct {
		name        string
		args        args
		respone     string
		headerName  string
		headerValue string
	}{
		{"base", args{[]byte("123"), map[string]string{"my-header": "header_content"}}, "OK", "my-header", "header_content"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				data, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.Equal(t, tt.args.data, data)
				assert.Equal(t, tt.headerValue, r.Header.Get(tt.headerName))
				fmt.Fprint(w, tt.respone)
			}))
			defer svr.Close()

			c, _ := NewClient(conf)
			c.address = svr.URL
			got, err := c.UpdateLogoPassData(tt.args.data, tt.args.headers)
			assert.NoError(t, err)

			data, err := io.ReadAll(got.Body)
			assert.NoError(t, err)
			defer got.Body.Close()
			if !reflect.DeepEqual(string(data), string(tt.respone)) {
				t.Errorf("Client.UpdateLogoPassData() = %v, want %v", got, tt.respone)
			}
		})
	}
}

func TestClient_UpdateCreditCardData(t *testing.T) {
	type args struct {
		data    []byte
		headers map[string]string
	}
	tests := []struct {
		name        string
		args        args
		respone     string
		headerName  string
		headerValue string
	}{
		{"base", args{[]byte("123"), map[string]string{"my-header": "header_content"}}, "OK", "my-header", "header_content"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				data, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.Equal(t, tt.args.data, data)
				assert.Equal(t, tt.headerValue, r.Header.Get(tt.headerName))
				fmt.Fprint(w, tt.respone)
			}))
			defer svr.Close()

			c, _ := NewClient(conf)
			c.address = svr.URL
			got, err := c.UpdateCreditCardData(tt.args.data, tt.args.headers)
			assert.NoError(t, err)

			data, err := io.ReadAll(got.Body)
			assert.NoError(t, err)
			defer got.Body.Close()
			if !reflect.DeepEqual(string(data), string(tt.respone)) {
				t.Errorf("Client.UpdateCreditCardData() = %v, want %v", got, tt.respone)
			}
		})
	}
}

func TestClient_UpdateTextData(t *testing.T) {
	type args struct {
		data    []byte
		headers map[string]string
	}
	tests := []struct {
		name        string
		args        args
		respone     string
		headerName  string
		headerValue string
	}{
		{"base", args{[]byte("123"), map[string]string{"my-header": "header_content"}}, "OK", "my-header", "header_content"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				data, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.Equal(t, tt.args.data, data)
				assert.Equal(t, tt.headerValue, r.Header.Get(tt.headerName))
				fmt.Fprint(w, tt.respone)
			}))
			defer svr.Close()

			c, _ := NewClient(conf)
			c.address = svr.URL
			got, err := c.UpdateTextData(tt.args.data, tt.args.headers)
			assert.NoError(t, err)

			data, err := io.ReadAll(got.Body)
			assert.NoError(t, err)
			defer got.Body.Close()
			if !reflect.DeepEqual(string(data), string(tt.respone)) {
				t.Errorf("Client.UpdateTextData() = %v, want %v", got, tt.respone)
			}
		})
	}
}

func TestClient_DeleteItem(t *testing.T) {

	type args struct {
		token string
		key   string
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		respCode int
	}{
		{"ok", args{"token", "111"}, false, http.StatusOK},
		{"notOK", args{"token", "111"}, true, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/api/item/111", r.URL.Path)
				assert.Equal(t, r.Header.Get("X-Auth-Token"), tt.args.token)
				w.WriteHeader(tt.respCode)
			}))
			defer svr.Close()

			c, _ := NewClient(conf)
			c.address = svr.URL
			if err := c.DeleteItem(tt.args.token, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateItem_Text(t *testing.T) {
	type args struct {
		token   string
		pass    string
		newItem types.GenericItem[*types.TextData]
	}

	textData := "text"

	textType := types.TextData(textData)

	tests := []struct {
		name     string
		args     args
		wantErr  bool
		respCode int
	}{
		{"ok", args{"token", "secret", types.GenericItem[*types.TextData]{Data: &textType}}, false, http.StatusOK},
		{"notok", args{"token", "secret", types.GenericItem[*types.TextData]{Data: &textType}}, true, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				data, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				text := types.TextItem{}

				err = json.Unmarshal(data, &text)
				assert.NoError(t, err)
				assert.NotEqual(t, text.Data, tt.args.newItem.Data)

				err = text.Data.Decrypt(tt.args.pass)
				assert.NoError(t, err)

				assert.Equal(t, string(text.Data), textData)

				w.WriteHeader(tt.respCode)
			}))
			defer svr.Close()

			c, _ := NewClient(conf)
			c.address = svr.URL
			if err := UpdateItem(tt.args.token, tt.args.pass, tt.args.newItem, c.UpdateTextData); (err != nil) != tt.wantErr {
				t.Errorf("UpdateItem() error = %v, wantErr %v", err, tt.wantErr)
			}
			textType = types.TextData(textData)
		})
	}
}

func TestUpdateItem_Logopass(t *testing.T) {
	type args struct {
		token   string
		pass    string
		newItem types.GenericItem[*types.LoginPassword]
	}

	tests := []struct {
		name     string
		args     args
		wantErr  bool
		respCode int
	}{
		{"ok", args{"token", "secret", types.GenericItem[*types.LoginPassword]{Data: &types.LoginPassword{Login: "1"}}}, false, http.StatusOK},
		{"notok", args{"token", "secret", types.GenericItem[*types.LoginPassword]{Data: &types.LoginPassword{Login: "1"}}}, true, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				data, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				text := types.LoginPasswordItem{}

				err = json.Unmarshal(data, &text)
				assert.NoError(t, err)
				assert.NotEqual(t, text.Data.Login, "1")

				err = text.Data.Decrypt(tt.args.pass)
				assert.NoError(t, err)

				assert.Equal(t, "1", text.Data.Login)

				w.WriteHeader(tt.respCode)
			}))
			defer svr.Close()

			c, _ := NewClient(conf)
			c.address = svr.URL
			if err := UpdateItem(tt.args.token, tt.args.pass, tt.args.newItem, c.UpdateLogoPassData); (err != nil) != tt.wantErr {
				t.Errorf("UpdateItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateItem_CreditCard(t *testing.T) {
	type args struct {
		token   string
		pass    string
		newItem types.GenericItem[*types.CreditCardData]
	}

	tests := []struct {
		name     string
		args     args
		wantErr  bool
		respCode int
	}{
		{"ok", args{"token", "secret", types.GenericItem[*types.CreditCardData]{Data: &types.CreditCardData{CVC: "1"}}}, false, http.StatusOK},
		{"notok", args{"token", "secret", types.GenericItem[*types.CreditCardData]{Data: &types.CreditCardData{CVC: "1"}}}, true, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				data, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				text := types.CreditCardItem{}

				err = json.Unmarshal(data, &text)
				assert.NoError(t, err)
				assert.NotEqual(t, text.Data.CVC, "1")

				err = text.Data.Decrypt(tt.args.pass)
				assert.NoError(t, err)

				assert.Equal(t, "1", text.Data.CVC)

				w.WriteHeader(tt.respCode)
			}))
			defer svr.Close()

			c, _ := NewClient(conf)
			c.address = svr.URL
			if err := UpdateItem(tt.args.token, tt.args.pass, tt.args.newItem, c.UpdateLogoPassData); (err != nil) != tt.wantErr {
				t.Errorf("UpdateItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateItem_Binary(t *testing.T) {
	type args struct {
		token   string
		pass    string
		newItem types.GenericItem[*types.BinaryData]
	}

	textData := "text"

	binaryType := types.BinaryData([]byte(textData))

	tests := []struct {
		name     string
		args     args
		wantErr  bool
		respCode int
	}{
		{"ok", args{"token", "secret", types.GenericItem[*types.BinaryData]{Item: types.Item{Type: types.TypeBinary}, Data: &binaryType}}, false, http.StatusOK},
		{"notok", args{"token", "secret", types.GenericItem[*types.BinaryData]{Item: types.Item{Type: types.TypeBinary}, Data: &binaryType}}, true, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				data, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.Contains(t, string(data), "Content-Type: application/octet-stream")
				assert.Contains(t, string(data), "Content-Type: application/json")
				assert.Contains(t, string(data), `{"Id":0,"key":"","info":"","type":"binary"}`)

				w.WriteHeader(tt.respCode)
			}))
			defer svr.Close()

			c, _ := NewClient(conf)
			c.address = svr.URL
			if err := UpdateItem(tt.args.token, tt.args.pass, tt.args.newItem, c.UpdateBinaryItem); (err != nil) != tt.wantErr {
				t.Errorf("UpdateItem() error = %v, wantErr %v", err, tt.wantErr)
			}
			binaryType = types.BinaryData([]byte(textData))
		})
	}
}

func TestCreateItem_Text(t *testing.T) {
	type args struct {
		token   string
		pass    string
		newItem types.GenericItem[*types.TextData]
	}

	textData := "text"

	textType := types.TextData(textData)

	tests := []struct {
		name     string
		args     args
		wantErr  bool
		respCode int
	}{
		{"ok", args{"token", "secret", types.GenericItem[*types.TextData]{Data: &textType}}, false, http.StatusCreated},
		{"notok", args{"token", "secret", types.GenericItem[*types.TextData]{Data: &textType}}, true, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				data, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				text := types.TextItem{}

				err = json.Unmarshal(data, &text)
				assert.NoError(t, err)
				assert.NotEqual(t, text.Data, tt.args.newItem.Data)

				err = text.Data.Decrypt(tt.args.pass)
				assert.NoError(t, err)

				assert.Equal(t, string(text.Data), textData)

				w.WriteHeader(tt.respCode)
			}))
			defer svr.Close()

			c, _ := NewClient(conf)
			c.address = svr.URL
			if err := CreateItem(tt.args.token, tt.args.pass, tt.args.newItem, c.CreateTextItem); (err != nil) != tt.wantErr {
				t.Errorf("CreateItem() error = %v, wantErr %v", err, tt.wantErr)
			}
			textType = types.TextData(textData)
		})
	}
}

func TestCreateItem_Logopass(t *testing.T) {
	type args struct {
		token   string
		pass    string
		newItem types.GenericItem[*types.LoginPassword]
	}

	tests := []struct {
		name     string
		args     args
		wantErr  bool
		respCode int
	}{
		{"ok", args{"token", "secret", types.GenericItem[*types.LoginPassword]{Data: &types.LoginPassword{Login: "1"}}}, false, http.StatusCreated},
		{"notok", args{"token", "secret", types.GenericItem[*types.LoginPassword]{Data: &types.LoginPassword{Login: "1"}}}, true, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				data, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				text := types.LoginPasswordItem{}

				err = json.Unmarshal(data, &text)
				assert.NoError(t, err)
				assert.NotEqual(t, text.Data.Login, "1")

				err = text.Data.Decrypt(tt.args.pass)
				assert.NoError(t, err)

				assert.Equal(t, "1", text.Data.Login)

				w.WriteHeader(tt.respCode)
			}))
			defer svr.Close()

			c, _ := NewClient(conf)
			c.address = svr.URL
			if err := CreateItem(tt.args.token, tt.args.pass, tt.args.newItem, c.CreateLoginPasswordItem); (err != nil) != tt.wantErr {
				t.Errorf("CreateItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateItem_CreditCard(t *testing.T) {
	type args struct {
		token   string
		pass    string
		newItem types.GenericItem[*types.CreditCardData]
	}

	tests := []struct {
		name     string
		args     args
		wantErr  bool
		respCode int
	}{
		{"ok", args{"token", "secret", types.GenericItem[*types.CreditCardData]{Data: &types.CreditCardData{CVC: "1"}}}, false, http.StatusCreated},
		{"notok", args{"token", "secret", types.GenericItem[*types.CreditCardData]{Data: &types.CreditCardData{CVC: "1"}}}, true, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				data, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				text := types.CreditCardItem{}

				err = json.Unmarshal(data, &text)
				assert.NoError(t, err)
				assert.NotEqual(t, text.Data.CVC, "1")

				err = text.Data.Decrypt(tt.args.pass)
				assert.NoError(t, err)

				assert.Equal(t, "1", text.Data.CVC)

				w.WriteHeader(tt.respCode)
			}))
			defer svr.Close()

			c, _ := NewClient(conf)
			c.address = svr.URL
			if err := CreateItem(tt.args.token, tt.args.pass, tt.args.newItem, c.CreateCreditCardItem); (err != nil) != tt.wantErr {
				t.Errorf("CreateItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateItem_Binary(t *testing.T) {
	type args struct {
		token   string
		pass    string
		newItem types.GenericItem[*types.BinaryData]
	}

	textData := "text"

	binaryType := types.BinaryData([]byte(textData))

	tests := []struct {
		name     string
		args     args
		wantErr  bool
		respCode int
	}{
		{"ok", args{"token", "secret", types.GenericItem[*types.BinaryData]{Item: types.Item{Type: types.TypeBinary}, Data: &binaryType}}, false, http.StatusCreated},
		{"notok", args{"token", "secret", types.GenericItem[*types.BinaryData]{Item: types.Item{Type: types.TypeBinary}, Data: &binaryType}}, true, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				data, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.Contains(t, string(data), "Content-Type: application/octet-stream")
				assert.Contains(t, string(data), "Content-Type: application/json")
				assert.Contains(t, string(data), `{"Id":0,"key":"","info":"","type":"binary"}`)

				w.WriteHeader(tt.respCode)
			}))
			defer svr.Close()

			c, _ := NewClient(conf)
			c.address = svr.URL
			if err := CreateItem(tt.args.token, tt.args.pass, tt.args.newItem, c.CreateBinaryItem); (err != nil) != tt.wantErr {
				t.Errorf("CreateItem() error = %v, wantErr %v", err, tt.wantErr)
			}
			binaryType = types.BinaryData([]byte(textData))
		})
	}
}

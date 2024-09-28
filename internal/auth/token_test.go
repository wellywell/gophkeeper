package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyUser(t *testing.T) {
	type args struct {
		r      *http.Request
		secret []byte
	}

	tests := []struct {
		name    string
		args    args
		token   string
		want    string
		wantErr bool
	}{
		{"good token", args{httptest.NewRequest(http.MethodGet, "/api/", nil), []byte("secret")}, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6InVzZXIifQ.2MU1vpPXOY4ypHv1bPVpTqOG0zQmk4-ZD8Qoze4xVKg", "user", false},
		{"bad token", args{httptest.NewRequest(http.MethodGet, "/api/", nil), []byte("secret")}, "bad", "", true},
		{"wrong secret", args{httptest.NewRequest(http.MethodGet, "/api/", nil), []byte("wrong")}, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6InVzZXIifQ.2MU1vpPXOY4ypHv1bPVpTqOG0zQmk4-ZD8Qoze4xVKg", "", true},
		{"empty token", args{httptest.NewRequest(http.MethodGet, "/api/", nil), []byte("secret")}, "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.args.r.Header.Set("X-Auth-Token", tt.token)

			got, err := VerifyUser(tt.args.r, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("VerifyUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildJWTString(t *testing.T) {
	type args struct {
		user   string
		secret []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"token", args{"user", []byte("secret")}, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6InVzZXIifQ.2MU1vpPXOY4ypHv1bPVpTqOG0zQmk4-ZD8Qoze4xVKg", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildJWTString(tt.args.user, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildJWTString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BuildJWTString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	type args struct {
		tokenString string
		secret      []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"userPresent", args{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6InVzZXIifQ.2MU1vpPXOY4ypHv1bPVpTqOG0zQmk4-ZD8Qoze4xVKg", []byte("secret")}, "user", false},
		{"wrongToken", args{"wrong", []byte("secret")}, "", true},
		{"wrongSecret", args{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6InVzZXIifQ.2MU1vpPXOY4ypHv1bPVpTqOG0zQmk4-ZD8Qoze4xVKg", []byte("wrong")}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUser(tt.args.tokenString, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetToken(t *testing.T) {
	type args struct {
		username string
		w        http.ResponseWriter
		secret   []byte
	}
	w := httptest.NewRecorder()
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"set token", args{"user", w, []byte("secret")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetToken(tt.args.username, tt.args.w, tt.args.secret); (err != nil) != tt.wantErr {
				t.Errorf("SetToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			token := w.Header().Get("X-Auth-Token")
			assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6InVzZXIifQ.2MU1vpPXOY4ypHv1bPVpTqOG0zQmk4-ZD8Qoze4xVKg", token)
		})
	}
}

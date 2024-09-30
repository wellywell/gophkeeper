package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAuthenticatedUser(t *testing.T) {
	type args struct {
		req *http.Request
	}

	r := httptest.NewRequest(http.MethodGet, "/api/", nil)

	const contextKey UserKey = "username"
	ctx := context.WithValue(r.Context(), contextKey, "user")
	rWithUser := r.WithContext(ctx)

	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{"with context", args{rWithUser}, "user", true},
		{"without context", args{r}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetAuthenticatedUser(tt.args.req)
			if got != tt.want {
				t.Errorf("GetAuthenticatedUser() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetAuthenticatedUser() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

type Mockhandler struct{}

func (h *Mockhandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

type TestMiddleWare struct {
	t            *testing.T
	expectedUser string
	expectedOk   bool
}

func (m TestMiddleWare) Handle(next http.Handler) http.Handler {

	test := func(w http.ResponseWriter, r *http.Request) {

		const contextKey UserKey = "username"
		user, ok := r.Context().Value(contextKey).(string)

		assert.Equal(m.t, m.expectedUser, user)
		assert.Equal(m.t, m.expectedOk, ok)
		next.ServeHTTP(w, r)

	}
	return http.HandlerFunc(test)
}

func TestAuthenticateMiddleware_Handle(t *testing.T) {

	type fields struct {
		Secret []byte
	}
	type args struct {
		next http.Handler
	}

	r := httptest.NewRequest(http.MethodGet, "/api/", nil)

	rWithToken := httptest.NewRequest(http.MethodGet, "/api/", nil)
	rWithToken.Header.Set("X-Auth-Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6InVzZXIifQ.2MU1vpPXOY4ypHv1bPVpTqOG0zQmk4-ZD8Qoze4xVKg")

	secret := "secret"
	wrongSecret := "wrong"

	handler := &Mockhandler{}

	tests := []struct {
		name               string
		fields             fields
		args               args
		request            *http.Request
		expectedStatusCode int
		expectedUser       string
		expextedUserOK     bool
	}{
		{"good request", fields{[]byte(secret)}, args{handler}, rWithToken, http.StatusOK, "user", true},
		{"bad request", fields{[]byte(secret)}, args{handler}, r, http.StatusUnauthorized, "", false},
		{"bad secret", fields{[]byte(wrongSecret)}, args{handler}, rWithToken, http.StatusUnauthorized, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tm := TestMiddleWare{t, tt.expectedUser, tt.expextedUserOK}
			h := tm.Handle(tt.args.next)

			m := AuthenticateMiddleware{
				Secret: tt.fields.Secret,
			}
			got := m.Handle(h)

			w := httptest.NewRecorder()
			got.ServeHTTP(w, tt.request)

			assert.Equal(t, w.Result().StatusCode, tt.expectedStatusCode)

		})
	}
}

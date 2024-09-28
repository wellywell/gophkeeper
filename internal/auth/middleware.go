package auth

import (
	"context"
	"net/http"
)

// AuthenticateMiddleware миддлвара для аутентификации пользователя
type AuthenticateMiddleware struct {
	Secret []byte
}

// UserKey ключ для хранения в контексте запроса имени пользователя
type UserKey string

const contextKey UserKey = "username"

// Handle хедер, проверяющий токен юзера и сохраняющий в контексте запроса юзернейм
func (m AuthenticateMiddleware) Handle(next http.Handler) http.Handler {

	authenticate := func(w http.ResponseWriter, r *http.Request) {

		user, err := VerifyUser(r, m.Secret)
		if err != nil {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), contextKey, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

	}
	return http.HandlerFunc(authenticate)
}

// GetAuthenticatedUser получает имя пользователя из контекста запроса
func GetAuthenticatedUser(req *http.Request) (string, bool) {
	user, ok := req.Context().Value(contextKey).(string)
	return user, ok

}

package auth

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

// Claims тип для работы с JWT
type Claims struct {
	jwt.RegisteredClaims
	Username string
}

// AuthHeader название заголовка для передачи токена
const AuthHeader = "X-Auth-Token"

// VerifyUser проверяет корректность токена и получает имя пользователя из него
func VerifyUser(r *http.Request, secret []byte) (string, error) {
	token := r.Header.Get(AuthHeader)
	if token != "" {
		user, err := GetUser(token, secret)
		if err != nil {
			return user, err
		}
		return user, nil
	}
	return "", fmt.Errorf("no auth token")
}

// BuildJWTString создаёт JWT-токен, содержащий имя пользователя
func BuildJWTString(user string, secret []byte) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{},

		Username: user,
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetUser парсит токен, проверяет его подпись и извлекает имя пользователя
func GetUser(tokenString string, secret []byte) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return secret, nil
		})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("token invalid")
	}

	return claims.Username, nil
}

// SetToken устанавливает заголовк авторизации с токеном
func SetToken(username string, w http.ResponseWriter, secret []byte) error {

	token, err := BuildJWTString(username, secret)
	if err != nil {
		return err
	}
	w.Header().Set(AuthHeader, token)
	return nil
}

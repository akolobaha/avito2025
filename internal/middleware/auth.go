package middleware

import (
	"avito2015/internal/token"
	"context"
	"net/http"
	"strings"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлечение токена из заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Проверка, начинается ли заголовок с "Bearer "
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		tokenModel, err := token.GetByToken(parts[1])
		if err != nil {
			http.Error(w, "Token does not exists", http.StatusUnauthorized)
			return
		}

		_, err = token.ValidateToken(tokenModel.Jwt)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
		}

		ctx := context.WithValue(r.Context(), "token", tokenModel)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}

// TODO: Bind user model to context
//func BindUserToContext(next http.HandlerFunc) http.HandlerFunc {
//
//}

package middleware

import (
	"avito2015/internal/token"
	"avito2015/internal/user"
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

		uRepo := user.NewUserRepository()
		uService := user.NewUserService(uRepo)
		tokenModel, usrModel, err := uService.GetUserAndTokenByJwt(parts[1])
		if err != nil {
			http.Error(w, "Token does not exists", http.StatusUnauthorized)
			return
		}

		_, err = token.ValidateToken(tokenModel.Jwt)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "token", tokenModel)
		ctx = context.WithValue(ctx, "user", usrModel)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}

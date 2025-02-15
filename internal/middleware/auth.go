package middleware

import (
	"avito2015/internal/token"
	"avito2015/internal/user"
	"avito2015/pkg/jsonresponse"
	"context"
	"errors"
	"net/http"
	"strings"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлечение токена из заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			jsonresponse.Error(w, errors.New("authorization header missing"), http.StatusUnauthorized)
			return
		}

		// Проверка, начинается ли заголовок с "Bearer "
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			jsonresponse.Error(w, errors.New("invalid authorization format"), http.StatusUnauthorized)
			return
		}

		uRepo := user.NewUserRepository()
		uService := user.NewUserService(uRepo)
		tokenModel, usrModel, err := uService.GetUserAndTokenByJwt(parts[1])
		if err != nil {
			jsonresponse.Error(w, errors.New("token does not exists"), http.StatusUnauthorized)
			return
		}

		s := token.NewService(nil)
		_, err = s.ValidateToken(tokenModel.Jwt)
		if err != nil {
			jsonresponse.Error(w, err, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "token", tokenModel)
		ctx = context.WithValue(ctx, "user", usrModel)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}

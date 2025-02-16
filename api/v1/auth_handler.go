package httphandler

import (
	"avito2015/internal/db"
	"avito2015/internal/user"
	"avito2015/pkg/jsonresponse"
	"encoding/json"
	"errors"
	"net/http"
)

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	var userReq user.AuthRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userReq); err != nil {
		http.Error(w, err.Error(), 400)
	}

	repo := user.NewUserRepository(db.DB)
	s := user.NewUserService(repo)

	token, err := s.CreateOrAuthUser(userReq.Username, userReq.Password)
	if err != nil {
		if errors.Is(err, user.InvalidPasswordError) {
			http.Error(w, err.Error(), 403)
			return
		}
		jsonresponse.Error(w, err, 400)
		return
	}

	jsonresponse.StatusOK(w, token)
}

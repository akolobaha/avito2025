package httphandler

import (
	"avito2015/internal/user"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func AuthHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("AuthHandler")

	var userReq user.AuthRequest

	//err := json.NewDecoder(r.Body).Decode(&userReq)
	//if err != nil {
	//	return
	//}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userReq); err != nil {
		http.Error(w, err.Error(), 400)
	}

	repo := user.NewUserRepository()
	s := user.NewUserService(repo)

	token, err := s.CreateOrAuthUser(userReq.Username, userReq.Password)
	if err != nil {
		if errors.Is(err, user.InvalidPasswordError) {
			http.Error(w, err.Error(), 403)
			return
		}
		fmt.Println(err.Error())
		http.Error(w, err.Error(), 400)
		return
	}

	renderJSON(w, token)
}

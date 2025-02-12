package httphandler

import (
	"avito2015/internal/user"
	"encoding/json"
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	repo := user.NewUserRepository()
	s := user.NewUserService(repo)

	err := s.CreateOrAuthUser(userReq.Username, userReq.Password)

	//_, err = w.Write([]byte(`{"message": "Authentication successful"}`))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	//w.WriteHeader()
	//fmt.Println("userReq:", userReq)
	//fmt.Println(s)
}

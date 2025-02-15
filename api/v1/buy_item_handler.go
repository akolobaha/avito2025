package httphandler

import (
	"avito2015/internal/merch"
	"avito2015/internal/user"
	"github.com/gorilla/mux"
	"net/http"
)

func BuyItemHandler(w http.ResponseWriter, r *http.Request) {
	usr, ok := r.Context().Value("user").(*user.User)
	if !ok {
		http.Error(w, "Unauthorized", 500)
		return
	}

	merchName := mux.Vars(r)["item"]

	mRepo := merch.NewMerchRepository()
	mService := merch.NewService(mRepo)

	err := mService.Buy(usr, merchName)
	if err != nil {
		jsonErrResp(w, err, http.StatusOK)
		return
	}

	jsonResp(w, MessageResp{
		"success",
	})
}

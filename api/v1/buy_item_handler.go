package httphandler

import (
	"avito2015/internal/db"
	"avito2015/internal/merch"
	"avito2015/internal/user"
	"avito2015/pkg/jsonresponse"
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

	mRepo := merch.NewMerchRepository(db.DB)
	mService := merch.NewService(mRepo)

	err := mService.Buy(usr, merchName)
	if err != nil {
		jsonresponse.Error(w, err, http.StatusOK)
		return
	}

	jsonresponse.StatusOK(w, jsonresponse.MessageResp{
		Message: "success",
	})
}

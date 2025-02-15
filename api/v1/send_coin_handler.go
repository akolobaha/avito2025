package httphandler

import (
	"avito2015/internal/transfer"
	"avito2015/internal/user"
	"encoding/json"
	"net/http"
)

func SendCoinHandler(w http.ResponseWriter, r *http.Request) {
	usr, ok := r.Context().Value("user").(*user.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	transferReq := transfer.CoinTransferReq{}
	if err := json.NewDecoder(r.Body).Decode(&transferReq); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	transferR := transfer.NewTransferRepository()
	transferS := transfer.NewTransferService(transferR)

	err := transferS.SendCoins(*usr, transferReq)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	resp := transfer.CoinTransferResp{}
	resp.Message = "Send successfully"

	renderJSON(w, resp)
}

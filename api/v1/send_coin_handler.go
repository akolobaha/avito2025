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
		http.Error(w, "Unauthorized", 500)
		return
	}

	transferReq := transfer.CoinTransferReq{}
	if err := json.NewDecoder(r.Body).Decode(&transferReq); err != nil {
		jsonErrResp(w, err, 400)
		return
	}

	transferR := transfer.NewTransferRepository()
	transferS := transfer.NewTransferService(transferR)

	err := transferS.SendCoins(*usr, transferReq)
	if err != nil {
		jsonErrResp(w, err, 500)
		return
	}

	resp := transfer.CoinTransferResp{}
	resp.Message = "Send successfully"

	jsonResp(w, resp)
}

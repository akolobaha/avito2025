package httphandler

import (
	"avito2015/internal/transfer"
	"avito2015/internal/user"
	"avito2015/pkg/jsonresponse"
	"encoding/json"
	"errors"
	"net/http"
)

func SendCoinHandler(w http.ResponseWriter, r *http.Request) {
	usr, ok := r.Context().Value("user").(*user.User)
	if !ok {
		jsonresponse.Error(w, errors.New("unauthorised"), 500)
		return
	}

	transferReq := transfer.CoinTransferReq{}
	if err := json.NewDecoder(r.Body).Decode(&transferReq); err != nil {
		jsonresponse.Error(w, err, 400)
		return
	}

	transferR := transfer.NewTransferRepository()
	transferS := transfer.NewTransferService(transferR)

	err := transferS.SendCoins(*usr, transferReq)
	if err != nil {
		jsonresponse.Error(w, err, 500)
		return
	}

	resp := transfer.CoinTransferResp{}
	resp.Message = "Send successfully"

	jsonresponse.StatusOK(w, resp)
}

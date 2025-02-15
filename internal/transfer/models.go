package transfer

type CoinTransfer struct {
	UserIdFrom string `db:"user_id_from"`
	UserIdTo   string `db:"user_id_to"`
	Coins      string `json:"coins"`
}

type CoinTransferReq struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type CoinTransferResp struct {
	Message string `json:"message"`
}

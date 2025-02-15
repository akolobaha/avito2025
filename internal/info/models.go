package info

type Resp struct {
	Coins       int             `json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistory     `json:"coinHistory"`
}

type InventoryItem struct {
	Type     string `json:"type" db:"name"`
	Quantity int    `json:"quantity" db:"count"`
}

type CoinHistory struct {
	Received []CoinReceived `json:"received"`
	Sent     []CoinSent     `json:"sent"`
}

type CoinReceived struct {
	FromUser string `json:"fromUser" db:"username"`
	Amount   int    `json:"amount" db:"coins"`
}

type CoinSent struct {
	ToUser string `json:"toUser" db:"username"`
	Amount int    `json:"amount" db:"coins"`
}

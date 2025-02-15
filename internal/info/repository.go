package info

import (
	"avito2015/internal/db"
	"avito2015/internal/user"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetMerchGroupedByItems(usr *user.User) ([]InventoryItem, error)
	GetTransactionsSent(user *user.User) ([]CoinSent, error)
	GetTransactionsReceive(user *user.User) ([]CoinReceived, error)
}

type infoRepositoryImpl struct {
	db *sqlx.DB
}

func NewInfoRepository() Repository {
	return &infoRepositoryImpl{db: db.DB}
}

func (r *infoRepositoryImpl) GetMerchGroupedByItems(usr *user.User) ([]InventoryItem, error) {
	query := `SELECT m.name, count(merch_id) FROM user_merch JOIN public.merch m on m.id = user_merch.merch_id WHERE user_merch.user_id = $1 GROUP BY merch_id, m.name`
	var inventoryList []InventoryItem
	err := r.db.Select(&inventoryList, query, usr.Id)
	if err != nil {
		return nil, err
	}

	return inventoryList, nil
}

func (r *infoRepositoryImpl) GetTransactionsSent(user *user.User) ([]CoinSent, error) {
	query := `SELECT u.username, coin_transfer.coins FROM coin_transfer JOIN "user" u on coin_transfer.user_id_to = u.id WHERE coin_transfer.user_id_from = $1;`
	var coinSentList []CoinSent
	err := r.db.Select(&coinSentList, query, user.Id)
	if err != nil {
		return nil, err
	}

	return coinSentList, nil
}

func (r *infoRepositoryImpl) GetTransactionsReceive(user *user.User) ([]CoinReceived, error) {
	query := `SELECT u.username, coin_transfer.coins FROM coin_transfer JOIN "user" u on coin_transfer.user_id_from = u.id WHERE coin_transfer.user_id_to = $1;`
	var coinReceiveList []CoinReceived
	err := r.db.Select(&coinReceiveList, query, user.Id)
	if err != nil {
		return nil, err
	}

	return coinReceiveList, nil
}

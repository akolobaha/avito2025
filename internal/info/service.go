package info

import (
	"avito2015/internal/user"
)

type Service struct {
	repo Repository
}

func NewInfoService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(user *user.User) (Resp, error) {
	// Покупки
	merchItems, err := s.repo.GetMerchGroupedByItems(user)
	if err != nil {
		return Resp{}, err
	}
	tSent, err := s.repo.GetTransactionsSent(user)
	tReceived, err := s.repo.GetTransactionsReceive(user)

	if err != nil {
		return Resp{}, err
	}
	history := CoinHistory{
		Received: tReceived,
		Sent:     tSent,
	}

	response := Resp{
		Coins:       user.Coins,
		Inventory:   merchItems,
		CoinHistory: history,
	}

	return response, nil
}

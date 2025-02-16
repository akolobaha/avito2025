package transfer

import (
	"avito2015/internal/db"
	"avito2015/internal/user"
	"database/sql"
	"errors"
)

var ErrorNotEnoughMoney = errors.New("not enough coins")

type Service struct {
	repo     Repository
	userRepo user.Repository
}

func NewTransferService(userRepo user.Repository, repo Repository) *Service {
	return &Service{userRepo: userRepo, repo: repo}
}

func (s *Service) SendCoins(from user.User, transferReq CoinTransferReq) error {
	// Проверим до начала, что денег достаточно
	if from.Coins < transferReq.Amount {
		return ErrorNotEnoughMoney
	}

	// Отдельным запросом получим пользователя из CoinTransferReq
	uRepo := user.NewUserRepository(db.DB)
	usrTo, err := uRepo.FindByUsername(transferReq.ToUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}

	// Транзакция отправки монет
	r := NewTransferRepository(db.DB)
	err = r.FromUserToUser(from, usrTo, transferReq.Amount)
	if err != nil {
		return err
	}

	return nil
}

package transfer

import (
	"avito2015/internal/db"
	"avito2015/internal/user"
	"github.com/jmoiron/sqlx"
	"log"
)

type Repository interface {
	FromUserToUser(from user.User, to user.User, amount int) error
}

type transferRepositoryImpl struct {
	db *sqlx.DB
}

func NewTransferRepository() Repository {
	return &transferRepositoryImpl{db: db.DB}
}

func (r *transferRepositoryImpl) FromUserToUser(from user.User, to user.User, amount int) error {
	// Начинаем транзакцию
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Вставка в журнал переводов
	_, err = tx.Exec("INSERT INTO coin_transfer(user_id_from, user_id_to, coins) VALUES ($1, $2, $3)", from.Id, to.Id, amount)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Обновляем количество монет у отправителя
	_, err = tx.Exec("UPDATE \"user\" SET coins = coins - $1 WHERE id = $2", amount, from.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Обновляем количество монет у получателя
	_, err = tx.Exec("UPDATE \"user\" SET coins = coins + $1 WHERE id = $2", amount, to.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Коммитим транзакцию
	if err = tx.Commit(); err != nil {
		return err
	}

	log.Println("Транзакция успешно выполнена!")

	return nil
}

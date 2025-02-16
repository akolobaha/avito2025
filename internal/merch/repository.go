package merch

import (
	"avito2015/internal/user"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	FindMerchByName(name string) (*Merch, error)
	Buy(usr *user.User, merch *Merch) error
}

type merchRepositoryImpl struct {
	db *sqlx.DB
}

func (m merchRepositoryImpl) FindMerchByName(name string) (*Merch, error) {
	var merch Merch
	query := `SELECT id, name, price FROM merch WHERE name = $1`
	err := m.db.Get(&merch, query, name)
	if err != nil {
		return &Merch{}, err
	}
	return &merch, nil
}

func (m merchRepositoryImpl) Buy(usr *user.User, merch *Merch) error {
	// Начинаем транзакцию

	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	// Обновляем количество монет у покупателя
	_, err = tx.Exec("UPDATE \"user\" SET coins = coins - $1 WHERE id = $2", merch.Price, usr.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Делаем запись в журнал мерча
	_, err = tx.Exec(`INSERT INTO user_merch(user_id, merch_id) VALUES ($1, $2)`, usr.Id, merch.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Коммитим транзакцию
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func NewMerchRepository(db *sqlx.DB) Repository {
	return &merchRepositoryImpl{db: db}
}

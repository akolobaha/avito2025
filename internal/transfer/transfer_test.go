package transfer

import (
	"avito2015/internal/db"
	"avito2015/internal/user"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromUserToUser_Success(t *testing.T) {
	// Create a mock DB connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock DB: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")
	repo := NewTransferRepository(sqlxDB)

	// Mock the expected SQL queries and responses for a successful transfer
	mock.ExpectBegin() // Start a transaction

	// Mock the first expected SQL query (INSERT into coin_transfer)
	mock.ExpectExec(`^INSERT INTO coin_transfer\(user_id_from, user_id_to, coins\) VALUES \(\$1, \$2, \$3\)$`).
		WithArgs(1, 2, 100).
		WillReturnResult(sqlmock.NewResult(1, 1)) // Insert into coin_transfer

	// Mock the second expected SQL query (UPDATE sender)
	mock.ExpectExec(`^UPDATE "user" SET coins = coins - \$1 WHERE id = \$2$`).
		WithArgs(100, 1).
		WillReturnResult(sqlmock.NewResult(1, 1)) // Update sender

	// Mock the third expected SQL query (UPDATE receiver)
	mock.ExpectExec(`^UPDATE "user" SET coins = coins \+ \$1 WHERE id = \$2$`).
		WithArgs(100, 2).
		WillReturnResult(sqlmock.NewResult(1, 1)) // Update receiver

	// Mock the commit of the transaction
	mock.ExpectCommit() // Commit the transaction

	// Define the users
	fromUser := user.User{Id: 1, Username: "user1", Coins: 500}
	toUser := user.User{Id: 2, Username: "user2", Coins: 100}

	// Call the function
	err = repo.FromUserToUser(fromUser, toUser, 100)

	// Assertions
	assert.NoError(t, err)                        // No error should occur
	assert.NoError(t, mock.ExpectationsWereMet()) // Ensure all expectations were met
}

func TestSendCoins_InsufficientFunds(t *testing.T) {
	// Create a mock DB connection
	dbMock, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock DB: %v", err)
	}
	defer dbMock.Close()

	// Mock the repositories
	sqlxDB := sqlx.NewDb(dbMock, "pgx")
	userRepo := user.NewUserRepository(sqlxDB)
	repo := NewTransferRepository(db.DB)
	service := NewTransferService(userRepo, repo)

	// Define the user and transfer request (Insufficient coins scenario)
	fromUser := user.User{Id: 1, Username: "sender", Coins: 50}
	transferReq := CoinTransferReq{ToUser: "recipient", Amount: 100}

	// Call the SendCoins method
	err = service.SendCoins(fromUser, transferReq)

	// Assertions
	assert.Error(t, err)                      // Error should occur due to insufficient funds
	assert.Equal(t, ErrorNotEnoughMoney, err) // The error should be "not enough coins"
}

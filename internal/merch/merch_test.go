package merch

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"

	"avito2015/internal/user"
)

// Mock Repository
type mockRepository struct {
	findMerchError error
	merchItem      *Merch // Указатель на структуру Merch
	buyError       error
}

func (m *mockRepository) FindMerchByName(name string) (*Merch, error) {
	if m.findMerchError != nil {
		return nil, m.findMerchError
	}
	return m.merchItem, nil
}

func (m *mockRepository) Buy(usr *user.User, merch *Merch) error {
	return m.buyError
}

func TestBuy(t *testing.T) {
	tests := []struct {
		name          string
		mockRepo      mockRepository
		userCoins     int
		merchName     string
		expectedError error
	}{
		{
			name: "Success",
			mockRepo: mockRepository{
				merchItem: &Merch{ID: 1, Name: "item1", Price: 50},
			},
			userCoins:     100,
			merchName:     "item1",
			expectedError: nil,
		},
		{
			name: "Merch not found",
			mockRepo: mockRepository{
				findMerchError: errors.New("merch not found"),
			},
			userCoins:     100,
			merchName:     "item2",
			expectedError: errors.New("merch not found"),
		},
		{
			name: "Not enough coins",
			mockRepo: mockRepository{
				merchItem: &Merch{ID: 1, Name: "item1", Price: 150}, // Цена больше, чем у пользователя
			},
			userCoins:     100,
			merchName:     "item1",
			expectedError: errors.New("not enough coins"),
		},
		{
			name: "Transaction error",
			mockRepo: mockRepository{
				merchItem: &Merch{ID: 1, Name: "item1", Price: 50},
				buyError:  errors.New("transaction failed"),
			},
			userCoins:     100,
			merchName:     "item1",
			expectedError: errors.New("transaction failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &tt.mockRepo
			service := NewService(repo)
			user := &user.User{Id: 1, Coins: tt.userCoins}

			err := service.Buy(user, tt.merchName)

			if err != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestFindMerchByName(t *testing.T) {
	// Create a mock DB connection
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock DB: %v", err)
	}
	defer dbMock.Close()

	// Create the merch repository with the mock DB
	sqlxDB := sqlx.NewDb(dbMock, "pgx")
	repo := NewMerchRepository(sqlxDB)

	// Define the expected SQL query and mock the response
	mock.ExpectQuery(`SELECT id, name, price FROM merch WHERE name = \$1`).
		WithArgs("T-shirt").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).
			AddRow(1, "T-shirt", 100)) // Mock the result for the merch item

	// Call the FindMerchByName method
	merch, err := repo.FindMerchByName("T-shirt")

	// Assertions
	assert.NoError(t, err)                        // No error should occur
	assert.NotNil(t, merch)                       // Merch should not be nil
	assert.Equal(t, "T-shirt", merch.Name)        // The name should match
	assert.Equal(t, "1", merch.ID)                // The ID should match
	assert.Equal(t, 100, merch.Price)             // The price should match
	assert.NoError(t, mock.ExpectationsWereMet()) // Ensure all expectations were met
}

func TestBuy_Success(t *testing.T) {
	// Create a mock DB connection
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock DB: %v", err)
	}
	defer dbMock.Close()

	// Create the merch repository with the mock DB
	sqlxDB := sqlx.NewDb(dbMock, "pgx")
	repo := NewMerchRepository(sqlxDB)

	// Mock the expected SQL queries for a successful purchase
	mock.ExpectBegin() // Start a transaction
	mock.ExpectExec(`UPDATE "user" SET coins = coins - \$1 WHERE id = \$2`).
		WithArgs(100, 1).                         // Mock the update for user coins (pass 1 as int)
		WillReturnResult(sqlmock.NewResult(1, 1)) // Mock successful update

	mock.ExpectExec(`INSERT INTO user_merch\(user_id, merch_id\) VALUES \(\$1, \$2\)`).
		WithArgs(1, 1).                           // Mock the insert into user_merch (user_id: 1, merch_id: 1)
		WillReturnResult(sqlmock.NewResult(1, 1)) // Mock successful insert

	mock.ExpectCommit() // Commit the transaction

	// Define the user and merch item
	usr := &user.User{Id: 1, Username: "user1", Coins: 200} // User ID is an integer
	merch := &Merch{ID: 1, Name: "T-shirt", Price: 100}     // Merch ID is an integer

	// Call the Buy method
	err = repo.Buy(usr, merch)

	// Assertions
	assert.NoError(t, err)                        // No error should occur
	assert.NoError(t, mock.ExpectationsWereMet()) // Ensure all expectations were met
}

func TestBuy_RollbackOnError(t *testing.T) {
	// Create a mock DB connection
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock DB: %v", err)
	}
	defer dbMock.Close()

	// Create the merch repository with the mock DB
	sqlxDB := sqlx.NewDb(dbMock, "pgx")
	repo := NewMerchRepository(sqlxDB)

	// Mock the expected SQL queries and simulate an error during the transaction
	mock.ExpectBegin() // Start a transaction
	mock.ExpectExec(`UPDATE "user" SET coins = coins - \$1 WHERE id = \$2`).
		WithArgs(100, 1).                         // Mock the update for user coins
		WillReturnResult(sqlmock.NewResult(1, 1)) // Mock successful update

	// Simulate an error during the insert into user_merch
	mock.ExpectExec(`INSERT INTO user_merch\(user_id, merch_id\) VALUES \(\$1, \$2\)`).
		WithArgs(1, 1). // Mock the insert into user_merch
		WillReturnError(fmt.Errorf("mock error during insert"))

	mock.ExpectRollback() // Rollback the transaction

	// Define the user and merch item
	usr := &user.User{Id: 1, Username: "user1", Coins: 200}
	merch := &Merch{ID: 1, Name: "T-shirt", Price: 100}

	// Call the Buy method
	err = repo.Buy(usr, merch)

	// Assertions
	assert.Error(t, err)                                     // An error should occur
	assert.Equal(t, "mock error during insert", err.Error()) // The error message should match the mock error
	assert.NoError(t, mock.ExpectationsWereMet())            // Ensure all expectations were met
}

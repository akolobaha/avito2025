package info_test

import (
	"errors"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"

	"avito2015/internal/info"
	"avito2015/internal/user"
)

// Mock Repository
type mockRepository struct {
	merchItems       []info.InventoryItem
	tSent            []info.CoinSent
	tReceived        []info.CoinReceived
	getMerchError    error
	getSentError     error
	getReceivedError error
}

func (m *mockRepository) GetMerchGroupedByItems(u *user.User) ([]info.InventoryItem, error) {
	return m.merchItems, m.getMerchError
}

func (m *mockRepository) GetTransactionsSent(u *user.User) ([]info.CoinSent, error) {
	return m.tSent, m.getSentError
}

func (m *mockRepository) GetTransactionsReceive(u *user.User) ([]info.CoinReceived, error) {
	return m.tReceived, m.getReceivedError
}

// Функция для сравнения двух объектов Resp
func respEqual(a, b info.Resp) bool {
	if a.Coins != b.Coins {
		return false
	}
	if len(a.Inventory) != len(b.Inventory) {
		return false
	}
	for i := range a.Inventory {
		if a.Inventory[i] != b.Inventory[i] {
			return false
		}
	}
	if len(a.CoinHistory.Received) != len(b.CoinHistory.Received) {
		return false
	}
	for i := range a.CoinHistory.Received {
		if a.CoinHistory.Received[i] != b.CoinHistory.Received[i] {
			return false
		}
	}
	if len(a.CoinHistory.Sent) != len(b.CoinHistory.Sent) {
		return false
	}
	for i := range a.CoinHistory.Sent {
		if a.CoinHistory.Sent[i] != b.CoinHistory.Sent[i] {
			return false
		}
	}
	return true
}

func TestGet(t *testing.T) {
	tests := []struct {
		name          string
		mockRepo      mockRepository
		expectedResp  info.Resp
		expectedError error
	}{
		{
			name: "Success",
			mockRepo: mockRepository{
				merchItems: []info.InventoryItem{
					{Type: "item1", Quantity: 5},
					{Type: "item2", Quantity: 3},
				},
				tSent: []info.CoinSent{
					{ToUser: "user1", Amount: 10},
				},
				tReceived: []info.CoinReceived{
					{FromUser: "user2", Amount: 20},
				},
			},
			expectedResp: info.Resp{
				Coins: 100,
				Inventory: []info.InventoryItem{
					{Type: "item1", Quantity: 5},
					{Type: "item2", Quantity: 3},
				},
				CoinHistory: info.CoinHistory{
					Received: []info.CoinReceived{
						{FromUser: "user2", Amount: 20},
					},
					Sent: []info.CoinSent{
						{ToUser: "user1", Amount: 10},
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "GetMerchGroupedByItems error",
			mockRepo: mockRepository{
				getMerchError: errors.New("error fetching merch"),
			},
			expectedResp:  info.Resp{},
			expectedError: errors.New("error fetching merch"),
		},
		{
			name: "GetTransactionsSent error",
			mockRepo: mockRepository{
				merchItems: []info.InventoryItem{
					{Type: "item1", Quantity: 5},
				},
				getSentError: errors.New("error fetching sent transactions"),
			},
			expectedResp:  info.Resp{},
			expectedError: errors.New("error fetching sent transactions"),
		},
		{
			name: "GetTransactionsReceive error",
			mockRepo: mockRepository{
				merchItems: []info.InventoryItem{
					{Type: "item1", Quantity: 5},
				},
				tSent: []info.CoinSent{
					{ToUser: "user1", Amount: 10},
				},
				getReceivedError: errors.New("error fetching received transactions"),
			},
			expectedResp:  info.Resp{},
			expectedError: errors.New("error fetching received transactions"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &tt.mockRepo
			service := info.NewInfoService(repo)
			user := &user.User{Coins: 100}

			resp, err := service.Get(user)

			if err != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}

			if !respEqual(resp, tt.expectedResp) {
				t.Errorf("expected response %v, got %v", tt.expectedResp, resp)
			}
		})
	}
}

func TestGetMerchGroupedByItems(t *testing.T) {
	// Create a mock DB connection
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock DB: %v", err)
	}
	defer dbMock.Close()

	// Create the info repository with the mock DB
	sqlxDB := sqlx.NewDb(dbMock, "pgx")
	repo := info.NewInfoRepository(sqlxDB)

	// Define the expected SQL query and mock the response
	mock.ExpectQuery(`SELECT m.name, count\(merch_id\) FROM user_merch JOIN public.merch m on m.id = user_merch.merch_id WHERE user_merch.user_id = \$1 GROUP BY merch_id, m.name`).
		WithArgs(1). // Mock user_id = 1
		WillReturnRows(sqlmock.NewRows([]string{"name", "count"}).
			AddRow("T-shirt", 2). // Mock item "T-shirt" with 2 counts
			AddRow("Sweater", 1)) // Mock item "Sweater" with 1 count

	// Define a user
	usr := &user.User{Id: 1}

	// Call the GetMerchGroupedByItems method
	items, err := repo.GetMerchGroupedByItems(usr)

	// Assertions
	assert.NoError(t, err)  // No error should occur
	assert.Len(t, items, 2) // Two items should be returned
	//assert.Equal(t, "T-shirt", items[0].Name)     // The first item should be "T-shirt"
	//assert.Equal(t, 2, items[0].Count)            // The count for "T-shirt" should be 2
	//assert.Equal(t, "Sweater", items[1].Name)     // The second item should be "Sweater"
	//assert.Equal(t, 1, items[1].Count)            // The count for "Sweater" should be 1
	assert.NoError(t, mock.ExpectationsWereMet()) // Ensure all expectations were met
}

func TestGetTransactionsSent(t *testing.T) {
	// Create a mock DB connection
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock DB: %v", err)
	}
	defer dbMock.Close()

	// Create the info repository with the mock DB
	sqlxDB := sqlx.NewDb(dbMock, "pgx")
	repo := info.NewInfoRepository(sqlxDB)

	// Define the expected SQL query and mock the response
	mock.ExpectQuery(`SELECT u.username, coin_transfer.coins FROM coin_transfer JOIN "user" u on coin_transfer.user_id_to = u.id WHERE coin_transfer.user_id_from = \$1`).
		WithArgs(1). // Mock user_id = 1
		WillReturnRows(sqlmock.NewRows([]string{"username", "coins"}).
			AddRow("user2", 50). // Mock transaction to "user2" with 50 coins
			AddRow("user3", 20)) // Mock transaction to "user3" with 20 coins

	// Define a user
	usr := &user.User{Id: 1}

	// Call the GetTransactionsSent method
	transactions, err := repo.GetTransactionsSent(usr)

	// Assertions
	assert.NoError(t, err)         // No error should occur
	assert.Len(t, transactions, 2) // Two transactions should be returned
	//assert.Equal(t, "user2", transactions[0].Username) // The first transaction should be to "user2"
	//assert.Equal(t, 50, transactions[0].Coins)         // The coin amount should be 50
	//assert.Equal(t, "user3", transactions[1].Username) // The second transaction should be to "user3"
	//assert.Equal(t, 20, transactions[1].Coins)         // The coin amount should be 20
	assert.NoError(t, mock.ExpectationsWereMet()) // Ensure all expectations were met
}

func TestGetTransactionsReceive(t *testing.T) {
	// Create a mock DB connection
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock DB: %v", err)
	}
	defer dbMock.Close()

	// Create the info repository with the mock DB
	sqlxDB := sqlx.NewDb(dbMock, "pgx")
	repo := info.NewInfoRepository(sqlxDB)

	// Define the expected SQL query and mock the response
	mock.ExpectQuery(`SELECT u.username, coin_transfer.coins FROM coin_transfer JOIN "user" u on coin_transfer.user_id_from = u.id WHERE coin_transfer.user_id_to = \$1`).
		WithArgs(1). // Mock user_id = 1
		WillReturnRows(sqlmock.NewRows([]string{"username", "coins"}).
			AddRow("user4", 30). // Mock transaction from "user4" with 30 coins
			AddRow("user5", 60)) // Mock transaction from "user5" with 60 coins

	// Define a user
	usr := &user.User{Id: 1}

	// Call the GetTransactionsReceive method
	transactions, err := repo.GetTransactionsReceive(usr)

	// Assertions
	assert.NoError(t, err)         // No error should occur
	assert.Len(t, transactions, 2) // Two transactions should be returned
	//assert.Equal(t, "user4", transactions[0].Username) // The first transaction should be from "user4"
	//assert.Equal(t, 30, transactions[0].Coins)         // The coin amount should be 30
	//assert.Equal(t, "user5", transactions[1].Username) // The second transaction should be from "user5"
	//assert.Equal(t, 60, transactions[1].Coins)         // The coin amount should be 60
	assert.NoError(t, mock.ExpectationsWereMet()) // Ensure all expectations were met
}

package info_test

import (
	"errors"
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

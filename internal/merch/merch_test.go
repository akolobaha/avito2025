package merch

import (
	"errors"
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
				merchItem: &Merch{ID: "1", Name: "item1", Price: 50},
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
				merchItem: &Merch{ID: "1", Name: "item1", Price: 150}, // Цена больше, чем у пользователя
			},
			userCoins:     100,
			merchName:     "item1",
			expectedError: errors.New("not enough coins"),
		},
		{
			name: "Transaction error",
			mockRepo: mockRepository{
				merchItem: &Merch{ID: "1", Name: "item1", Price: 50},
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

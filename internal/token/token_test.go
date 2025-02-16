package token_test

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"testing"

	"avito2015/internal/token"
)

// Mock Repository
type mockRepository struct {
	saveError   error
	getToken    *token.Token
	getError    error
	updateError error
}

func (m *mockRepository) Save(t token.Token) error {
	return m.saveError
}

func (m *mockRepository) Get(userID int) (*token.Token, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	return m.getToken, nil
}

func (m *mockRepository) Update(t token.Token) (*token.Token, error) {
	if m.updateError != nil {
		return nil, m.updateError
	}
	return &t, nil
}

func TestSaveToken(t *testing.T) {
	os.Setenv("TOKEN_SECRET", "testsecret")

	tests := []struct {
		name          string
		mockRepo      mockRepository
		userId        int64
		expectedError error
	}{
		{
			name: "Success",
			mockRepo: mockRepository{
				getToken: &token.Token{
					Jwt:      "valid.jwt.token",
					UserId:   1,
					IsActive: true,
				},
			},
			userId:        1,
			expectedError: nil,
		},
		{
			name: "Token generation error",
			mockRepo: mockRepository{
				saveError: errors.New("error saving token"),
			},
			userId:        1,
			expectedError: errors.New("error saving token"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &tt.mockRepo
			service := token.NewService(repo)

			_, err := service.SaveToken(tt.userId)

			if err != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestGetByUserId(t *testing.T) {
	tests := []struct {
		name          string
		mockRepo      mockRepository
		userId        int
		expectedToken *token.Token
		expectedError error
	}{
		{
			name: "Token not found",
			mockRepo: mockRepository{
				getError: errors.New("token not found"),
			},
			userId:        1,
			expectedToken: nil,
			expectedError: errors.New("token not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &tt.mockRepo
			service := token.NewService(repo)

			token, err := service.GetByUserId(tt.userId)

			if err != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}

			if token != nil && *token != *tt.expectedToken {
				t.Errorf("expected token %v, got %v", tt.expectedToken, token)
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	os.Setenv("TOKEN_SECRET", "testsecret")

	tests := []struct {
		name           string
		tokenString    string
		expectedClaims *token.Claims
		expectedError  error
	}{
		{
			name: "Valid token",
			tokenString: func() string {
				claims := &token.Claims{
					RegisteredClaims: jwt.RegisteredClaims{
						Issuer: "1",
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				signedToken, _ := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
				return signedToken
			}(),
			expectedClaims: &token.Claims{RegisteredClaims: jwt.RegisteredClaims{Issuer: "1"}},
			expectedError:  nil,
		},
		{
			name:           "Invalid token",
			tokenString:    "invalid.token.string",
			expectedClaims: nil,
			expectedError:  errors.New("invalid character '\\u008a' looking for beginning of value"),
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			service := token.NewService(&mockRepository{})

			claims, err := service.ValidateToken(tt.tokenString)

			if err != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}

			if claims != nil && !claimsEqual(claims, tt.expectedClaims) {
				t.Errorf("expected claims %v, got %v", tt.expectedClaims, claims)
			}
		})
	}
}

// Функция для сравнения Claims
func claimsEqual(a, b *token.Claims) bool {
	if a.Issuer != b.Issuer {
		return false
	}
	if a.IssuedAt != nil && b.IssuedAt != nil {
		if !a.IssuedAt.Time.Equal(b.IssuedAt.Time) {
			return false
		}
	} else if a.IssuedAt != nil || b.IssuedAt != nil {
		return false
	}
	if a.ExpiresAt != nil && b.ExpiresAt != nil {
		if !a.ExpiresAt.Time.Equal(b.ExpiresAt.Time) {
			return false
		}
	} else if a.ExpiresAt != nil || b.ExpiresAt != nil {
		return false
	}
	return true
}

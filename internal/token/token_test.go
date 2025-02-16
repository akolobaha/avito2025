package token_test

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
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

func TestGet(t *testing.T) {
	// Create a mock DB connection
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock DB: %v", err)
	}
	defer dbMock.Close()

	// Initialize the repository with the mock DB
	sqlxDB := sqlx.NewDb(dbMock, "pgx")
	repo := token.NewTokenRepository(sqlxDB)

	// Define the expected SQL query and mock the response
	mock.ExpectQuery(`SELECT jwt, user_id, is_active FROM "user_token" WHERE user_id = \$1 AND is_active = true`).
		WithArgs(1). // Pass the user ID
		WillReturnRows(sqlmock.NewRows([]string{"jwt", "user_id", "is_active"}).
			AddRow("mock-jwt-token", 1, true)) // Mock token data

	// Call the Get method
	token, err := repo.Get(1)

	// Assertions
	assert.NoError(t, err)                        // No error should occur
	assert.NotNil(t, token)                       // The returned token should not be nil
	assert.Equal(t, "mock-jwt-token", token.Jwt)  // The token should match the mock token
	assert.Equal(t, int64(1), token.UserId)       // The user ID should match the mock user ID
	assert.NoError(t, mock.ExpectationsWereMet()) // Ensure all expectations were met
}

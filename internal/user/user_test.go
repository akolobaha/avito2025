package user

import (
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestSave(t *testing.T) {
	// Create a mock DB connection
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")

	// Create the repository
	repo := NewUserRepository(sqlxDB)

	// Define the expected SQL query and mock the response
	mock.ExpectQuery(`INSERT INTO "user" \(.+\) VALUES \(.+\) RETURNING id`).
		WithArgs("testuser", "testpassword").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Call the Save function
	lastInsertID, err := repo.Save("testuser", "testpassword")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, int64(1), lastInsertID)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByUsername(t *testing.T) {
	// Create a mock DB connection
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")

	// Create the repository
	repo := NewUserRepository(sqlxDB)

	// Define the expected SQL query and mock the response
	mock.ExpectQuery(`^SELECT id, username, coins, password, created_at FROM "user" WHERE username = \$1$`).
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "coins", "password", "created_at"}).
			AddRow(1, "testuser", 100, "testpassword", time.Now()))

	// Call the FindByUsername function
	user, err := repo.FindByUsername("testuser")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, 1, user.Id)
	assert.Equal(t, 100, user.Coins)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByToken(t *testing.T) {
	// Create a mock DB connection
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")

	// Create the repository
	repo := NewUserRepository(sqlxDB)

	// Define the expected SQL query and mock the response
	mock.ExpectQuery(`^SELECT ut.jwt, ut.user_id, ut.is_active, u.id, u.username, u.password, u.coins, u.created_at FROM user_token ut JOIN public."user" u ON u.id = ut.user_id WHERE ut.jwt = \$1 AND ut.is_active = true$`).
		WithArgs("testtoken").
		WillReturnRows(sqlmock.NewRows([]string{
			"jwt", "user_id", "is_active", "id", "username", "password", "coins", "created_at"}).
			AddRow("testtoken", 1, true, 1, "testuser", "testpassword", 100, time.Now()))

	// Call the GetByToken function
	tokenM, usr, err := repo.GetByToken("testtoken")

	// Ensure that no error occurred
	assert.NoError(t, err)

	// Ensure tokenM and usr are not nil
	assert.NotNil(t, tokenM)
	assert.NotNil(t, usr)

	// Assertions for tokenM
	assert.Equal(t, "testtoken", tokenM.Jwt)
	assert.Equal(t, int64(1), tokenM.UserId)
	assert.Equal(t, true, tokenM.IsActive)

	// Assertions for usr
	assert.Equal(t, "testuser", usr.Username)
	assert.Equal(t, 1, usr.Id)
	assert.Equal(t, 100, usr.Coins)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

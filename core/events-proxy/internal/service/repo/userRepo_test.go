package repo

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func setupUserRepo(t *testing.T) (userRepo, sqlmock.Sqlmock) {
	db, m, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	r := NewUserRepo(sqlxDB)

	return r, m
}

func Test_userRepo_GetByUsername(t *testing.T) {
	repo, mDB := setupUserRepo(t)

	username := "username1"
	aesKey := "aesKey"
	totpKey := "totpKey"

	expectedRows := mDB.NewRows([]string{"username", "aes_key", "totp_key"}).AddRow(username, aesKey, totpKey)
	mDB.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(expectedRows)

	record, err := repo.GetByUsername(context.Background(), username)
	require.NoError(t, err)
	require.EqualValues(t, username, record.Username)
	require.EqualValues(t, aesKey, record.AesKey)
	require.EqualValues(t, totpKey, record.TotpKey)
}

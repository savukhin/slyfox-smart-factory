package repo

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var statementBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type UserRecord struct {
	Username string `db:"username"`
	TotpKey  string `db:"totp_key"`
	AesKey   string `db:"aes_key"`
}

type UserRepo interface {
	GetByUsername(ctx context.Context, username string) (UserRecord, error)
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) userRepo {
	return userRepo{
		db: db,
	}
}

func (repo *userRepo) GetByUsername(ctx context.Context, username string) (record UserRecord, err error) {
	sql, args, err := statementBuilder.
		Select("username", "aes_key", "totp_key").
		From("users").
		ToSql()

	if err != nil {
		return
	}

	row := repo.db.QueryRowxContext(ctx, sql, args)
	if row.Err() != nil {
		err = row.Err()
		return
	}

	err = row.StructScan(&record)
	return
}

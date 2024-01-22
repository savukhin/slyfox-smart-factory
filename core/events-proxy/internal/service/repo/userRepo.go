package repo

import (
	"context"
	"eventsproxy/internal/domain"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var statementBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type UserRepo interface {
	GetByUsername(ctx context.Context, username string) (domain.User, error)
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) userRepo {
	return userRepo{
		db: db,
	}
}

func (repo *userRepo) GetByUsername(ctx context.Context, username string) (record domain.User, err error) {
	sql, args, err := statementBuilder.
		Select("username", "aes_key", "totp_key").
		From("users").
		Where(squirrel.Eq{"username": username}).
		ToSql()

	if err != nil {
		return
	}

	row := repo.db.QueryRowxContext(ctx, sql, args...)
	if row.Err() != nil {
		err = row.Err()
		fmt.Println("err", err)
		return
	}

	err = row.StructScan(&record)
	return
}

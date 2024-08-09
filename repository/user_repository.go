package repository

import (
	"context"
	"database/sql"

	"github.com/gunawan98/golang-restfull-api/model/domain"
)

type UserRepository interface {
	Save(ctx context.Context, tx *sql.Tx, user domain.User) domain.User
	Update(ctx context.Context, tx *sql.Tx, user domain.User) domain.User
	Delete(ctx context.Context, tx *sql.Tx, userId int)
	FindById(ctx context.Context, tx *sql.Tx, userId int) (domain.User, error)
	FindAll(ctx context.Context, tx *sql.Tx) []domain.User
	FindByUsername(ctx context.Context, tx *sql.Tx, username string) (domain.User, error)
}

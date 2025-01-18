package postgresql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/muhrifqii/tuskar/domain"
	"github.com/muhrifqii/tuskar/internal/sqler"
	"go.uber.org/zap"
)

type (
	UserRepository struct {
		db *sqler.SqlxWrapper
	}
)

func NewUserRepository(db *sqlx.DB, zap *zap.Logger) *UserRepository {
	return &UserRepository{
		db: sqler.NewSqlxWrapper(db, zap),
	}
}

func (r *UserRepository) GetByUsername(c context.Context, username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE username = $1", username)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil // Or any default value for YourStruct
	}
	return &user, err
}

func (r *UserRepository) CreateUser(c context.Context, user *domain.User) error {

	_, err := r.db.NamedExec("INSERT INTO users (username, a_password, first_name, last_name) VALUES (:username, :a_password, :first_name, :last_name)", user)

	return err
}

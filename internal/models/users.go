package models

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// interface to help with testing later
type UserModelInterface interface {
	Insert(email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
}

type User struct {
	ID             int
	Email          string
	HashedPassword []byte
	CreatedAt      time.Time
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (um *UserModel) Insert(email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	statement := `INSERT INTO users (email, password, created_at)
    VALUES($1, $2, NOW())`

	_, err = um.DB.Exec(context.Background(), statement, email, string(hashedPassword))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && strings.Contains(pgErr.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		} else {
			return err
		}
	}
	return nil
}

func (um *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (um *UserModel) Exists(id int) (bool, error) {
	return true, nil
}

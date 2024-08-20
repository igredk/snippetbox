package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	query := `INSERT INTO users (name, email, hashed_password, created)
    VALUES($1, $2, $3, NOW()) RETURNING id`

	var id int
	err = m.DB.QueryRow(context.Background(), query, name, email, string(hashedPassword)).Scan(&id)
	if err != nil {
		var ucError *pgconn.PgError
		if errors.As(err, &ucError) {
			if ucError.Code == UniqueConstraintViolation {
				return ErrDuplicateEmail
			}
		}
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}

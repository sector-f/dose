package sqlite

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// Auth implements the dose.AuthService interface
type Auth struct {
	db *sql.DB
}

func (a *Auth) CheckAuthentication(username, password string) (bool, error) {
	var dbPass string

	row := a.db.QueryRow("SELECT password FROM users WHERE username = ?", username)

	err := row.Scan(&dbPass)
	switch err {
	case sql.ErrNoRows:
		return false, errors.New("Invalid user")
	case nil:
		break
	default:
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(password))
	if err != nil {
		return false, err
	}

	return true, nil
}

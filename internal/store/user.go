package store

import (
	"database/sql"
	"errors"
	"fmt"
)

type User struct {
	Id   int64
	Name string
}

type UserStore struct {
	db *sql.DB
}

// TODO: auth is still just dummy values
func (u *UserStore) GetUser(temp string) (*User, error) {
	user := User{}
	query := `SELECT u.id, u.name FROM users u WHERE u.name = ?`
	err := u.db.QueryRow(query, []interface{}{temp}...).Scan(&user.Id, &user.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
	}
	return &user, nil
}

package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var AnonymousUser = &User{}

type User struct {
	Id          int64    `json:"id"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Password    Password `json:"-"`
	CreatedAt   int64    `json:"created_at"`
	UpdatedAt   int64    `json:"updated_at"`
	Activated   bool     `json:"activated"`
	Permissions []string `json:"-"`
	Version     int      `json:"-"`
}

type UserStore struct {
	db *sql.DB
}
type Password struct {
	plaintext *string
	hash      string
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (p *Password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = string(hash)
	return nil
}
func (p *Password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func (u *UserStore) GetUser(email string) (*User, error) {
	user := User{Password: Password{}}
	query := `SELECT u.id, u.email, u.password_hash FROM users u WHERE u.email = ?`
	err := u.db.QueryRow(query, []any{email}...).Scan(&user.Id, &user.Name, &user.Password.hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
        return nil, fmt.Errorf("searching for user: %w", err)
	}
	return &user, nil
}

func (u *UserStore) GetForToken(tokenPlaintext string) (*User, error) {
	// Calculate the SHA-256 hash of the plaintext token provided by the client.
	// Remember that this returns a byte *array* with length 32, not a slice.
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))
	// Set up the SQL query.
	query := `
SELECT users.id, users.created_at, users.name, users.email, users.password_hash, users.activated, users.version
FROM users
INNER JOIN sessions
ON users.id = sessions.user_id
WHERE sessions.token_hash = ? 
AND sessions.expiry > ?`
	// Create a slice containing the query arguments. Notice how we use the [:] operator
	// to get a slice containing the token hash, rather than passing in the array (which
	// is not supported by the pq driver), and that we pass the current time as the
	// value to check against the token expiry.
	args := []interface{}{tokenHash[:], time.Now().UnixMilli()}
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Execute the query, scanning the return values into a User struct. If no matching
	// record is found we return an ErrRecordNotFound error.
	err := u.db.QueryRowContext(ctx, query, args...).Scan(
		&user.Id,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Return the matching user.
	return &user, nil
}

func (u *UserStore) Insert(user *User) error {
	args := []any{user.Name, user.Email, user.Password.hash, time.Now().UnixMilli(), time.Now().UnixMilli(), 0, 1}
	_, err := u.db.Exec("INSERT INTO users (name, email, password_hash, created_at, updated_at, activated, version) VALUES (?, ?, ?, ?, ?, ?, ?)", args...)
	return err
}

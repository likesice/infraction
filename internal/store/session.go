package store

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"
)

type SessionStore struct {
	db *sql.DB
}
type Session struct {
	Token     string `json:"token"`
	TokenHash []byte `json:"-"`
	User      int64  `json:"-"`
	Expiry    int64  `json:"expiry"`
	CreatedAt int64  `json:"-"`
	UpdatedAt int64  `json:"-"`
}

func (s *SessionStore) New(user *User) (*Session, error) {
	session := Session{
		User:      user.Id,
		Expiry:    time.Now().Add(time.Hour * 24).UnixMilli(),
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		//TODO: add context
		return nil, err
	}

	session.Token = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(session.Token))
	session.TokenHash = hash[:]
	s.Insert(&session)

	return &session, nil
}

func (s *SessionStore) Insert(session *Session) error {
	stmt := `INSERT INTO sessions (token_hash, user_id, expiry, created_at, updated_at) VALUES (?, ?, ?, ?, ?);`

	tx, err := s.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	args := []interface{}{session.TokenHash, session.User, session.Expiry, session.CreatedAt, session.UpdatedAt}
	exec, err := tx.Exec(stmt, args...)
	n, err := exec.RowsAffected()
	if err != nil || n != 1 {
		err := tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		err := tx.Rollback()
		return err
	}
	return nil
}

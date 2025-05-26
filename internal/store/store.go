package store

import (
	"database/sql"
	"errors"
	"time"
)

var ErrRecordNotFound = errors.New("record not found")

// used for testing purposes
var now = time.Now

type Store struct {
	GroupStore   GroupStore
	UserStore    UserStore
	SessionStore SessionStore
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		GroupStore:   GroupStore{db: db},
		UserStore:    UserStore{db: db},
		SessionStore: SessionStore{db: db},
	}
}

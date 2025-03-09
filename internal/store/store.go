package store

import (
	"database/sql"
	"time"
)

// used for testing purposes
var now = time.Now

type Store struct {
	GroupStore GroupStore
	UserStore  UserStore
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		GroupStore: GroupStore{db: db},
		UserStore:  UserStore{db: db},
	}
}

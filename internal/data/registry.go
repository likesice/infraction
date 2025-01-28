package data

import "database/sql"

type Registry struct {
	InfractionRepository InfractionRepository
	UserRepository       UserRepository
}

func NewRegistry(db *sql.DB) *Registry {
	return &Registry{
		InfractionRepository: InfractionRepository{db: db},
		UserRepository:       UserRepository{db: db},
	}
}

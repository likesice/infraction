package data

import "database/sql"

type Repositories struct {
	Infractions  InfractionRepository
	Transactions TransactionRepository
}

func NewRepositories(db *sql.DB) Repositories {
	return Repositories{
		Infractions:  InfractionRepository{DB: db},
		Transactions: TransactionRepository{DB: db},
	}
}

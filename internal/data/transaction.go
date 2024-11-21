package data

import (
	"database/sql"
	"infraction.mageis.net/internal/data/validator"
	ierrors "infraction.mageis.net/internal/errors"
	"time"
)

type Transaction struct {
	Id        int64     `json:"id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type TransactionRepository struct {
	DB *sql.DB
}

func (t *Transaction) Validate(v *validator.Validator) bool {
	v.Check(t.Amount > 0, "amount", "amount must be greater than zero")
	return v.Valid()
}

func (t *TransactionRepository) Insert(transaction *Transaction, infractionId int64) error {
	query := `insert into "transaction" (infraction_id, amount, created_at) 
              values (?, ?, ?)
              returning id, created_at`

	args := []interface{}{infractionId, transaction.Amount, time.Now().UnixMilli()}
	var createdAt int64
	err := t.DB.QueryRow(query, args...).Scan(&transaction.Id, &createdAt)
	if err != nil {
		return ierrors.ErrDbFailure.Wrap(err)
	}
	transaction.CreatedAt = time.UnixMilli(createdAt)
	return nil
}

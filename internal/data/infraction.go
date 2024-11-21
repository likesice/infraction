package data

import (
	"database/sql"
	"strings"
	"time"

	"infraction.mageis.net/internal/data/validator"
	ierrors "infraction.mageis.net/internal/errors"
)

type SplitKind int32

const (
	EvenSplit SplitKind = iota
)

type Infraction struct {
	Id           int64         `json:"id"`
	Name         string        `json:"name"`
	Split        SplitKind     `json:"split"`
	Transactions []Transaction `json:"transactions,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	Version      int32         `json:"version"`
}

type InfractionRepository struct {
	DB *sql.DB
}

func (i *Infraction) Validate(v *validator.Validator) bool {
	v.Check(len(i.Name) > 2, "name", "name must be longer than 2 bytes")
	return v.Valid()
}

func (i *InfractionRepository) Insert(infraction *Infraction) error {
	query := `insert into infraction (name, split, created_at, version) 
			  values  (?,?,?,?) 
              returning id, created_at, version`
	args := []interface{}{infraction.Name, infraction.Split, time.Now().UnixMilli(), 1}

	var createdAt int64
	tx, err := i.DB.Begin()
	if err != nil {
		return ierrors.ErrDbFailure.Wrap(err)
	}
	defer tx.Rollback()
	err = tx.QueryRow(query, args...).Scan(&infraction.Id, &createdAt, &infraction.Version)
	if err != nil {
		return ierrors.ErrDbFailure.Wrap(err)
	}
	infraction.CreatedAt = time.UnixMilli(createdAt)

	if len(infraction.Transactions) != 0 {
		queryTransaction := `insert into "transaction" (infraction_id, amount, created_at) VALUES `
		var argsTransaction []interface{}
		for _, transaction := range infraction.Transactions {
			queryTransaction += "(?, ?, ?), "
			argsTransaction = append(argsTransaction, infraction.Id)
			argsTransaction = append(argsTransaction, transaction.Amount)
			argsTransaction = append(argsTransaction, time.Now().UnixMilli())
		}
		queryTransaction, _ = strings.CutSuffix(queryTransaction, ", ")
		queryTransaction += " returning id, amount, created_at"
		rows, err := tx.Query(queryTransaction, argsTransaction...)
		if err != nil {
			return ierrors.ErrDbFailure.Wrap(err)
		}
		var transactions []Transaction
		defer rows.Close()
		for rows.Next() {
			var transaction Transaction
			var createdAt int64
			err := rows.Scan(&transaction.Id, &transaction.Amount, &createdAt)
			if err != nil {
				return ierrors.ErrDbFailure.Wrap(err)
			}
			transaction.CreatedAt = time.UnixMilli(createdAt)
			transactions = append(transactions, transaction)
		}
		infraction.Transactions = transactions
	}
	err = tx.Commit()
	if err != nil {
		return ierrors.ErrDbFailure.Wrap(err)
	}
	return nil
}

func (i *InfractionRepository) Select(id int64) (*Infraction, error) {
	query := `select i.id, i.created_at, i.split, i.name, i.version, t.id, t.amount, t.created_at 
              from infraction i left join "transaction" t on i.id = t.infraction_id 
              where i.id = ?`

	var infraction Infraction
	var createdAt int64
	rows, err := i.DB.Query(query, id)
	if err != nil {
		return nil, ierrors.ErrDbFailure.Wrap(err)
	}

	defer rows.Close()
	for rows.Next() {
		var transaction Transaction
		nillableTransaction := struct {
			id        *int64
			amount    *float64
			createdAt *int64
		}{}
		err := rows.Scan(&infraction.Id,
			&createdAt,
			&infraction.Split,
			&infraction.Name,
			&infraction.Version,
			&nillableTransaction.id,
			&nillableTransaction.amount,
			&nillableTransaction.createdAt)
		if err != nil {
			return nil, ierrors.ErrDbFailure.Wrap(err)
		}
		if nillableTransaction.id != nil &&
			nillableTransaction.amount != nil &&
			nillableTransaction.createdAt != nil {
			transaction.Id = *nillableTransaction.id
			transaction.Amount = *nillableTransaction.amount
			transaction.CreatedAt = time.UnixMilli(*nillableTransaction.createdAt)
			infraction.Transactions = append(infraction.Transactions, transaction)
		}
		infraction.CreatedAt = time.UnixMilli(createdAt)
	}

	if infraction.Version > 0 {
		return &infraction, nil
	}
	return nil, ierrors.ErrNoInfractionFound
}

func (i *InfractionRepository) SelectAll() (*[]Infraction, error) {
	query := `select i.id, i.version, i.created_at, i.name, i.split from infraction i`

	rows, err := i.DB.Query(query)
	if err != nil {
		return nil, ierrors.ErrDbFailure.Wrap(err)
	}

	infractions := make([]Infraction, 0)
	defer rows.Close()
	for rows.Next() {
		var infraction Infraction
		var createdAt int64
		err := rows.Scan(
			&infraction.Id,
			&infraction.Version,
			&createdAt,
			&infraction.Name,
			&infraction.Split)
		if err != nil {
			return nil, ierrors.ErrDbFailure.Wrap(err)
		}
		infraction.CreatedAt = time.UnixMilli(createdAt)

		infractions = append(infractions, infraction)
	}

	return &infractions, nil
}

func (i *InfractionRepository) Delete(id int64) error {
	query := `delete from infraction where id = ?`
	tx, err := i.DB.Begin()
	if err != nil {
		return ierrors.ErrDbFailure.Wrap(err)
	}
	defer tx.Rollback()

	res, err := tx.Exec(query, id)
	if err != nil {
		return ierrors.ErrDbFailure.Wrap(err)
	}
	if rows, _ := res.RowsAffected(); rows != 1 {
		return ierrors.ErrNoInfractionFound
	}
	err = tx.Commit()
	if err != nil {
		return ierrors.ErrDbFailure.Wrap(err)
	}
	return nil
}

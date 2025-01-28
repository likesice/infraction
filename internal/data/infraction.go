package data

import (
	"database/sql"
	"time"
)

type InfractionRepository struct {
	db *sql.DB
}
type Infraction struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	User      int64  `json:"-"`
	Version   int64  `json:"version"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func (i *InfractionRepository) Insert(infraction *Infraction) error {
	infraction.CreatedAt = time.Now().UnixMilli()
	infraction.UpdatedAt = time.Now().UnixMilli()
	infraction.Version = 1

	query := `
		INSERT INTO infractions (name, version, created_at, updated_at, user_id)
		VALUES (?, ?, ?, ?, ?)
        RETURNING id, version;`

	args := []interface{}{infraction.Name, infraction.Version, infraction.CreatedAt, infraction.UpdatedAt, infraction.User}

	return i.db.QueryRow(query, args...).Scan(&infraction.Id, &infraction.Version)
}

func (i *InfractionRepository) SelectAll(user *User) ([]*Infraction, error) {
	query := `SELECT name, version, created_at, updated_at, user_id 
FROM infractions WHERE user_id = ?;`

	args := []interface{}{user.Id}

	rows, err := i.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	infractions := make([]*Infraction, 0)
	for rows.Next() {
		var infraction Infraction
		err = rows.Scan(&infraction.Name, &infraction.Version,
			&infraction.CreatedAt, &infraction.UpdatedAt, &infraction.User)

		if err != nil {
			return nil, err
		}
		infractions = append(infractions, &infraction)
	}
	return infractions, nil
}

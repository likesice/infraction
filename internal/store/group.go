package store

import (
	"context"
	"database/sql"
)

type GroupStore struct {
	db *sql.DB
}
type Group struct {
	Id        int64   `json:"id"`
	Name      string  `json:"name"`
	Members   []int64 `json:"-"`
	Version   int64   `json:"version"`
	CreatedAt int64   `json:"created_at"`
	UpdatedAt int64   `json:"updated_at"`
}

func (i *GroupStore) Insert(group *Group) error {
	group.CreatedAt = now().UnixMilli()
	group.UpdatedAt = now().UnixMilli()
	group.Version = 1

	query := `
		INSERT INTO groups (name, version, created_at, updated_at)
		VALUES (?, ?, ?, ?)
        RETURNING id, version;`

	args := []interface{}{group.Name, group.Version, group.CreatedAt, group.UpdatedAt}

	tx, err := i.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}

	err = tx.QueryRow(query, args...).Scan(&group.Id, &group.Version)
	if err != nil {
		err := tx.Rollback()
		return err
	}
	query = `
		INSERT INTO groups_users (user_id, group_id)
		VALUES (?, ?);`
	args = []interface{}{group.Members[0], group.Id}
	exec, err := tx.Exec(query, args...)
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

func (i *GroupStore) SelectAll(user *User) ([]*Group, error) {
	query := `SELECT id, name, version, created_at, updated_at, user_id
FROM groups g JOIN groups_users gu ON g.id = gu.group_id WHERE gu.user_id = ?;`

	args := []interface{}{user.Id}

	rows, err := i.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make([]*Group, 0)
	for rows.Next() {
		var group Group
		var userId int64
		err = rows.Scan(&group.Id, &group.Name, &group.Version,
			&group.CreatedAt, &group.UpdatedAt, &userId)

		if err != nil {
			return nil, err
		}
		group.Members = append(group.Members, userId)
		groups = append(groups, &group)
	}
	return groups, nil
}

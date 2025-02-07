// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: role.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const getRoleByName = `-- name: GetRoleByName :one
SELECT id FROM role WHERE role.name = $1
`

func (q *Queries) GetRoleByName(ctx context.Context, name string) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getRoleByName, name)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

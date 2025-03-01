// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: users.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    role_id, 
    email,
    password_hash)
	VALUES ($1, $2, $3) RETURNING id, role_id, email, password_hash, is_enabled, created_at
`

type CreateUserParams struct {
	RoleID       uuid.UUID `json:"role_id"`
	Email        string    `json:"email"`
	PasswordHash []byte    `json:"password_hash"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.RoleID, arg.Email, arg.PasswordHash)
	var i User
	err := row.Scan(
		&i.ID,
		&i.RoleID,
		&i.Email,
		&i.PasswordHash,
		&i.IsEnabled,
		&i.CreatedAt,
	)
	return i, err
}

const deleteUserById = `-- name: DeleteUserById :exec
DELETE FROM users WHERE id = $1
`

func (q *Queries) DeleteUserById(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteUserById, id)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT u.id AS user_id ,e.id AS employee_id , e.name,e.avatar , u.email,
u.password_hash, r.name AS role_name , e.job_title , e.department
FROM users u
JOIN role r ON u.role_id = r.id
JOIN employee e ON u.id = e.user_id
WHERE u.email = $1
`

type GetUserByEmailRow struct {
	UserID       uuid.UUID `json:"user_id"`
	EmployeeID   uuid.UUID `json:"employee_id"`
	Name         string    `json:"name"`
	Avatar       string    `json:"avatar"`
	Email        string    `json:"email"`
	PasswordHash []byte    `json:"password_hash"`
	RoleName     string    `json:"role_name"`
	JobTitle     string    `json:"job_title"`
	Department   string    `json:"department"`
}

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i GetUserByEmailRow
	err := row.Scan(
		&i.UserID,
		&i.EmployeeID,
		&i.Name,
		&i.Avatar,
		&i.Email,
		&i.PasswordHash,
		&i.RoleName,
		&i.JobTitle,
		&i.Department,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, role_id, email, password_hash, is_enabled, created_at FROM  users  WHERE users.id = $1
`

func (q *Queries) GetUserById(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.RoleID,
		&i.Email,
		&i.PasswordHash,
		&i.IsEnabled,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByToken = `-- name: GetUserByToken :one
SELECT u.id AS user_id, u.email, u.is_enabled,
u.password_hash,r.name AS role_name
FROM users u
INNER JOIN token t ON u.id = t.user_id
INNER JOIN role r ON u.role_id = r.id
WHERE t.hash = $1
AND t.scope = $2
AND t.expiry > $3
`

type GetUserByTokenParams struct {
	Hash   []byte    `json:"hash"`
	Scope  string    `json:"scope"`
	Expiry time.Time `json:"expiry"`
}

type GetUserByTokenRow struct {
	UserID       uuid.UUID `json:"user_id"`
	Email        string    `json:"email"`
	IsEnabled    bool      `json:"is_enabled"`
	PasswordHash []byte    `json:"password_hash"`
	RoleName     string    `json:"role_name"`
}

func (q *Queries) GetUserByToken(ctx context.Context, arg GetUserByTokenParams) (GetUserByTokenRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByToken, arg.Hash, arg.Scope, arg.Expiry)
	var i GetUserByTokenRow
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.IsEnabled,
		&i.PasswordHash,
		&i.RoleName,
	)
	return i, err
}

const updateUserById = `-- name: UpdateUserById :exec
UPDATE users
SET
    role_id = $1,
    email = $2,
    is_enabled = $3,
    password_hash = $4
WHERE
    id = $5
`

type UpdateUserByIdParams struct {
	RoleID       uuid.UUID `json:"role_id"`
	Email        string    `json:"email"`
	IsEnabled    bool      `json:"is_enabled"`
	PasswordHash []byte    `json:"password_hash"`
	ID           uuid.UUID `json:"id"`
}

func (q *Queries) UpdateUserById(ctx context.Context, arg UpdateUserByIdParams) error {
	_, err := q.db.ExecContext(ctx, updateUserById,
		arg.RoleID,
		arg.Email,
		arg.IsEnabled,
		arg.PasswordHash,
		arg.ID,
	)
	return err
}

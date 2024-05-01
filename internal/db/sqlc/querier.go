// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateLeave(ctx context.Context, arg CreateLeaveParams) error
	CreateToken(ctx context.Context, arg CreateTokenParams) error
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteLeave(ctx context.Context, id uuid.UUID) error
	GetAllLeaves(ctx context.Context) ([]GetAllLeavesRow, error)
	GetLeaveById(ctx context.Context, id uuid.UUID) (Leave, error)
	GetLeaveByIdDetailed(ctx context.Context, id uuid.UUID) (GetLeaveByIdDetailedRow, error)
	GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error)
	GetUserByToken(ctx context.Context, arg GetUserByTokenParams) (GetUserByTokenRow, error)
	UpdateLeave(ctx context.Context, arg UpdateLeaveParams) error
}

var _ Querier = (*Queries)(nil)

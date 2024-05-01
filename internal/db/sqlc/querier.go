// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateContract(ctx context.Context, arg CreateContractParams) error
	CreateEmployee(ctx context.Context, arg CreateEmployeeParams) error
	CreateLeave(ctx context.Context, arg CreateLeaveParams) error
	CreateToken(ctx context.Context, arg CreateTokenParams) error
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteContract(ctx context.Context, id uuid.UUID) error
	DeleteEmployeeById(ctx context.Context, id uuid.UUID) error
	DeleteLeave(ctx context.Context, id uuid.UUID) error
	GetAllContracts(ctx context.Context) ([]GetAllContractsRow, error)
	GetAllEmployees(ctx context.Context) ([]GetAllEmployeesRow, error)
	GetAllLeaves(ctx context.Context) ([]GetAllLeavesRow, error)
	GetContractById(ctx context.Context, id uuid.UUID) (Contract, error)
	GetContractByIdDetailed(ctx context.Context, id uuid.UUID) (GetContractByIdDetailedRow, error)
	GetEmployeeById(ctx context.Context, id uuid.UUID) (Employee, error)
	GetEmployeeByIdDetailed(ctx context.Context, id uuid.UUID) (GetEmployeeByIdDetailedRow, error)
	GetLeaveById(ctx context.Context, id uuid.UUID) (Leave, error)
	GetLeaveByIdDetailed(ctx context.Context, id uuid.UUID) (GetLeaveByIdDetailedRow, error)
	GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error)
	GetUserByToken(ctx context.Context, arg GetUserByTokenParams) (GetUserByTokenRow, error)
	UpdateContract(ctx context.Context, arg UpdateContractParams) error
	UpdateEmployeeById(ctx context.Context, arg UpdateEmployeeByIdParams) error
	UpdateLeave(ctx context.Context, arg UpdateLeaveParams) error
}

var _ Querier = (*Queries)(nil)

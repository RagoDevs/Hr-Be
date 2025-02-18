package db

import (
	"context"
	"database/sql"
)

type Store interface {
	Querier
	TxnCreateUserEmployee(ctx context.Context, user CreateUserParams, emp CreateEmployeeParams) error
	TxnUpdateUserEmployee(ctx context.Context, user UpdateUserByIdParams, emp UpdateEmployeeByIdParams) error
	TxnDeleteUserEmployee(ctx context.Context, args TxnUserEmployeeDelete) error
	TxnAcceptedRejectLeave(ctx context.Context, args ApproveRejectLeaveParams) error
	GetAllPayroll(ctx context.Context, arg GetAllPayrollParams) ([]GetAllPayrollRow, error)
}

type SQLStore struct {
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

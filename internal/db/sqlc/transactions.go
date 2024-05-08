package db

import (
	"context"

	"github.com/google/uuid"
)

type TxnUserEmployeeDelete struct {
	UserId     uuid.UUID
	EmployeeId uuid.UUID
}

func (store *SQLStore) TxnCreateUserEmployee(ctx context.Context, user CreateUserParams, emp CreateEmployeeParams) error {

	tx, err := store.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := New(tx)

	usr, err := qtx.CreateUser(ctx, user)

	if err != nil {
		return err
	}

	emp.UserID = usr.ID

	if err = qtx.CreateEmployee(ctx, emp); err != nil {
		return err
	}

	return tx.Commit()

}

func (store *SQLStore) TxnUpdateUserEmployee(ctx context.Context, user UpdateUserByIdParams, emp UpdateEmployeeByIdParams) error {

	tx, err := store.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := New(tx)

	err = qtx.UpdateUserById(ctx, user)

	if err != nil {
		return err
	}
	if err = qtx.UpdateEmployeeById(ctx, emp); err != nil {
		return err
	}

	return tx.Commit()

}

func (store *SQLStore) TxnDeleteUserEmployee(ctx context.Context, args TxnUserEmployeeDelete) error {

	tx, err := store.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := New(tx)

	err = qtx.DeleteEmployeeById(ctx, args.EmployeeId)

	if err != nil {
		return err
	}
	if err = qtx.DeleteUserById(ctx, args.UserId); err != nil {
		return err
	}

	return tx.Commit()

}

func (store *SQLStore) TxnAcceptedRejectLeave(ctx context.Context, args ApproveRejectLeaveParams) error {

	tx, err := store.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := New(tx)

	err = qtx.ApproveRejectLeave(ctx, args)

	if err != nil {
		return err
	}

	if args.Approved {

		leave, err := qtx.GetLeaveById(ctx, args.ID)

		if err != nil {
			return err
		}

		emp := PresentAbsentEmployeeByIdParams{
			IsPresent: false,
			ID:        leave.EmployeeID,
		}

		err = qtx.PresentAbsentEmployeeById(ctx, emp)

		if err != nil {
			return err
		}

	}

	return tx.Commit()

}

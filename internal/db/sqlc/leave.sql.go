// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: leave.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const approveRejectLeave = `-- name: ApproveRejectLeave :exec
UPDATE leave
SET approved = $1,
    seen = TRUE
WHERE id = $2
`

type ApproveRejectLeaveParams struct {
	Approved bool      `json:"approved"`
	ID       uuid.UUID `json:"id"`
}

func (q *Queries) ApproveRejectLeave(ctx context.Context, arg ApproveRejectLeaveParams) error {
	_, err := q.db.ExecContext(ctx, approveRejectLeave, arg.Approved, arg.ID)
	return err
}

const createLeave = `-- name: CreateLeave :exec
INSERT INTO leave (
    employee_id, 
    approved_by_id,
    description,
    start_date,
    end_date,
    leave_count)
	VALUES ($1, $2, $3, $4, $5, $6)
`

type CreateLeaveParams struct {
	EmployeeID   uuid.UUID `json:"employee_id"`
	ApprovedByID uuid.UUID `json:"approved_by_id"`
	Description  string    `json:"description"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	LeaveCount   int16     `json:"leave_count"`
}

func (q *Queries) CreateLeave(ctx context.Context, arg CreateLeaveParams) error {
	_, err := q.db.ExecContext(ctx, createLeave,
		arg.EmployeeID,
		arg.ApprovedByID,
		arg.Description,
		arg.StartDate,
		arg.EndDate,
		arg.LeaveCount,
	)
	return err
}

const deleteLeave = `-- name: DeleteLeave :exec
DELETE FROM leave WHERE id = $1
`

func (q *Queries) DeleteLeave(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteLeave, id)
	return err
}

const getAllAApprovers = `-- name: GetAllAApprovers :many
SELECT id, name 
FROM employee 
WHERE LOWER(job_title) = 'hr'
`

type GetAllAApproversRow struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (q *Queries) GetAllAApprovers(ctx context.Context) ([]GetAllAApproversRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllAApprovers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllAApproversRow{}
	for rows.Next() {
		var i GetAllAApproversRow
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllLeavesApproved = `-- name: GetAllLeavesApproved :many
SELECT 
    l.id AS leave_id,
    e.name AS employee_name,
    e.department,
    u.email AS employee_email,
    l.employee_id, 
    l.approved_by_id ,
    ep.name AS approved_by_name,
    up.email AS approved_by_email,
    l.approved, 
    l.description, 
    l.start_date, 
    l.end_date, 
    l.leave_count, 
    l.created_at, 
    l.seen 
FROM 
    leave l
JOIN 
    employee e ON e.id = l.employee_id
JOIN
    users u ON e.user_id = u.id
JOIN 
    employee ep ON l.approved_by_id = ep.id
JOIN
    users up ON ep.user_id = up.id

WHERE l.seen = TRUE AND l.approved = TRUE

ORDER BY l.created_at DESC
`

type GetAllLeavesApprovedRow struct {
	LeaveID         uuid.UUID `json:"leave_id"`
	EmployeeName    string    `json:"employee_name"`
	Department      string    `json:"department"`
	EmployeeEmail   string    `json:"employee_email"`
	EmployeeID      uuid.UUID `json:"employee_id"`
	ApprovedByID    uuid.UUID `json:"approved_by_id"`
	ApprovedByName  string    `json:"approved_by_name"`
	ApprovedByEmail string    `json:"approved_by_email"`
	Approved        bool      `json:"approved"`
	Description     string    `json:"description"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	LeaveCount      int16     `json:"leave_count"`
	CreatedAt       time.Time `json:"created_at"`
	Seen            bool      `json:"seen"`
}

func (q *Queries) GetAllLeavesApproved(ctx context.Context) ([]GetAllLeavesApprovedRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllLeavesApproved)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllLeavesApprovedRow{}
	for rows.Next() {
		var i GetAllLeavesApprovedRow
		if err := rows.Scan(
			&i.LeaveID,
			&i.EmployeeName,
			&i.Department,
			&i.EmployeeEmail,
			&i.EmployeeID,
			&i.ApprovedByID,
			&i.ApprovedByName,
			&i.ApprovedByEmail,
			&i.Approved,
			&i.Description,
			&i.StartDate,
			&i.EndDate,
			&i.LeaveCount,
			&i.CreatedAt,
			&i.Seen,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllLeavesUnSeenTopBottomSeen = `-- name: GetAllLeavesUnSeenTopBottomSeen :many
SELECT 
    l.id AS leave_id,
    e.name AS employee_name,
    e.department,
    u.email AS employee_email,
    l.employee_id, 
    l.approved_by_id ,
    ep.name AS approved_by_name,
    up.email AS approved_by_email,
    l.approved, 
    l.description, 
    l.start_date, 
    l.end_date, 
    l.leave_count, 
    l.created_at, 
    l.seen 
FROM 
    leave l
JOIN 
    employee e ON e.id = l.employee_id
JOIN
    users u ON e.user_id = u.id
JOIN 
    employee ep ON l.approved_by_id = ep.id
JOIN
    users up ON ep.user_id = up.id

ORDER BY l.seen ASC, l.created_at DESC
`

type GetAllLeavesUnSeenTopBottomSeenRow struct {
	LeaveID         uuid.UUID `json:"leave_id"`
	EmployeeName    string    `json:"employee_name"`
	Department      string    `json:"department"`
	EmployeeEmail   string    `json:"employee_email"`
	EmployeeID      uuid.UUID `json:"employee_id"`
	ApprovedByID    uuid.UUID `json:"approved_by_id"`
	ApprovedByName  string    `json:"approved_by_name"`
	ApprovedByEmail string    `json:"approved_by_email"`
	Approved        bool      `json:"approved"`
	Description     string    `json:"description"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	LeaveCount      int16     `json:"leave_count"`
	CreatedAt       time.Time `json:"created_at"`
	Seen            bool      `json:"seen"`
}

func (q *Queries) GetAllLeavesUnSeenTopBottomSeen(ctx context.Context) ([]GetAllLeavesUnSeenTopBottomSeenRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllLeavesUnSeenTopBottomSeen)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllLeavesUnSeenTopBottomSeenRow{}
	for rows.Next() {
		var i GetAllLeavesUnSeenTopBottomSeenRow
		if err := rows.Scan(
			&i.LeaveID,
			&i.EmployeeName,
			&i.Department,
			&i.EmployeeEmail,
			&i.EmployeeID,
			&i.ApprovedByID,
			&i.ApprovedByName,
			&i.ApprovedByEmail,
			&i.Approved,
			&i.Description,
			&i.StartDate,
			&i.EndDate,
			&i.LeaveCount,
			&i.CreatedAt,
			&i.Seen,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLeaveById = `-- name: GetLeaveById :one
SELECT id, employee_id, approved_by_id, approved, description, start_date, end_date, leave_count, created_at, seen FROM  leave WHERE leave.id = $1
`

func (q *Queries) GetLeaveById(ctx context.Context, id uuid.UUID) (Leave, error) {
	row := q.db.QueryRowContext(ctx, getLeaveById, id)
	var i Leave
	err := row.Scan(
		&i.ID,
		&i.EmployeeID,
		&i.ApprovedByID,
		&i.Approved,
		&i.Description,
		&i.StartDate,
		&i.EndDate,
		&i.LeaveCount,
		&i.CreatedAt,
		&i.Seen,
	)
	return i, err
}

const getLeavesByEmployeeId = `-- name: GetLeavesByEmployeeId :many
SELECT 
    l.id AS leave_id,
    e.name AS approved_by_name,
    e.department,
    u.email AS approved_by_email,
    l.approved, 
    l.description, 
    l.start_date, 
    l.end_date, 
    l.leave_count, 
    l.created_at, 
    l.seen 
FROM 
    leave l

JOIN 
    employee e ON l.approved_by_id = e.id

JOIN
    users u ON e.user_id = u.id

WHERE 
    l.employee_id = $1
`

type GetLeavesByEmployeeIdRow struct {
	LeaveID         uuid.UUID `json:"leave_id"`
	ApprovedByName  string    `json:"approved_by_name"`
	Department      string    `json:"department"`
	ApprovedByEmail string    `json:"approved_by_email"`
	Approved        bool      `json:"approved"`
	Description     string    `json:"description"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	LeaveCount      int16     `json:"leave_count"`
	CreatedAt       time.Time `json:"created_at"`
	Seen            bool      `json:"seen"`
}

func (q *Queries) GetLeavesByEmployeeId(ctx context.Context, employeeID uuid.UUID) ([]GetLeavesByEmployeeIdRow, error) {
	rows, err := q.db.QueryContext(ctx, getLeavesByEmployeeId, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetLeavesByEmployeeIdRow{}
	for rows.Next() {
		var i GetLeavesByEmployeeIdRow
		if err := rows.Scan(
			&i.LeaveID,
			&i.ApprovedByName,
			&i.Department,
			&i.ApprovedByEmail,
			&i.Approved,
			&i.Description,
			&i.StartDate,
			&i.EndDate,
			&i.LeaveCount,
			&i.CreatedAt,
			&i.Seen,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateLeave = `-- name: UpdateLeave :exec
UPDATE leave
SET approved_by_id = $1,
    approved = $2,
    description= $3,
    start_date = $4,
    end_date = $5,
    leave_count = $6,
    seen = $7
WHERE id = $8
`

type UpdateLeaveParams struct {
	ApprovedByID uuid.UUID `json:"approved_by_id"`
	Approved     bool      `json:"approved"`
	Description  string    `json:"description"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	LeaveCount   int16     `json:"leave_count"`
	Seen         bool      `json:"seen"`
	ID           uuid.UUID `json:"id"`
}

func (q *Queries) UpdateLeave(ctx context.Context, arg UpdateLeaveParams) error {
	_, err := q.db.ExecContext(ctx, updateLeave,
		arg.ApprovedByID,
		arg.Approved,
		arg.Description,
		arg.StartDate,
		arg.EndDate,
		arg.LeaveCount,
		arg.Seen,
		arg.ID,
	)
	return err
}

// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: employee.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createEmployee = `-- name: CreateEmployee :exec
INSERT INTO employee (
    user_id, 
    name, 
    dob, 
    avatar, 
    phone, 
    gender, 
    job_title, 
    department, 
    address, 
    joining_date
)
VALUES (
    $1, 
    $2, 
    $3, 
    $4, 
    $5, 
    $6, 
    $7, 
    $8, 
    $9, 
    $10
)
`

type CreateEmployeeParams struct {
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Dob         time.Time `json:"dob"`
	Avatar      string    `json:"avatar"`
	Phone       string    `json:"phone"`
	Gender      string    `json:"gender"`
	JobTitle    string    `json:"job_title"`
	Department  string    `json:"department"`
	Address     string    `json:"address"`
	JoiningDate time.Time `json:"joining_date"`
}

func (q *Queries) CreateEmployee(ctx context.Context, arg CreateEmployeeParams) error {
	_, err := q.db.ExecContext(ctx, createEmployee,
		arg.UserID,
		arg.Name,
		arg.Dob,
		arg.Avatar,
		arg.Phone,
		arg.Gender,
		arg.JobTitle,
		arg.Department,
		arg.Address,
		arg.JoiningDate,
	)
	return err
}

const deleteEmployeeById = `-- name: DeleteEmployeeById :exec
DELETE FROM employee
WHERE id = $1
`

func (q *Queries) DeleteEmployeeById(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteEmployeeById, id)
	return err
}

const getAllEmployees = `-- name: GetAllEmployees :many
SELECT
    e.id AS employee_id,
    e.name AS employee_name,
    e.job_title ,
    e.department,
    is_present 

FROM
    employee e

ORDER BY
    e.created_at DESC
`

type GetAllEmployeesRow struct {
	EmployeeID   uuid.UUID `json:"employee_id"`
	EmployeeName string    `json:"employee_name"`
	JobTitle     string    `json:"job_title"`
	Department   string    `json:"department"`
	IsPresent    bool      `json:"is_present"`
}

func (q *Queries) GetAllEmployees(ctx context.Context) ([]GetAllEmployeesRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllEmployees)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllEmployeesRow{}
	for rows.Next() {
		var i GetAllEmployeesRow
		if err := rows.Scan(
			&i.EmployeeID,
			&i.EmployeeName,
			&i.JobTitle,
			&i.Department,
			&i.IsPresent,
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

const getEmployeeById = `-- name: GetEmployeeById :one
SELECT id, user_id, name, dob, avatar, phone, gender, job_title, department, address, is_present, joining_date, created_at FROM  employee  WHERE employee.id = $1
`

func (q *Queries) GetEmployeeById(ctx context.Context, id uuid.UUID) (Employee, error) {
	row := q.db.QueryRowContext(ctx, getEmployeeById, id)
	var i Employee
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Dob,
		&i.Avatar,
		&i.Phone,
		&i.Gender,
		&i.JobTitle,
		&i.Department,
		&i.Address,
		&i.IsPresent,
		&i.JoiningDate,
		&i.CreatedAt,
	)
	return i, err
}

const getEmployeeByIdDetailed = `-- name: GetEmployeeByIdDetailed :one
SELECT
    e.id AS employee_id,
    e.user_id ,
    u.email ,
    u.is_enabled ,
    u.role_id ,
    r.name AS role_name,
    e.name AS employee_name,
    e.dob ,
    e.avatar ,
    e.phone ,
    e.gender ,
    e.job_title ,
    e.department ,
    e.is_present, 
    e.address ,
    e.joining_date ,
    e.created_at 
FROM
    employee e
JOIN
    users u ON e.user_id = u.id

JOIN 
    role r ON r.id =  u.role_id

WHERE 
    e.id = $1
`

type GetEmployeeByIdDetailedRow struct {
	EmployeeID   uuid.UUID `json:"employee_id"`
	UserID       uuid.UUID `json:"user_id"`
	Email        string    `json:"email"`
	IsEnabled    bool      `json:"is_enabled"`
	RoleID       uuid.UUID `json:"role_id"`
	RoleName     string    `json:"role_name"`
	EmployeeName string    `json:"employee_name"`
	Dob          time.Time `json:"dob"`
	Avatar       string    `json:"avatar"`
	Phone        string    `json:"phone"`
	Gender       string    `json:"gender"`
	JobTitle     string    `json:"job_title"`
	Department   string    `json:"department"`
	IsPresent    bool      `json:"is_present"`
	Address      string    `json:"address"`
	JoiningDate  time.Time `json:"joining_date"`
	CreatedAt    time.Time `json:"created_at"`
}

func (q *Queries) GetEmployeeByIdDetailed(ctx context.Context, id uuid.UUID) (GetEmployeeByIdDetailedRow, error) {
	row := q.db.QueryRowContext(ctx, getEmployeeByIdDetailed, id)
	var i GetEmployeeByIdDetailedRow
	err := row.Scan(
		&i.EmployeeID,
		&i.UserID,
		&i.Email,
		&i.IsEnabled,
		&i.RoleID,
		&i.RoleName,
		&i.EmployeeName,
		&i.Dob,
		&i.Avatar,
		&i.Phone,
		&i.Gender,
		&i.JobTitle,
		&i.Department,
		&i.IsPresent,
		&i.Address,
		&i.JoiningDate,
		&i.CreatedAt,
	)
	return i, err
}

const presentAbsentEmployeeById = `-- name: PresentAbsentEmployeeById :exec
UPDATE employee
SET
    is_present = $1
WHERE
    id = $2
`

type PresentAbsentEmployeeByIdParams struct {
	IsPresent bool      `json:"is_present"`
	ID        uuid.UUID `json:"id"`
}

func (q *Queries) PresentAbsentEmployeeById(ctx context.Context, arg PresentAbsentEmployeeByIdParams) error {
	_, err := q.db.ExecContext(ctx, presentAbsentEmployeeById, arg.IsPresent, arg.ID)
	return err
}

const updateEmployeeById = `-- name: UpdateEmployeeById :exec
UPDATE employee
SET
    name = $1,
    dob = $2,
    avatar = $3,
    phone = $4,
    gender = $5,
    job_title = $6,
    department = $7,
    address = $8,
    joining_date = $9,
    is_present = $10

WHERE
    id = $11
`

type UpdateEmployeeByIdParams struct {
	Name        string    `json:"name"`
	Dob         time.Time `json:"dob"`
	Avatar      string    `json:"avatar"`
	Phone       string    `json:"phone"`
	Gender      string    `json:"gender"`
	JobTitle    string    `json:"job_title"`
	Department  string    `json:"department"`
	Address     string    `json:"address"`
	JoiningDate time.Time `json:"joining_date"`
	IsPresent   bool      `json:"is_present"`
	ID          uuid.UUID `json:"id"`
}

func (q *Queries) UpdateEmployeeById(ctx context.Context, arg UpdateEmployeeByIdParams) error {
	_, err := q.db.ExecContext(ctx, updateEmployeeById,
		arg.Name,
		arg.Dob,
		arg.Avatar,
		arg.Phone,
		arg.Gender,
		arg.JobTitle,
		arg.Department,
		arg.Address,
		arg.JoiningDate,
		arg.IsPresent,
		arg.ID,
	)
	return err
}

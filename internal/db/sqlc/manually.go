package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const getAllPayroll = `-- name: GetAllPayroll :many
SELECT 
    count(*) OVER(), payroll.id, payroll.employee_id, payroll.basic_salary, payroll.tin, payroll.bank_name, payroll.bank_account, payroll.is_active, payroll.created_at, payroll.updated_at, 
    employee.name AS employee_name,
    employee.department
FROM payroll
JOIN employee ON payroll.employee_id = employee.id
WHERE payroll.is_active = TRUE
AND (
    ($3 = '' OR payroll.bank_name ILIKE '%' || $3 || '%') AND
    ($4 = '' OR employee.department ILIKE '%' || $3 || '%') AND
    ($5 = '' OR employee.name ILIKE '%' || $3 || '%')
)

ORDER BY 
    payroll.created_at DESC 
LIMIT 
    $1 
OFFSET 
    $2
`

type GetAllPayrollParams struct {
	Limit        int32  `json:"limit"`
	Offset       int32  `json:"offset"`
	EmployeeName string `json:"employee_name"`
	Department   string `json:"department"`
	BankName     string `json:"bank_name"`
}

type GetAllPayrollRow struct {
	Count        int64     `json:"count"`
	ID           uuid.UUID `json:"id"`
	EmployeeID   uuid.UUID `json:"employee_id"`
	BasicSalary  string    `json:"basic_salary"`
	Tin          string    `json:"tin"`
	BankName     string    `json:"bank_name"`
	BankAccount  string    `json:"bank_account"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	EmployeeName string    `json:"employee_name"`
	Department   string    `json:"department"`
}

func (q *Queries) GetAllPayroll(ctx context.Context, arg GetAllPayrollParams) ([]GetAllPayrollRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllPayroll,
		arg.Limit,
		arg.Offset,
		arg.BankName,
		arg.Department,
		arg.EmployeeName,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllPayrollRow{}
	for rows.Next() {
		var i GetAllPayrollRow
		if err := rows.Scan(
			&i.ID,
			&i.EmployeeID,
			&i.BasicSalary,
			&i.Tin,
			&i.BankName,
			&i.BankAccount,
			&i.IsActive,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.EmployeeName,
			&i.Department,
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

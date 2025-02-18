-- name: CreatePayroll :exec
INSERT INTO payroll (
    employee_id, basic_salary, tin, bank_name, bank_account
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: UpdatePayroll :exec
UPDATE payroll
SET
    basic_salary = $2,
    tin = $3,
    bank_name = $4,
    bank_account = $5,
    is_active = $6,
    updated_at = NOW()
WHERE id = $1;


-- name: GetPayroll :one
SELECT 
    payroll.id AS payroll_id,
    payroll.employee_id,
    payroll.basic_salary,
    payroll.tin,
    payroll.bank_name,
    payroll.bank_account,
    payroll.is_active AS payroll_is_active,
    payroll.created_at AS payroll_created_at,
    payroll.updated_at AS payroll_updated_at,
    employee.user_id,
    employee.name AS employee_name,
    employee.dob,
    employee.avatar,
    employee.phone,
    employee.gender,
    employee.job_title,
    employee.department,
    employee.address,
    employee.is_present AS employee_is_present,
    employee.joining_date AS employee_joining_date,
    employee.created_at AS employee_created_at
FROM payroll
JOIN employee ON payroll.employee_id = employee.id
WHERE payroll.id = $1
AND payroll.is_active = TRUE;





-- name: GetAllPayroll :many
SELECT 
    payroll.*, 
    employee.name AS employee_name,
    employee.department
FROM payroll
JOIN employee ON payroll.employee_id = employee.id
AND payroll.is_active = TRUE;


-- name: DeletePayroll :exec
UPDATE payroll
SET
    is_active = FALSE
WHERE id = $1;


-- name: IsEmployeeExisting :one
SELECT EXISTS (
    SELECT 1 FROM payroll WHERE employee_id = $1
);

-- name: JustGetPayroll :one
SELECT *
FROM payroll;
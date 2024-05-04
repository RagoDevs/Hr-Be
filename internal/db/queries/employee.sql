-- name: CreateEmployee :exec
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
);

-- name: GetEmployeeById :one
SELECT * FROM  employee  WHERE employee.id = $1;


-- name: UpdateEmployeeById :exec
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
    joining_date = $9

WHERE
    id = $10;



-- name: DeleteEmployeeById :exec
DELETE FROM employee
WHERE id = $1;



-- name: GetAllEmployees :many
SELECT
    e.id AS employee_id,
    e.name AS employee_name,
    e.job_title ,
    e.department,
    is_present 

FROM
    employee e

ORDER BY
    e.created_at DESC;



-- name: GetEmployeeByIdDetailed :one
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
    e.id = $1;

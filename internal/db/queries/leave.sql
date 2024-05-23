-- name: CreateLeave :exec
INSERT INTO leave (
    employee_id, 
    approved_by_id,
    leave_type,
    description,
    start_date,
    end_date,
    leave_count)
	VALUES ($1, $2, $3, $4, $5, $6, $7);


-- name: GetLeaveById :one
SELECT * FROM  leave WHERE leave.id = $1;


-- name: UpdateLeave :exec
UPDATE leave
SET approved_by_id = $1,
    approved = $2,
    leave_type = $3,
    description= $4,
    start_date = $5,
    end_date = $6,
    leave_count = $7,
    seen = $8
WHERE id = $9;


-- name: DeleteLeave :exec
DELETE FROM leave WHERE id = $1;


-- name: GetAllLeavesUnSeenTopBottomSeen :many
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
    l.leave_type,
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

ORDER BY l.seen ASC, l.created_at DESC;


-- name: GetAllLeavesApproved :many
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
    l.leave_type,
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

ORDER BY l.created_at DESC;


-- name: GetLeavesByEmployeeId :many
SELECT 
    l.id AS leave_id,
    e.name AS approved_by_name,
    e.department,
    u.email AS approved_by_email,
    l.approved,
    l.leave_type,
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
    l.employee_id = $1;


-- name: ApproveRejectLeave :exec
UPDATE leave
SET approved = $1,
    seen = TRUE
WHERE id = $2;


-- name: GetAllAApprovers :many
SELECT id, name 
FROM employee 
WHERE LOWER(job_title) = 'hr';


-- name: GetEmployeesOnLeave :many
SELECT
    e.id AS employee_id,
    e.name AS employee_name, 
    l.leave_type,
    l.description,
    l.start_date,
    l.end_date
FROM
    leave l
JOIN
    employee e ON l.employee_id = e.id
WHERE
    l.start_date <= CURRENT_DATE AND
    l.end_date >= CURRENT_DATE AND
    l.approved = TRUE;


-- name: GetEmployeesOnLeaveInthreedays :many
SELECT
    e.id AS employee_id,
    e.name AS employee_name, 
    l.leave_type,
    l.description,
    l.start_date,
    l.end_date,
    l.leave_count,
    l.approved,
    l.seen
FROM
    leave l
JOIN
    employee e ON l.employee_id = e.id
WHERE
    l.start_date >= CURRENT_DATE + INTERVAL '1 day' AND
    l.start_date < CURRENT_DATE + INTERVAL '4 days' AND
    l.approved = TRUE;

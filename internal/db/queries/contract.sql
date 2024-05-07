-- name: CreateContract :exec
INSERT INTO contract (
    employee_id, 
    contract_type,
    period,  
    start_date,
    end_date,
    attachment)
	VALUES ($1, $2, $3, $4, $5, $6);


-- name: GetContractById :one
SELECT * FROM  contract  WHERE contract.id = $1;


-- name: UpdateContract :exec
UPDATE contract
SET contract_type = $1,
    period = $2,
    start_date = $3,
    end_date = $4,
    attachment = $5
WHERE id = $6;


-- name: DeleteContract :exec
DELETE FROM contract WHERE id = $1;


-- name: GetContractsOfEmployeeByEmployeeId :many
SELECT
    c.id AS contract_id,
    c.contract_type,
    c.period,  
    c.start_date,
    c.end_date,
    c.attachment
    
FROM 
    contract c
JOIN 
    employee e ON e.id = c.employee_id
JOIN
    users u ON e.user_id = u.id
WHERE 
    e.id = $1;



-- -- name: GetContractsOfEmployeeByEmployeeId :many
-- SELECT
--     c.id AS contract_id,
--     c.employee_id, 
--     e.name AS employee_name,
--     u.email AS employee_email,
--     e.avatar,
--     e.job_title,
--     e.department,
--     c.contract_type,
--     c.period,  
--     c.start_date,
--     c.end_date,
--     c.attachment
    
-- FROM 
--     contract c
-- JOIN 
--     employee e ON e.id = c.employee_id
-- JOIN
--     users u ON e.user_id = u.id
-- WHERE 
--     e.id = $1;
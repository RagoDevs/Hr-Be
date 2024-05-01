-- name: CreateUser :one
INSERT INTO users (
    role_id, 
    email,
    password_hash)
	VALUES ($1, $2, $3) RETURNING *;


-- name: GetUserByEmail :one
SELECT u.id AS user_id ,e.id AS employee_id , e.name,e.avatar , u.email,
u.password_hash, r.name AS role_name , e.job_title , e.department
FROM users u
JOIN role r ON u.role_id = r.id
JOIN employee e ON u.id = e.user_id
WHERE u.email = $1;


-- name: GetUserByToken :one
SELECT u.id AS user_id, u.email, u.is_enabled,
u.password_hash,r.name AS role_name
FROM users u
INNER JOIN token t ON u.id = t.user_id
INNER JOIN role r ON u.role_id = r.id
WHERE t.hash = $1
AND t.scope = $2
AND t.expiry > $3;


-- name: GetUserById :one
SELECT * FROM  users  WHERE users.id = $1;


-- name: UpdateUserById :exec
UPDATE users
SET
    role_id = $1,
    email = $2,
    is_enabled = $3
WHERE
    id = $4;


-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;



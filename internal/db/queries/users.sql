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



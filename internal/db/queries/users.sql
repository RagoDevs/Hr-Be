-- name: CreateUser :one
INSERT INTO users (
    role_id, 
    email,
    password_hash)
	VALUES ($1, $2, $3) RETURNING *;

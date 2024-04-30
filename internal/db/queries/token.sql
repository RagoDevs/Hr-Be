-- name: CreateToken :exec
INSERT INTO token (hash, user_id, expiry, scope)
VALUES ($1, $2, $3, $4);

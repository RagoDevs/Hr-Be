-- name: GetRoleByName :one
SELECT id FROM role WHERE role.name = $1;
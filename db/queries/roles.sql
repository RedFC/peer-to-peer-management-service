-- name: CreateRole :one
INSERT INTO roles (name, name_hash, description, is_deleted, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListRoles :many
SELECT * FROM roles WHERE is_deleted = FALSE ORDER BY created_at DESC;

-- name: GetRoleByID :one
SELECT * FROM roles WHERE id = $1 AND is_deleted = FALSE;

-- name: GetRoleByName :one
SELECT * FROM roles WHERE name_hash = $1;

-- name: UpdateRole :one
UPDATE roles
SET name = $2, name_hash = $3, description = $4, updated_at = $5
WHERE id = $1
RETURNING *;

-- name: DeleteRole :exec
UPDATE roles
SET is_deleted = TRUE, updated_at = $2
WHERE id = $1;


-- name: CreateGroup :one
INSERT INTO groups (name, name_hash, metadata, is_deleted, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListGroups :many
SELECT * FROM groups WHERE is_deleted = FALSE
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetTotalGroupsCount :one
SELECT COUNT(*) AS total_count FROM groups WHERE is_deleted = FALSE;

-- name: GetGroupByID :one
SELECT * FROM groups WHERE id = $1;

-- name: GetGroupByName :one
SELECT * FROM groups WHERE name_hash = $1;

-- name: UpdateGroup :one
UPDATE groups
SET name = $2, name_hash = $3, metadata = $4, updated_at = $5
WHERE id = $1
RETURNING *;

-- name: DeleteGroup :exec
UPDATE groups
SET is_deleted = TRUE, updated_at = $2
WHERE id = $1;

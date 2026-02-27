-- name: AssignGroupToUser :exec
INSERT INTO user_groups (user_id, group_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemoveGroupFromUser :exec
DELETE FROM user_groups WHERE user_id = $1 AND group_id = $2;

-- name: GetUserGroups :many
SELECT g.* FROM groups g
JOIN user_groups ug ON ug.group_id = g.id
WHERE ug.user_id = $1;

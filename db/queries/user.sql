-- name: CreateUser :one
INSERT INTO users (email, email_hash, password, first_name, last_name, is_deleted, is_active, is_password_reset, last_login, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetUserByEmailHash :one
SELECT * FROM users WHERE email_hash = $1;

-- name: GetUserByID :one
SELECT 
    u.*, 
    COALESCE(
        (SELECT json_agg(
            json_build_object(
                'id', r.id,
                'name', r.name,
                'description', r.description,
                'is_deleted', r.is_deleted,
                'created_at', r.created_at,
                'updated_at', r.updated_at
            )
        ) FROM roles r 
         INNER JOIN user_roles ur ON ur.role_id = r.id 
         WHERE ur.user_id = u.id), 
        '[]'
    ) AS roles,
    COALESCE(
        (SELECT json_agg(
            json_build_object(
                'id', g.id,
                'name', g.name,
                'metadata', g.metadata,
                'is_deleted', g.is_deleted,
                'created_at', g.created_at,
                'updated_at', g.updated_at
            )
        ) FROM groups g 
         INNER JOIN user_groups ug ON ug.group_id = g.id 
         WHERE ug.user_id = u.id), 
        '[]'
    ) AS groups,
    COALESCE(
        (SELECT row_to_json(p) FROM profiles p WHERE p.user_id = u.id LIMIT 1),
        '{}'::json
    ) AS profile
FROM users u
WHERE u.id = $1;

-- name: ListUsers :many
SELECT 
    u.*, 
    COALESCE(
        (
            SELECT json_agg(
                json_build_object(
                'id', r.id,
                'name', r.name,
                'description', r.description,
                'is_deleted', r.is_deleted,
                'created_at', r.created_at,
                'updated_at', r.updated_at
            )
            )
            FROM roles r
            INNER JOIN user_roles ur ON ur.role_id = r.id
            WHERE ur.user_id = u.id 
              AND r.is_deleted = FALSE
        ), 
        '[]'
    ) AS roles,
    COALESCE(
        (
            SELECT json_agg(
                json_build_object(
                'id', g.id,
                'name', g.name,
                'metadata', g.metadata,
                'is_deleted', g.is_deleted,
                'created_at', g.created_at,
                'updated_at', g.updated_at
            )
            )
            FROM groups g
            INNER JOIN user_groups ug ON ug.group_id = g.id
            WHERE ug.user_id = u.id
              AND g.is_deleted = FALSE
        ), 
        '[]'
    ) AS groups
FROM users u
WHERE 
    u.is_deleted = FALSE   -- ✅ filter main users table
    AND NOT EXISTS (
        SELECT 1
        FROM roles r
        INNER JOIN user_roles ur ON ur.role_id = r.id
        WHERE ur.user_id = u.id
          AND r.name_hash = ANY($3::text[])
          AND r.is_deleted = FALSE
    )
    AND (
        $4::uuid = '00000000-0000-0000-0000-000000000000'
        OR EXISTS (
            SELECT 1 FROM user_roles ur WHERE ur.user_id = u.id AND ur.role_id = $4::uuid
        )
    )
    AND (
        $5::uuid = '00000000-0000-0000-0000-000000000000'
        OR EXISTS (
            SELECT 1 FROM user_groups ug WHERE ug.user_id = u.id AND ug.group_id = $5::uuid
        )
    )   
ORDER BY u.created_at DESC
LIMIT $1 OFFSET $2;


-- name: GetTotalUsersCount :one
SELECT COUNT(*) AS total_count
FROM users u
    WHERE 
        u.is_deleted = FALSE   -- ✅ filter main users table
    AND NOT EXISTS (
        SELECT 1
        FROM roles r
        INNER JOIN user_roles ur ON ur.role_id = r.id
        WHERE ur.user_id = u.id
        AND r.name_hash = ANY($1::text[])
    )
    AND (
        $2::uuid = '00000000-0000-0000-0000-000000000000'
        OR EXISTS (
            SELECT 1 FROM user_roles ur WHERE ur.user_id = u.id AND ur.role_id = $2::uuid
        )
    )
    AND (
        $3::uuid = '00000000-0000-0000-0000-000000000000'
        OR EXISTS (
            SELECT 1 FROM user_groups ug WHERE ug.user_id = u.id AND ug.group_id = $3::uuid
        )
    ); 

-- name: DeleteUser :exec
UPDATE users
SET is_deleted = TRUE, updated_at = $2
WHERE id = $1;


-- name: UpdateUserProfile :one
UPDATE users
SET first_name = $2, last_name = $3, last_login = $4, updated_at = $5
WHERE id = $1
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET email = $2, email_hash = $3, first_name = $4, last_name = $5, last_login = $6, updated_at = $7
WHERE id = $1
RETURNING *;

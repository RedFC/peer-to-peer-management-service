-- name: CreateProfile :one
INSERT INTO profiles (user_id, full_name, full_address, phone)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetProfileByUserID :one
SELECT * FROM profiles WHERE user_id = $1;

-- name: UpdateProfile :exec
UPDATE profiles
SET full_name = $2, full_address = $3, phone = $4
WHERE user_id = $1;

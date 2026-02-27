-- name: CreateSubscription :one
INSERT INTO subscriptions (user_id, group_id, subscription_status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListSubscriptions :many
SELECT
  s.id AS subscription_id,
  s.subscription_status,
  s.created_at AS subscription_created_at,
  s.updated_at AS subscription_updated_at,
  -- user details
  u.id AS user_id,
  u.first_name AS user_first_name,
  u.last_name AS user_last_name,
  u.email AS user_email,
  u.created_at AS user_created_at,
  u.updated_at AS user_updated_at,
  -- group details
  g.id AS group_id,
  g.name AS group_name,
  g.metadata AS group_metadata,
  g.created_at AS group_created_at,
  g.updated_at AS group_updated_at
FROM subscriptions s
JOIN users u ON s.user_id = u.id
JOIN groups g ON s.group_id = g.id
WHERE s.subscription_status = 'active' OR s.subscription_status = 'inactive'
ORDER BY s.created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetTotalSubscriptionsCount :one
SELECT COUNT(*) AS total_count FROM subscriptions WHERE is_deleted = FALSE;

-- name: GetSubscriptionByID :one
SELECT
  s.id AS subscription_id,
  s.subscription_status,
  s.created_at AS subscription_created_at,
  s.updated_at AS subscription_updated_at,
  -- user details
  u.id AS user_id,
  u.first_name AS user_first_name,
  u.last_name AS user_last_name,
  u.email AS user_email,
  u.created_at AS user_created_at,
  u.updated_at AS user_updated_at,
  -- group details
  g.id AS group_id,
  g.name AS group_name,
  g.metadata AS group_metadata,
  g.created_at AS group_created_at,
  g.updated_at AS group_updated_at
FROM subscriptions s
JOIN users u ON s.user_id = u.id
JOIN groups g ON s.group_id = g.id
WHERE s.id = $1
  AND s.is_deleted = FALSE;

-- name: GetSubscriptionByUserIDAndGroupID :many
SELECT * FROM subscriptions WHERE user_id = $1 AND group_id = $2 AND is_deleted = FALSE;

-- name: GetSubscriptionByUserIDAndGroupIDAndStatus :one
SELECT * FROM subscriptions WHERE user_id = $1 AND group_id = $2 AND subscription_status = $3;

-- name: GetSubscriptionsByUserID :many
SELECT * FROM subscriptions WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateSubscription :one
UPDATE subscriptions
SET subscription_status = $2, updated_at = $3
WHERE id = $1
RETURNING *;

-- name: RevokeSubscription :one
UPDATE subscriptions
SET subscription_status = 'inactive', updated_at = $2
WHERE id = $1
RETURNING *;

-- name: DeleteSubscription :exec
UPDATE subscriptions
SET is_deleted = TRUE, updated_at = $2, subscription_status = 'deleted'
WHERE id = $1;
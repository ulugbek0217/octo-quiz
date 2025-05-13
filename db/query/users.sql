-- name: CreateUser :one
INSERT INTO users (
    user_id, telegram_username, full_name, username, "role", phone
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE user_id = $1
LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;
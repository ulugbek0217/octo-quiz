-- name: CreateTeacher :one
INSERT INTO teachers (
    telegram_id, full_name
) VALUES (
    $1, $2
)
RETURNING *;

-- name: DeleteTeacher :exec
DELETE FROM teachers
WHERE id = $1;
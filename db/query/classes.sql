-- name: CreateClass :one
INSERT INTO classes (
    class_name, teacher_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: ListClasses :many
SELECT * FROM classes
WHERE teacher_id = $1
LIMIT $2
OFFSET $3;

-- name: DeleteClass :exec
DELETE FROM classes
WHERE class_id = $1;

-- name: CreateTestGroup :one
INSERT INTO test_groups (
    teacher_id, group_name
) VALUES (
    $1, $2
)
RETURNING *;
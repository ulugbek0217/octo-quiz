-- name: CreateTest :one
INSERT INTO test_groups (
    teacher_id, group_name
) VALUES (
    $1, $2
) RETURNING *;
INSERT INTO tests (
    group_id, english_word, uzbek_word
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetTestsByGroup :many
SELECT * FROM tests
WHERE group_id = $1;

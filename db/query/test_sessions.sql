-- name: NewTestSession :one
INSERT INTO test_sessions (
    student_id, test_set_id, start_time, correct_count, incorrect_count, completed
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;
-- name: NewStudentProgress :one
INSERT INTO student_progress (
    student_id, test_set_id, words_id, correct_count, incorrect_count
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;
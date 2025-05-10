-- name: AddStudentToClass :one
INSERT INTO class_students (
    class_id, student_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: ListClassStudents :many
SELECT * FROM class_students
WHERE class_id = $1
LIMIT $2
OFFSET $3;

-- name: DeleteStudentFromClass :exec
DELETE FROM class_students
WHERE student_id = $1;
// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: classes.sql

package db

import (
	"context"
)

const createClass = `-- name: CreateClass :one
INSERT INTO classes (
    class_name, teacher_id
) VALUES (
    $1, $2
) RETURNING class_id, class_name, teacher_id
`

type CreateClassParams struct {
	ClassName string `json:"class_name"`
	TeacherID int64  `json:"teacher_id"`
}

func (q *Queries) CreateClass(ctx context.Context, db DBTX, arg CreateClassParams) (Class, error) {
	row := db.QueryRow(ctx, createClass, arg.ClassName, arg.TeacherID)
	var i Class
	err := row.Scan(&i.ClassID, &i.ClassName, &i.TeacherID)
	return i, err
}

const deleteClass = `-- name: DeleteClass :exec
DELETE FROM classes
WHERE class_id = $1
`

func (q *Queries) DeleteClass(ctx context.Context, db DBTX, classID int64) error {
	_, err := db.Exec(ctx, deleteClass, classID)
	return err
}

const listClasses = `-- name: ListClasses :many
SELECT class_id, class_name, teacher_id FROM classes
WHERE teacher_id = $1
LIMIT $2
OFFSET $3
`

type ListClassesParams struct {
	TeacherID int64 `json:"teacher_id"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

func (q *Queries) ListClasses(ctx context.Context, db DBTX, arg ListClassesParams) ([]Class, error) {
	rows, err := db.Query(ctx, listClasses, arg.TeacherID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Class{}
	for rows.Next() {
		var i Class
		if err := rows.Scan(&i.ClassID, &i.ClassName, &i.TeacherID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

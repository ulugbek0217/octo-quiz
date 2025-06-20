// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Class struct {
	ClassID   int64              `json:"class_id"`
	ClassName string             `json:"class_name"`
	TeacherID int64              `json:"teacher_id"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}

type ClassStudent struct {
	ClassID   int64              `json:"class_id"`
	StudentID int64              `json:"student_id"`
	AddedAt   pgtype.Timestamptz `json:"added_at"`
}

type ClassTestSet struct {
	ClassID   int64 `json:"class_id"`
	TestSetID int64 `json:"test_set_id"`
}

type StudentProgress struct {
	ProgressID     int64       `json:"progress_id"`
	StudentID      int64       `json:"student_id"`
	TestSetID      int64       `json:"test_set_id"`
	WordsID        int64       `json:"words_id"`
	CorrectCount   pgtype.Int4 `json:"correct_count"`
	IncorrectCount pgtype.Int4 `json:"incorrect_count"`
}

type TestProgress struct {
	ProgressID int64 `json:"progress_id"`
	StudentID  int64 `json:"student_id"`
	TestSetID  int64 `json:"test_set_id"`
	WordsID    int64 `json:"words_id"`
	Completed  bool  `json:"completed"`
}

type TestSession struct {
	SessionID int64 `json:"session_id"`
	StudentID int64 `json:"student_id"`
	TestSetID int64 `json:"test_set_id"`
	// Unix timestamp
	StartTime      pgtype.Timestamptz `json:"start_time"`
	CorrectCount   int32              `json:"correct_count"`
	IncorrectCount int32              `json:"incorrect_count"`
	// True if all words completed
	Completed bool `json:"completed"`
}

type TestSet struct {
	TestSetID   int64  `json:"test_set_id"`
	TestSetName string `json:"test_set_name"`
	CreatorID   int64  `json:"creator_id"`
	IsPublic    bool   `json:"is_public"`
	// Seconds, NULL if no limit
	TimeLimit pgtype.Int4 `json:"time_limit"`
}

// CHECK (role IN ("student", "teacher"))
type User struct {
	UserID           int64       `json:"user_id"`
	TelegramUsername pgtype.Text `json:"telegram_username"`
	FullName         string      `json:"full_name"`
	Username         string      `json:"username"`
	// Must be student or teacher
	Role      string             `json:"role"`
	Phone     string             `json:"phone"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}

type Word struct {
	WordsID     int64  `json:"words_id"`
	TestSetID   int64  `json:"test_set_id"`
	EnglishWord string `json:"english_word"`
	UzbekWord   string `json:"uzbek_word"`
}

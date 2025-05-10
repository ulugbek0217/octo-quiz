package handlers

import (
	"sync"

	db "github.com/ulugbek0217/octo-quiz/db/sqlc"
)

type App struct {
	Store db.Store
	WG    *sync.WaitGroup
	Root  int64
}

type State int

const (
	NewUser State = iota
	RegisterFullName
	RegisterUsername
	RegisterRole
	RegisterPhoneNumber
	DoTest
)

package handlers

import (
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/fsm"
	db "github.com/ulugbek0217/octo-quiz/db/sqlc"
)

type App struct {
	Store db.Store
	WG    *sync.WaitGroup
	Root  int64
	F     *fsm.FSM
	B     *bot.Bot
}

// Available user states
const (
	StateDefault fsm.StateID = "default"
	// StateStart       fsm.StateID = "start"
	StateAskName     fsm.StateID = "ask_name"
	StateAskUsername fsm.StateID = "ask_username"
	StateAskRole     fsm.StateID = "ask_role"
	StateAskPhone    fsm.StateID = "ask_phone"
	StateFinish      fsm.StateID = "finish"
)

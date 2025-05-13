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
	StateAskName            fsm.StateID = "ask_name"
	StateAskUsername        fsm.StateID = "ask_username"
	StateAskRole            fsm.StateID = "ask_role"
	StateAskPhone           fsm.StateID = "ask_phone"
	StateFinishRegistration fsm.StateID = "finish_registration"

	StateAskTestSetName               fsm.StateID = "ask_test_set_name"
	StateAskTestSetType               fsm.StateID = "ask_test_set_type"
	StateAskTestSetTimeLimitAndFinish fsm.StateID = "ask_test_set_time_limit_and_finish"
	// StateFinishTestSetCreating fsm.StateID = "finish_test_set_creating"
)

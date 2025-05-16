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

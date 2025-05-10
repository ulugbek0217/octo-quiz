package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/ulugbek0217/octo-quiz/db/sqlc"
	"github.com/ulugbek0217/octo-quiz/handlers"
	"github.com/ulugbek0217/octo-quiz/util"
)

func main() {
	err := util.LoadEnv(".env")
	if err != nil {
		log.Fatalf("err loading .env: %v\n", err)
	}
	dbPath := os.Getenv("DB_URL")
	rootID, err := strconv.ParseInt(os.Getenv("ROOT"), 10, 64)
	if err != nil {
		log.Fatalf("couldn't parse root id: %v\n", err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	conn, err := pgxpool.New(context.Background(), dbPath)
	if err != nil {
		log.Fatalf("err connecting to db: %v\n", err)
	}
	defer conn.Close()

	store := db.NewStore(conn)
	wg := sync.WaitGroup{}
	app := handlers.App{
		Store: store,
		WG:    &wg,
		Root:  rootID,
	}

	options := []bot.Option{
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, app.Start),
	}

	bot, err := bot.New(os.Getenv("BOT_TOKEN"), options...)
	if err != nil {
		log.Fatalf("err creating bot: %v\n", err)
	}
	log.Println("Starting...")
	bot.Start(ctx)
}

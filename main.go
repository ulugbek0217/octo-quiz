package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/fsm"
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
	rootID, err := strconv.ParseInt(os.Getenv("ROOT_USER"), 10, 64)
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
	wg := sync.WaitGroup{}

	store := db.NewStore(conn)
	app := handlers.App{
		Store: store,
		Root:  rootID,
		WG:    &wg,
	}
	app.F = fsm.New(
		handlers.StateDefault,
		map[fsm.StateID]fsm.Callback{
			handlers.StateAskName:            app.CallbackName,
			handlers.StateAskUsername:        app.CallbackUsername,
			handlers.StateAskRole:            app.CallbackRole,
			handlers.StateAskPhone:           app.CallbackPhone,
			handlers.StateFinishRegistration: app.CallbackFinish,

			handlers.StateAskTestSetName:               app.CallbackTestSetName,
			handlers.StateAskTestSetType:               app.CallbackTestSetType,
			handlers.StateAskTestSetTimeLimitAndFinish: app.CallbackTestSetTimeLimit,
			handlers.StateAskClassName:                 app.CallbackAskClassName,
			// handlers.StateFinishTestSetCreating:        app.CallbackFinishTestSetCreating,
		},
	)
	options := []bot.Option{
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, app.Start),
		bot.WithCallbackQueryDataHandler("create_test_set", bot.MatchTypeExact, app.CreateTestSet),
		bot.WithCallbackQueryDataHandler("test_sets_page", bot.MatchTypePrefix, app.TeacherTestSetsList),
		bot.WithCallbackQueryDataHandler("teacher_test_set", bot.MatchTypePrefix, app.TeacherTestSetOptions),
		bot.WithCallbackQueryDataHandler("insert_words_into", bot.MatchTypePrefix, app.TeacherInsertWordsIntoTestSet),
		bot.WithCallbackQueryDataHandler("finish_inserting_words_into", bot.MatchTypePrefix, app.TeacherFinishWordInserting),
		bot.WithCallbackQueryDataHandler("delete_test_set", bot.MatchTypePrefix, app.TeacherDeleteTestSet),
		bot.WithCallbackQueryDataHandler("teacher_create_class", bot.MatchTypeExact, app.TeacherCreateClass),
		bot.WithCallbackQueryDataHandler("teacher_classes_page", bot.MatchTypePrefix, app.TeacherClassesList),
		bot.WithCallbackQueryDataHandler("teacher_class", bot.MatchTypePrefix, app.TeacherClassOptions),
		bot.WithMessageTextHandler("attc", bot.MatchTypeCommand, app.AddTestSetToClass),
		bot.WithMessageTextHandler("astc", bot.MatchTypeCommand, app.AddStudentToClass),
		bot.WithCallbackQueryDataHandler("student_test_sets_page", bot.MatchTypePrefix, app.StudentTestSetsList),
		// bot.WithCallbackQueryDataHandler("dashboard", bot.MatchTypeExact)
		// bot.WithCallbackQueryDataHandler("teacher_test_sets_list", bot.MatchTypeExact, app.TeacherTestSetsList),
		bot.WithDefaultHandler(app.MainHandler),
		// bot.WithWorkers(10),
	}
	bot, err := bot.New(os.Getenv("BOT_TOKEN"), options...)
	if err != nil {
		log.Fatalf("err creating bot: %v\n", err)
	}

	app.B = bot

	log.Println("Starting...")
	app.B.Start(ctx)
}

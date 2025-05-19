package handlers

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/ulugbek0217/octo-quiz/builder"
	db "github.com/ulugbek0217/octo-quiz/db/sqlc"
)

func (app *App) StudentTestSetsList(ctx context.Context, b *bot.Bot, u *models.Update) {
	if u.CallbackQuery == nil {
		return
	}
	userID := u.CallbackQuery.From.ID
	chatID := userID
	offset, err := strconv.ParseInt(strings.TrimPrefix(u.CallbackQuery.Data, "test_sets_page_"), 10, 32)
	if err != nil {
		log.Fatalf("err parsing offset: %v\n", err)
	}

	testSetsCount, err := app.Store.GetTestSetsCount(ctx, app.Store.Pool, userID)
	if err != nil {
		log.Fatalf("err counting test sets: %v\n", err)
	}
	if testSetsCount == 0 {
		app.B.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: u.CallbackQuery.ID,
			Text:            "Testlar mavjud emas ❌",
			ShowAlert:       true,
		})

		msgToDelete[userID] = append(msgToDelete[userID], u.CallbackQuery.Message.Message.ID)
		app.DashBoard(ctx, b, u, userID)
		return
	}
	log.Printf("next button offset: %d\n", offset)
	arg := db.ListTestSetsByCreatorIDParams{
		CreatorID: userID,
		Limit:     5,
		Offset:    int32(offset),
	}
	testSetsList, err := app.Store.ListTestSetsByCreatorID(ctx, app.Store.Pool, arg)
	if err != nil {
		log.Fatalf("err getting test sets: %v\n", err)
	}

	isLastPage := false
	if math.Ceil(float64(testSetsCount)/float64(arg.Limit))-1 == float64(offset) {
		isLastPage = true
	}
	var ids []int64
	var testSetNames string
	for id, set := range testSetsList {
		ids = append(ids, set.TestSetID)
		testSetNames = fmt.Sprintf("%s\n%d. %s", testSetNames, int(offset)+id+1, set.TestSetName)
	}
	fmt.Printf("got test sets: %v\n", ids)
	options := builder.InlinePaginatorParams{
		ItemCallback:      "teacher_test_set",
		NavigatorCallback: "test_sets_page",
	}
	msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        testSetNames,
		ReplyMarkup: builder.NewInlinePaginator(ids, arg.Offset, isLastPage, options),
	})
	if err != nil {
		log.Printf("err sending sets list: %v\n", err)
	}
	b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    chatID,
		MessageID: u.CallbackQuery.Message.Message.ID,
	})
	msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
}

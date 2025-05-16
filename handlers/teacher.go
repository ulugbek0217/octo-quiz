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

func (app *App) CreateTestSet(ctx context.Context, b *bot.Bot, u *models.Update) {
	userID := u.CallbackQuery.From.ID
	chatID := userID
	app.F.Transition(userID, StateAskTestSetName, userID, chatID)
	app.B.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    chatID,
		MessageID: u.CallbackQuery.Message.Message.ID,
	})
}

func (app *App) TeacherTestSetsList(ctx context.Context, b *bot.Bot, u *models.Update) {
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
		return
	}
	log.Printf("next button offset: %d\n", offset)
	arg := db.GetTestSetsByCreatorIDParams{
		CreatorID: userID,
		Limit:     5,
		Offset:    int32(offset),
	}
	testSetsList, err := app.Store.GetTestSetsByCreatorID(ctx, app.Store.Pool, arg)
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

	msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        testSetNames,
		ReplyMarkup: builder.TeacherInlineKeyboardPaginator(ids, arg.Offset, isLastPage),
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

func (app *App) TeacherTestSetOptions(ctx context.Context, b *bot.Bot, u *models.Update) {
	if u.CallbackQuery == nil {
		return
	}

	testSetID, err := strconv.ParseInt(strings.TrimPrefix(u.CallbackQuery.Data, "teacher_test_set_"), 10, 64)
	if err != nil {
		log.Fatalf("err converting test set id: %v\n", err)
	}

	userID := u.CallbackQuery.From.ID

	testSet, err := app.Store.GetTestSetByID(ctx, app.Store.Pool, testSetID)
	if err != nil {
		log.Fatalf("err getting test set from db: %v\n", err)
	}

	b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      userID,
		MessageID:   u.CallbackQuery.Message.Message.ID,
		Text:        testSet.TestSetName,
		ReplyMarkup: builder.TeacherInlineKeyboardTestSetOptions(testSetID),
	})

}

func (app *App) InsertWordsIntoTestSet(ctx context.Context, b *bot.Bot, u *models.Update) {
	if u.CallbackQuery == nil {
		return
	}
	userID := u.CallbackQuery.From.ID
	chatID := userID

	app.F.Transition(userID, StateInsertWordsIntoTestSet, userID, chatID)
}

func (app *App) FinishInsertingWords(ctx context.Context, b *bot.Bot, u *models.Update) {
	if u.CallbackQuery == nil {
		return
	}

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: u.CallbackQuery.ID,
		Text:            "Saqlandi ✅",
	})

	app.F.Transition(u.CallbackQuery.From.ID, StateDefault)
	app.TeacherTestSetOptions(ctx, b, u)
}

func (app *App) TeacherFinishWordInserting(ctx context.Context, b *bot.Bot, u *models.Update) {

}

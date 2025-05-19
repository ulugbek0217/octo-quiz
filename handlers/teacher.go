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

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        testSet.TestSetName,
		ReplyMarkup: builder.TeacherInlineKeyboardTestSetOptions(testSetID),
	})

	b.DeleteMessages(ctx, &bot.DeleteMessagesParams{
		ChatID:     userID,
		MessageIDs: msgToDelete[userID],
	})
}

func (app *App) TeacherInsertWordsIntoTestSet(ctx context.Context, b *bot.Bot, u *models.Update) {
	if u.CallbackQuery == nil {
		return
	}
	userID := u.CallbackQuery.From.ID
	chatID := userID

	cData := strings.TrimPrefix(u.CallbackQuery.Data, "insert_words_into_")
	testSetID, err := strconv.ParseInt(cData, 10, 64)
	if err != nil {
		log.Fatalf("err parsing test set ID: %v\n", err)
	}

	app.F.Transition(userID, StateInsertWordsIntoTestSet, userID, chatID)
	app.F.Set(userID, "insert_words_into", testSetID)
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: u.CallbackQuery.ID,
	})

	msg, _ := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    userID,
		MessageID: u.CallbackQuery.Message.Message.ID,
		Text:      "So'zlarni kiriting. \nFormat: inglizcha so'z#o'zbekcha so'z",
	})
	msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
}

func (app *App) TeacherFinishWordInserting(ctx context.Context, b *bot.Bot, u *models.Update) {
	fmt.Println("in finish func")
	if u.CallbackQuery == nil {
		return
	}

	cData := strings.TrimPrefix(u.CallbackQuery.Data, "finish_inserting_words_into_")
	testSetID, err := strconv.ParseInt(cData, 10, 64)
	if err != nil {
		log.Fatalf("err getting test set id in TeacherFinishWordInserting: %v\n", err)
	}

	app.F.Transition(u.CallbackQuery.From.ID, StateDefault)
	u.CallbackQuery.Data = fmt.Sprintf("teacher_test_set_%d", testSetID)
	app.TeacherTestSetOptions(ctx, b, u)
	fmt.Println("finish func")
}

func (app *App) TeacherDeleteTestSet(ctx context.Context, b *bot.Bot, u *models.Update) {
	if u.CallbackQuery == nil {
		return
	}

	cData := strings.TrimPrefix(u.CallbackQuery.Data, "delete_test_set_")
	testSetID, err := strconv.ParseInt(cData, 10, 64)
	if err != nil {
		log.Fatalf("err parsing test set id: %v\n", err)
	}

	err = app.Store.DeleteTestSet(ctx, app.Store.Pool, testSetID)
	if err != nil {
		log.Fatalf("err deleting test set: %v\n", err)
	}

	u.CallbackQuery.Data = "test_sets_page_0"
	app.TeacherTestSetsList(ctx, b, u)
}

func (app *App) TeacherCreateClass(ctx context.Context, b *bot.Bot, u *models.Update) {
	if u.CallbackQuery == nil {
		return
	}
	userID := u.CallbackQuery.From.ID
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: u.CallbackQuery.ID,
	})
	app.F.Transition(userID, StateAskClassName, userID)
	msgToDelete[userID] = append(msgToDelete[userID], u.CallbackQuery.Message.Message.ID)
}

func (app *App) TeacherClassesList(ctx context.Context, b *bot.Bot, u *models.Update) {
	if u.CallbackQuery == nil {
		return
	}
	userID := u.CallbackQuery.From.ID
	chatID := userID
	offset, err := strconv.ParseInt(strings.TrimPrefix(u.CallbackQuery.Data, "teacher_classes_page_"), 10, 32)
	if err != nil {
		log.Fatalf("err parsing offset: %v\n", err)
	}

	classesCount, err := app.Store.ClassesCount(ctx, app.Store.Pool, userID)
	if err != nil {
		log.Fatalf("err counting classes: %v\n", err)
	}
	if classesCount == 0 {
		app.B.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: u.CallbackQuery.ID,
			Text:            "Sinflar mavjud emas ❌",
			ShowAlert:       true,
		})

		msgToDelete[userID] = append(msgToDelete[userID], u.CallbackQuery.Message.Message.ID)
		app.DashBoard(ctx, b, u, userID)
		return
	}
	log.Printf("next button offset: %d\n", offset)
	arg := db.ListClassesByTeacherIDParams{
		TeacherID: userID,
		Limit:     5,
		Offset:    int32(offset),
	}
	classesList, err := app.Store.ListClassesByTeacherID(ctx, app.Store.Pool, arg)
	if err != nil {
		log.Fatalf("err getting test sets: %v\n", err)
	}

	isLastPage := false
	if math.Ceil(float64(classesCount)/float64(arg.Limit))-1 == float64(offset) {
		isLastPage = true
	}
	var ids []int64
	var classNames string
	for id, class := range classesList {
		ids = append(ids, class.ClassID)
		classNames = fmt.Sprintf("%s\n%d. %s", classNames, int(offset)+id+1, class.ClassName)
	}
	fmt.Printf("got classes: %v\n", ids)
	options := builder.InlinePaginatorParams{
		ItemCallback:      "teacher_class",
		NavigatorCallback: "teacher_classes_page",
	}
	msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        classNames,
		ReplyMarkup: builder.NewInlinePaginator(ids, arg.Offset, isLastPage, options),
	})
	if err != nil {
		log.Printf("err sending classes list: %v\n", err)
	}
	b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    chatID,
		MessageID: u.CallbackQuery.Message.Message.ID,
	})
	msgToDelete[userID] = append(msgToDelete[userID], msg.ID)

}

func (app *App) TeacherClassOptions(ctx context.Context, b *bot.Bot, u *models.Update) {
	if u.CallbackQuery == nil {
		return
	}

	userID := u.CallbackQuery.From.ID
	classID, err := strconv.ParseInt(strings.TrimPrefix(u.CallbackQuery.Data, "teacher_class_"), 10, 64)
	if err != nil {
		log.Fatalf("err parsing class id:  %v\n", err)
	}

	class, err := app.Store.GetClassByID(ctx, app.Store.Pool, classID)
	if err != nil {
		log.Fatalf("err getting class by id: %v\n", err)
	}

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: u.CallbackQuery.ID,
	})
	b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    userID,
		MessageID: u.CallbackQuery.Message.Message.ID,
	})

	msg, _ := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userID,
		Text:        class.ClassName,
		ReplyMarkup: builder.TeacherInlineKeyboardClassOptions(classID),
	})

	// msgToDelete[userID] = append(msgToDelete[userID], u.CallbackQuery.Message.Message.ID)
	msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
}

func (app *App) AddTestSetToClass(ctx context.Context, b *bot.Bot, u *models.Update) {
	if u.Message == nil {
		return
	}

	args := strings.Split(u.Message.Text, " ")
	if len(args) != 3 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: u.Message.From.ID,
			Text:   "Xato buyruq. Namuna: /acct test-id class-id",
		})
		return
	}
	test_id, errTestID := strconv.ParseInt(args[1], 10, 64)
	class_id, errClassID := strconv.ParseInt(args[2], 10, 64)
	if errTestID != nil || errClassID != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: u.Message.From.ID,
			Text:   "Ma'lumotlarni to'g'ri kiriting.",
		})
	}

	var isExist bool
	app.Store.Pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM test_sets WHERE test_set_id = $1)", test_id).Scan(&isExist)
	if !isExist {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: u.Message.From.ID,
			Text:   "Bunday test mavjud emas.",
		})
		return
	}

	app.Store.Pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM classes WHERE class_id = $1)", class_id).Scan(&isExist)
	if !isExist {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: u.Message.From.ID,
			Text:   "Bunday sinf mavjud emas.",
		})
		return
	}

	err := app.Store.AddTestSetToClass(ctx, app.Store.Pool, db.AddTestSetToClassParams{
		ClassID:   class_id,
		TestSetID: test_id,
	})

	if err != nil {
		log.Fatalf("err inserting test to class: %v\n", err)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: u.Message.From.ID,
		Text:   "Qo'shildi",
	})
}

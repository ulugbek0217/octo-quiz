package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ulugbek0217/octo-quiz/builder"
	db "github.com/ulugbek0217/octo-quiz/db/sqlc"
	"github.com/ulugbek0217/octo-quiz/util"
)

var msgToDelete = map[int64][]int{}

func (app *App) MainHandler(ctx context.Context, b *bot.Bot, u *models.Update) {
	if u.Message == nil && u.CallbackQuery == nil {
		return
	}

	userID := int64(0)
	chatID := int64(0)

	// Determine userID and chatID based on update type
	if u.Message != nil {
		userID = u.Message.From.ID
		chatID = u.Message.Chat.ID
	} else if u.CallbackQuery != nil {
		userID = u.CallbackQuery.From.ID
		chatID = userID
	}

	state := app.F.Current(userID)
	log.Printf("Current state for user %d: %s", userID, state)

	switch state {
	case StateDefault:
		app.DashBoard(ctx, b, u, chatID)
	case StateAskName:
		app.askName(ctx, u, userID, chatID)
	case StateAskUsername:
		app.askUsername(ctx, u, userID, chatID)
	case StateAskRole:
		app.askRole(ctx, u, userID, chatID)
	case StateAskPhone:
		app.askPhone(ctx, u, userID, chatID)
	case StateFinishRegistration:
		app.registrationFinish(ctx, u, userID, chatID)
	case StateAskTestSetName:
		app.testSetName(ctx, u, userID, chatID)
	case StateAskTestSetType:
		app.testSetType(ctx, u, userID, chatID)
	case StateAskTestSetTimeLimitAndFinish:
		app.testSetTimeLimitAndFinish(ctx, u, userID, chatID)
	case StateInsertWordsIntoTestSet:
		app.insertWordsIntoTestSet(ctx, u, userID, chatID)
	default:
		log.Printf("Unexpected state: %s", state)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Xatolik yuz berdi. Iltimos, qaytadan boshlang.",
		})
		app.F.Transition(userID, StateDefault, userID, chatID)
	}

}

func (app *App) askName(ctx context.Context, u *models.Update, args ...any) {
	userID := args[0].(int64)
	chatID := args[1].(int64)
	if u.Message == nil || u.Message.Text == "" {
		msg, _ := app.B.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Iltimos, ism va familiyangizni kiriting.",
		})
		msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
		return
	}
	fullName := u.Message.Text
	if len(fullName) < 5 {
		msg, _ := app.B.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Ism va familiyangizni to'liq kiriting (kamida 5 belgi).",
		})
		msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
		msgToDelete[userID] = append(msgToDelete[userID], u.Message.ID)
		return
	}
	log.Print("Got name: ", fullName)
	app.F.Set(userID, "fullName", fullName)
	app.F.Transition(userID, StateAskUsername, userID, chatID)
	msgToDelete[userID] = append(msgToDelete[userID], u.Message.ID)
}

func (app *App) askUsername(ctx context.Context, u *models.Update, args ...any) {
	userID := args[0].(int64)
	chatID := args[1].(int64)
	// Handle callback query for "generate_username"
	if u.CallbackQuery != nil && u.CallbackQuery.Data == "generate_username" {
		generatedUsername := fmt.Sprintf("user%d", userID) // Implement your own username generation logic
		app.F.Set(userID, "username", generatedUsername)
		app.B.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: u.CallbackQuery.ID,
		})
		msgToDelete[userID] = append(msgToDelete[userID], u.CallbackQuery.Message.Message.ID)
		app.F.Transition(userID, StateAskRole, userID, chatID)
		return
	}

	// Handle text input for username
	if u.Message == nil || u.Message.Text == "" {
		msg, _ := app.B.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Iltimos, username kiriting yoki 'Generatsiya qilish' tugmasini bosing.",
		})
		msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
		return
	}

	username := u.Message.Text
	// Validate username (lowercase letters and numbers only)
	if !util.IsValidUsername(username) {
		msg, _ := app.B.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Username faqat kichik ingliz harflari va raqamlardan iborat bo'lishi kerak. Namuna: abu6500",
		})
		msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
		msgToDelete[userID] = append(msgToDelete[userID], u.Message.ID)
		return
	}
	log.Print("Got username: ", username)
	app.F.Set(userID, "username", username)
	app.F.Transition(userID, StateAskRole, userID, chatID)
	msgToDelete[userID] = append(msgToDelete[userID], u.Message.ID)
}

func (app *App) askRole(ctx context.Context, u *models.Update, args ...any) {
	userID := args[0].(int64)
	chatID := args[1].(int64)
	if u.CallbackQuery == nil || !strings.HasPrefix(u.CallbackQuery.Data, "role_") {
		msg, _ := app.B.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Iltimos, tugmalardan birini tanlang.",
		})
		msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
		return
	}
	role := strings.Split(u.CallbackQuery.Data, "_")[1]
	log.Print("Got role: ", role)
	app.F.Set(userID, "role", role)
	app.F.Transition(userID, StateAskPhone, userID, chatID)
	app.B.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: u.CallbackQuery.ID,
	})
	msgToDelete[userID] = append(msgToDelete[userID], u.CallbackQuery.Message.Message.ID)
}

func (app *App) askPhone(ctx context.Context, u *models.Update, args ...any) {
	userID := args[0].(int64)
	chatID := args[1].(int64)

	if u.Message == nil || u.Message.Contact == nil {
		msg, _ := app.B.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Iltimos, telefon raqamingizni yuborish uchun tugmadan foydalaning.",
		})
		msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
		msgToDelete[userID] = append(msgToDelete[userID], u.Message.ID)
		return
	}
	phone := u.Message.Contact.PhoneNumber
	// phone = strings.TrimPrefix(phone, "+")

	log.Print("Got phone: ", phone)
	app.F.Set(userID, "phone", phone)
	app.F.Transition(userID, StateFinishRegistration, userID, chatID)
	msgToDelete[userID] = append(msgToDelete[userID], u.Message.ID)
}

func (app *App) registrationFinish(ctx context.Context, u *models.Update, args ...any) {
	userID := args[0].(int64)
	chatID := args[1].(int64)
	if u.CallbackQuery == nil {
		msg, _ := app.B.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Tugmalardan foydalaning!",
		})
		msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
		return
	}
	if u.CallbackQuery.Data == "register_again" {
		app.F.Set(userID, "fullName", "")
		app.F.Set(userID, "username", "")
		app.F.Set(userID, "role", "")
		app.F.Set(userID, "phone", 0)

		app.B.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: u.CallbackQuery.ID,
		})

		msgToDelete[userID] = append(msgToDelete[userID], u.CallbackQuery.Message.Message.ID)
		app.B.DeleteMessages(ctx, &bot.DeleteMessagesParams{
			ChatID:     chatID,
			MessageIDs: msgToDelete[userID],
		})
		delete(msgToDelete, userID)

		app.F.Transition(userID, StateAskName, userID, chatID)
		return
	}
	if u.CallbackQuery.Data == "registration_done" {
		app.B.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: u.CallbackQuery.ID,
			Text:            "Muvaffaqiyatli ro'yxatdan o'tdingiz ✅",
		})
		app.F.Reset(userID)
		msgToDelete[userID] = append(msgToDelete[userID], u.CallbackQuery.Message.Message.ID)

		telegram_username := u.CallbackQuery.From.Username
		fullName, _ := app.F.Get(userID, "fullName")
		username, _ := app.F.Get(userID, "username")
		role, _ := app.F.Get(userID, "role")
		phone, _ := app.F.Get(userID, "phone")
		arg := db.CreateUserParams{
			UserID: userID,
			TelegramUsername: pgtype.Text{
				String: telegram_username,
				Valid:  true,
			},
			FullName: fullName.(string),
			Username: username.(string),
			Role:     role.(string),
			Phone:    phone.(string),
		}
		_, err := app.Store.CreateUser(ctx, app.Store.Pool, arg)
		if err != nil {
			log.Fatalf("err creating user: %v\n", err)
		}
	}
	log.Println(msgToDelete)
	app.B.DeleteMessages(ctx, &bot.DeleteMessagesParams{
		ChatID:     chatID,
		MessageIDs: msgToDelete[userID],
	})
	delete(msgToDelete, userID)
	app.DashBoard(ctx, app.B, u, chatID)
}

func (app *App) testSetName(ctx context.Context, u *models.Update, args ...any) {
	userID := args[0].(int64)
	chatID := args[1].(int64)

	if u.Message.Text == "" {
		msg, _ := app.B.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Iltimos test to'plami nomini kiriting.",
		})
		msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
		return
	}
	name := u.Message.Text
	log.Printf("New test set name: %s\n", name)
	app.F.Set(userID, "new_test_set_name", name)
	app.F.Set(userID, "new_test_set_creator_id", userID)
	msgToDelete[userID] = append(msgToDelete[userID], u.Message.ID)
	app.F.Transition(userID, StateAskTestSetType, userID, chatID)
}

func (app *App) testSetType(ctx context.Context, u *models.Update, args ...any) {
	userID := args[0].(int64)
	chatID := args[1].(int64)

	if u.CallbackQuery == nil || u.CallbackQuery.Data == "" {
		msg, _ := app.B.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Iltimos test to'plami uchun turni tanlang.",
		})
		msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
		return
	}
	app.B.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: u.CallbackQuery.ID,
	})
	testType := strings.TrimPrefix(u.CallbackQuery.Data, "test_set_type_")
	log.Printf("New test set type: %s\n", testType)
	app.F.Set(userID, "new_test_set_type", testType)
	msgToDelete[userID] = append(msgToDelete[userID], u.CallbackQuery.Message.Message.ID)
	app.F.Transition(userID, StateAskTestSetTimeLimitAndFinish, userID, chatID)
}

func (app *App) testSetTimeLimitAndFinish(ctx context.Context, u *models.Update, args ...any) {
	userID := args[0].(int64)
	chatID := args[1].(int64)

	if u.Message == nil || u.Message.Text == "" {
		msg, _ := app.B.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Iltimos test to'plami uchun vaqt limitini kiriting.",
		})
		msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
		return
	}
	timeLimit := u.Message.Text
	log.Printf("New test set time limit: %s\n", timeLimit)

	testSetName, _ := app.F.Get(userID, "new_test_set_name")
	var isPublic bool
	if testType, ok := app.F.Get(userID, "new_test_set_type"); ok {
		switch testType {
		case "public":
			isPublic = true
		case "private":
			isPublic = false
		default:
			log.Fatal("err test type(ispublic)")
		}
	}
	testSetCreatorID, _ := app.F.Get(userID, "new_test_set_creator_id")
	timeLimitInt, err := strconv.ParseInt(timeLimit, 10, 32)
	if err != nil {
		log.Fatalf("err parsing time limit: %v\n", err)
	}

	arg := db.CreateTestSetParams{
		TestSetName: testSetName.(string),
		CreatorID:   testSetCreatorID.(int64),
		IsPublic:    isPublic,
		TimeLimit: pgtype.Int4{
			Int32: int32(timeLimitInt),
			Valid: true,
		},
	}

	_, err = app.Store.CreateTestSet(ctx, app.Store.Pool, arg)
	if err != nil {
		log.Fatalf("err writing test set to db: %v\n", err)
	}

	// app.F.Set(userID, "new_test_set_time_limit", timeLimit)
	msgToDelete[userID] = append(msgToDelete[userID], u.Message.ID)
	app.F.Transition(userID, StateDefault, userID, chatID)
	app.DashBoard(ctx, app.B, u, chatID)
	app.B.DeleteMessages(ctx, &bot.DeleteMessagesParams{
		ChatID:     chatID,
		MessageIDs: msgToDelete[userID],
	})
	delete(msgToDelete, userID)
}

func (app *App) insertWordsIntoTestSet(ctx context.Context, u *models.Update, args ...any) {
	if u.Message == nil {
		return
	}

	userID := args[0].(int64)
	chatID := args[1].(int64)

	testSetID, ok := app.F.Get(userID, "insert_words_into")
	if !ok {
		app.B.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "error getting test id to insert into",
		})
		return
	}
	userMessage := u.Message.Text
	wordLines := strings.Split(userMessage, "\n")
	// var wordsMap map[string]string

	for _, line := range wordLines {
		words := strings.Split(line, "#")
		app.WG.Add(1)
		go func() {
			defer app.WG.Done()
			arg := db.InsertWordsParams{
				TestSetID:   testSetID.(int64),
				EnglishWord: words[0],
				UzbekWord:   words[1],
			}

			_, err := app.Store.InsertWords(ctx, app.Store.Pool, arg)
			if err != nil {
				log.Fatalf("error inserting words [%s, %s]: %v\n", words[0], words[1], err)
			}
		}()
		// wordsMap[words[0]] = words[1]
	}
	app.WG.Wait()

	// keyboard, err := builder.NewInlineKeyboardBuilderFromJson(builder.KeyboardFinishOrInsertWordsButtons)
	// if err != nil {
	// 	log.Fatalf("err building keyboard: %v\n", err)
	// }
	// markup := keyboard.Build()
	app.B.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        "So'zlar saqlandi ✅",
		ReplyMarkup: builder.TeacherInlineKeyboardInsertWordsOrFinish(testSetID.(int64)),
	})
}

// func (app *App) testSetCreatingFinish(ctx context.Context, u *models.Update, args ...any) {
// 	userID := args[0].(int64)
// 	chatID := args[1].(int64)

// 	if u.Message == nil || u.Message.Text == "" {
// 		msg, _ := app.B.SendMessage(ctx, &bot.SendMessageParams{
// 			ChatID: chatID,
// 			Text:   "Iltimos test to'plami uchun vaqt limitini kiriting.",
// 		})
// 		msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
// 		return
// 	}
// 	name := u.Message.Text
// 	log.Printf("New test set name: %s\n", name)
// 	app.F.Set(userID, "new_test_set_name", name)
// 	msgToDelete[userID] = append(msgToDelete[userID], u.Message.ID)
// 	app.F.Transition(userID, StateFinishTestSetCreating, userID, chatID)
// }

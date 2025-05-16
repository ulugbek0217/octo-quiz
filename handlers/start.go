package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/jackc/pgx/v5"
	"github.com/ulugbek0217/octo-quiz/builder"
	db "github.com/ulugbek0217/octo-quiz/db/sqlc"
)

func (app *App) Start(ctx context.Context, b *bot.Bot, u *models.Update) {
	// msgToDelete[u.Message.From.ID] = append(msgToDelete[u.Message.From.ID], u.Message.ID)
	_, newUser := app.isNewUser(u.Message.From.ID)
	fmt.Println(newUser)
	if newUser {
		// Transition to StateAskName and trigger the callback
		app.F.Transition(u.Message.From.ID, StateAskName, u.Message.From.ID, u.Message.Chat.ID)
		return
	}

	app.DashBoard(ctx, b, u, u.Message.From.ID)

}

func (app *App) DashBoard(ctx context.Context, b *bot.Bot, u *models.Update, args ...any) {

	var chatID int64
	if len(args) != 0 {
		chatID = args[0].(int64)
	} else {
		chatID = u.CallbackQuery.From.ID
	}

	user, _ := app.isNewUser(chatID)
	if u.CallbackQuery != nil {
		msgToDelete[user.UserID] = append(msgToDelete[user.UserID], u.CallbackQuery.Message.Message.ID)
	}
	b.DeleteMessages(ctx, &bot.DeleteMessagesParams{
		ChatID:     chatID,
		MessageIDs: msgToDelete[user.UserID],
	})

	kbd, err := builder.NewInlineKeyboardBuilder(builder.KeyboardStudentMainMenuInlineButtons)
	if user.Role == "teacher" {
		kbd, err = builder.NewInlineKeyboardBuilder(builder.KeyboardTeacherMainMenuInlineButtons)
	}
	if err != nil {
		log.Fatal(err)
	}
	markup := kbd.Build()
	app.B.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        "Ishchi stol\n\nDashboard",
		ReplyMarkup: markup,
	})
}

func (app *App) isNewUser(userID int64) (db.User, bool) {
	user, err := app.Store.GetUser(context.Background(), app.Store.Pool, userID)
	if err == pgx.ErrNoRows {
		return db.User{}, true
	} else if err == nil {
		return user, false
	}
	log.Fatalf("err getting user from db: %v\n", err)
	return db.User{}, false
}

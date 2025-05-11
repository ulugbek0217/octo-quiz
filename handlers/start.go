package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/fsm"
	"github.com/jackc/pgx/v5"
)

func (app *App) Start(ctx context.Context, b *bot.Bot, u *models.Update) {
	msgToDelete[u.Message.From.ID] = append(msgToDelete[u.Message.From.ID], u.Message.ID)
	if _, newUser := app.isNewUser(u.Message.From.ID); newUser {
		// Transition to StateAskName and trigger the callback
		app.F.Transition(u.Message.From.ID, StateAskName, u.Message.From.ID, u.Message.Chat.ID)
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: u.Message.Chat.ID,
		Text:   "Octo Quiz'ga xush kelibsiz!\nWelcome to Octo Quiz!",
	})
}

func (app *App) CallbackName(f *fsm.FSM, args ...any) {
	userID := args[0].(int64)
	chatID := args[1].(int64)
	msg, _ := app.B.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Octo Quiz'ga xush kelibsiz.\n\nRo'yxatdan o'tish uchun ism va familiyangizni kiriting.",
	})
	msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
}

func (app *App) CallbackUsername(f *fsm.FSM, args ...any) {
	userID := args[0].(int64)
	chatID := args[1].(int64)
	msg, _ := app.B.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "O'zingiz uchun username tanlang. Username kichik ingliz harflari va raqamlardan iborat bo'lishi kerak.\nNamuna: abu6500.",
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{
						Text:         "Generatsiya qilish",
						CallbackData: "generate_username",
					},
				},
			},
		},
	})
	msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
}

func (app *App) CallbackRole(f *fsm.FSM, args ...any) {
	userID := args[0].(int64)
	chatID := args[1].(int64)
	msg, _ := app.B.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Siz o'qituvchimisiz yoki o'quvchi?",
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{
						Text:         "O'quvchi",
						CallbackData: "role_student",
					},
					{
						Text:         "O'qituvchi",
						CallbackData: "role_teacher",
					},
				},
			},
		},
	})
	msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
}

func (app *App) CallbackPhone(f *fsm.FSM, args ...any) {
	userID := args[0].(int64)
	chatID := args[1].(int64)
	msg, _ := app.B.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Telefon raqamingizni yuboring.",
		ReplyMarkup: models.ReplyKeyboardMarkup{
			Keyboard: [][]models.KeyboardButton{
				{
					{
						Text:           "Telefon raqamingizni yuboring.",
						RequestContact: true,
					},
				},
			},
			ResizeKeyboard:  true,
			OneTimeKeyboard: true,
		},
	})
	msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
}

func (app *App) CallbackFinish(f *fsm.FSM, args ...any) {
	userID := args[0].(int64)
	chatID := args[1].(int64)

	fullName, _ := app.F.Get(userID, "fullName")
	username, _ := app.F.Get(userID, "username")
	role, _ := app.F.Get(userID, "role")
	phone, _ := app.F.Get(userID, "phone")

	userData := fmt.Sprintf("<b>Ma'lumotlar</b>\nIsm, familiya: %s\nUsername: %s\nRol: %s\nTelefon raqam: %s", fullName, username, role, phone)
	msg, _ := app.B.SendMessage(context.Background(), &bot.SendMessageParams{
		Text:      userData,
		ChatID:    chatID,
		ParseMode: models.ParseModeHTML,
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{
						Text:         "Qayta ro'yxatdan o'tish",
						CallbackData: "register_again",
					},
					{
						Text:         "Tasdiqlash",
						CallbackData: "registration_done",
					},
				},
			},
		},
	})
	msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
}
func (app *App) isNewUser(userID int64) (string, bool) {
	user, err := app.Store.GetUser(context.Background(), app.Store.Pool, userID)
	if err == pgx.ErrNoRows {
		return "", true
	} else if err == nil {
		return user.Role, false
	}
	log.Fatalf("err getting user from db: %v\n", err)
	return "", false
}

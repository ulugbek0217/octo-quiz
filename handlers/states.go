package handlers

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/fsm"
)

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

	StateInsertWordsIntoTestSet          fsm.StateID = "insert_words_into_test_set"
	StateFinishInsertingWordsIntoTestSet fsm.StateID = "finish_inserting_words_into_test_set"
)

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

func (app *App) CallbackTestSetName(f *fsm.FSM, args ...any) {
	userID := args[0].(int64)
	chatID := args[1].(int64)

	msg, _ := app.B.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Yangi test to'plami nomini kiriting.",
	})

	msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
}

func (app *App) CallbackTestSetType(f *fsm.FSM, args ...any) {
	userID := args[0].(int64)
	chatID := args[1].(int64)

	msg, _ := app.B.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Ushbu test to'plami turini tanlang",
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{
						Text:         "Ommaviy",
						CallbackData: "test_set_type_public",
					},
					{
						Text:         "Yashirin",
						CallbackData: "test_set_type_private",
					},
				},
			},
		},
	})

	msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
}

func (app *App) CallbackTestSetTimeLimit(f *fsm.FSM, args ...any) {
	userID := args[0].(int64)
	chatID := args[1].(int64)

	msg, _ := app.B.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Ushbu test to'plami uchun vaqt limitini belgilang.\nO'lchov birligi - soniya. Namuna: 120, ya'ni 120 soniya.\nDefault: 0",
	})

	msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
}

func (app *App) CallbackFinishTestSetCreating(f *fsm.FSM, args ...any) {
	userID := args[0].(int64)
	chatID := args[1].(int64)

	msg, _ := app.B.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Test to'plami yaratildi.",
	})

	msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
}

func (app *App) CallbackWaitForWords(f *fsm.FSM, args ...any) {
	userID := args[0].(int64)
	chatID := args[1].(int64)

	msg, _ := app.B.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Ushbu test to'plami uchun so'zlarni kiriting.\nFormat: en#uz",
	})

	msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
}

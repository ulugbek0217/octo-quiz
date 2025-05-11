package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
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
		b.SendMessage(ctx, &bot.SendMessageParams{
			Text:   "Default state",
			ChatID: chatID,
		})
	case StateAskName:
		if u.Message == nil || u.Message.Text == "" {
			msg, _ := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Iltimos, ism va familiyangizni kiriting.",
			})
			msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
			return
		}
		fullName := u.Message.Text
		if len(fullName) < 5 {
			msg, _ := b.SendMessage(ctx, &bot.SendMessageParams{
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
	case StateAskUsername:
		// Handle callback query for "generate_username"
		if u.CallbackQuery != nil && u.CallbackQuery.Data == "generate_username" {
			generatedUsername := fmt.Sprintf("user%d", userID) // Implement your own username generation logic
			app.F.Set(userID, "username", generatedUsername)
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: u.CallbackQuery.ID,
			})
			msgToDelete[userID] = append(msgToDelete[userID], u.CallbackQuery.Message.Message.ID)
			app.F.Transition(userID, StateAskRole, userID, chatID)
			return
		}

		// Handle text input for username
		if u.Message == nil || u.Message.Text == "" {
			msg, _ := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Iltimos, username kiriting yoki 'Generatsiya qilish' tugmasini bosing.",
			})
			msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
			return
		}

		username := u.Message.Text
		// Validate username (lowercase letters and numbers only)
		if !util.IsValidUsername(username) {
			msg, _ := b.SendMessage(ctx, &bot.SendMessageParams{
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
	case StateAskRole:
		if u.CallbackQuery == nil || !strings.HasPrefix(u.CallbackQuery.Data, "role_") {
			msg, _ := b.SendMessage(ctx, &bot.SendMessageParams{
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
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: u.CallbackQuery.ID,
		})
		msgToDelete[userID] = append(msgToDelete[userID], u.CallbackQuery.Message.Message.ID)
	case StateAskPhone:
		if u.Message == nil || u.Message.Contact == nil {
			msg, _ := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Iltimos, telefon raqamingizni yuborish uchun tugmadan foydalaning.",
			})
			msgToDelete[userID] = append(msgToDelete[userID], msg.ID)
			msgToDelete[userID] = append(msgToDelete[userID], u.Message.ID)
			return
		}
		phone := u.Message.Contact.PhoneNumber
		phone = strings.TrimPrefix(phone, "+")

		log.Print("Got phone: ", phone)
		app.F.Set(userID, "phone", phone)
		app.F.Transition(userID, StateFinish, userID, chatID)
		msgToDelete[userID] = append(msgToDelete[userID], u.Message.ID)
	case StateFinish:
		if u.CallbackQuery == nil {
			msg, _ := b.SendMessage(ctx, &bot.SendMessageParams{
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

			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: u.CallbackQuery.ID,
			})

			msgToDelete[userID] = append(msgToDelete[userID], u.CallbackQuery.Message.Message.ID)
			b.DeleteMessages(ctx, &bot.DeleteMessagesParams{
				ChatID:     chatID,
				MessageIDs: msgToDelete[userID],
			})
			delete(msgToDelete, userID)

			app.F.Transition(userID, StateAskName, userID, chatID)
			return
		}
		if u.CallbackQuery.Data == "registration_done" {
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: u.CallbackQuery.ID,
				Text:            "Muvaffaqiyatli ro'yxatdan o'tdingiz ✅",
			})
			app.F.Reset(userID)
			msgToDelete[userID] = append(msgToDelete[userID], u.CallbackQuery.Message.Message.ID)
		}
		log.Println(msgToDelete)
		b.DeleteMessages(ctx, &bot.DeleteMessagesParams{
			ChatID:     chatID,
			MessageIDs: msgToDelete[userID],
		})
		delete(msgToDelete, userID)
	default:
		log.Printf("Unexpected state: %s", state)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Xatolik yuz berdi. Iltimos, qaytadan boshlang.",
		})
		app.F.Transition(userID, StateDefault, userID, chatID)
	}

}

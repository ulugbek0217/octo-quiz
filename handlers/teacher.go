package handlers

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
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
	userID := u.CallbackQuery.From.ID
	chatID := userID

	testSetsList, err := app.Store.GetTestSetByCreatorID(ctx, app.Store.Pool, userID)

}

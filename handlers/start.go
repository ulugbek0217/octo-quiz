package handlers

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/jackc/pgx/v5"
)

func (app *App) Start(ctx context.Context, b *bot.Bot, u *models.Update) {
	if app.isNewUser(u.Message.From.ID) {
		log.Print("New User")
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: u.Message.Chat.ID,
		Text:   "Welcome to Octo Quiz!",
	})

}

func (app *App) Register(ctx context.Context, b *bot.Bot, u *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: u.Message.Chat.ID,
		Text:   "Ism va familiyangizni yozing.",
	})
}

func (app *App) isNewUser(userID int64) bool {
	_, err := app.Store.GetUser(context.Background(), app.Store.Pool, userID)
	return err == pgx.ErrNoRows
}

package handlers

import (
	"github.com/stefanitsky/gachinator-telegram-bot/db"
	"github.com/stefanitsky/gachinator-telegram-bot/translations"
	tele "gopkg.in/telebot.v3"
)

// Start handler initiates first user communication with a bot
func StartHandler(c tele.Context) error {
	user := c.Sender()
	db := c.Get("db").(db.DB)

	langCode, err := db.GetUserLangCode(user.ID)
	if err != nil {
		return err
	}

	return c.Send(translations.GetMessage(translations.Welcome, translations.LangCode(langCode)))
}

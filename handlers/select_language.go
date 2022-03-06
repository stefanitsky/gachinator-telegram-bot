package handlers

import (
	"github.com/stefanitsky/gachinator-telegram-bot/db"
	"github.com/stefanitsky/gachinator-telegram-bot/menu"
	"github.com/stefanitsky/gachinator-telegram-bot/translations"
	tele "gopkg.in/telebot.v3"
)

// Select language handler returns language menu
func SelectLangHandler(c tele.Context) error {
	user := c.Sender()

	db := c.Get("db").(db.DB)

	langCode, err := db.GetUserLangCode(user.ID)
	if err != nil {
		langCode = translations.DefaultLangCode
	}

	return c.Send(translations.GetMessage(translations.SelectLangMenu, translations.LangCode(langCode)), menu.LangMenu)
}

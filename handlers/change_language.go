package handlers

import (
	"github.com/stefanitsky/gachinator-telegram-bot/db"
	"github.com/stefanitsky/gachinator-telegram-bot/translations"
	tele "gopkg.in/telebot.v3"
)

// Language button handler changes user default language
func ChangeLangHandler(c tele.Context) error {
	user := c.Sender()

	db := c.Get("db").(db.DB)
	lc := c.Get("lc").(translations.LangCode)

	if err := db.SetUserLangCode(user.ID, lc); err != nil {
		c.Reply(translations.GetErrorMessage(lc))
		return err
	}

	return c.Reply(translations.GetMessage(translations.SelectNewLang, lc))
}

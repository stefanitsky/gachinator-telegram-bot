package handlers

import (
	"github.com/stefanitsky/gachinator"
	"github.com/stefanitsky/gachinator-telegram-bot/db"
	"github.com/stefanitsky/gachinator-telegram-bot/translations"
	tele "gopkg.in/telebot.v3"
)

// Text handler does the main job, it "gachinates" a text
func OnTextHandler(c tele.Context) error {
	user := c.Sender()

	db := c.Get("db").(db.DB)

	langCode, err := db.GetUserLangCode(user.ID)
	if err != nil {
		langCode = translations.DefaultLangCode
	}

	langConfig, err := gachinator.FindLangConfig(string(langCode))
	if err != nil {
		c.Send(err.Error())
		langConfig = &gachinator.RussianConfig
	}

	gachinated := gachinator.Gachinate([]byte(c.Text()), *langConfig)

	return c.Send(string(gachinated))
}

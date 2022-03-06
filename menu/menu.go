package menu

import tele "gopkg.in/telebot.v3"

var (
	LangMenu = &tele.ReplyMarkup{}
	RuBtn    = LangMenu.Text("🇷🇺")
	EnBtn    = LangMenu.Text("🇺🇸")
)

func init() {
	LangMenu.Reply(
		LangMenu.Row(RuBtn),
		LangMenu.Row(EnBtn),
	)
}

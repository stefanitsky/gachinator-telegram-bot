package menu

import tele "gopkg.in/telebot.v3"

var (
	LangMenu = &tele.ReplyMarkup{}
	RuBtn    = LangMenu.Text("π·πΊ")
	EnBtn    = LangMenu.Text("πΊπΈ")
)

func init() {
	LangMenu.Reply(
		LangMenu.Row(RuBtn),
		LangMenu.Row(EnBtn),
	)
}

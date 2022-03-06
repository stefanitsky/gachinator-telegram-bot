package main

import (
	"github.com/stefanitsky/gachinator-telegram-bot/app"
	"github.com/stefanitsky/gachinator-telegram-bot/handlers"
	"github.com/stefanitsky/gachinator-telegram-bot/menu"
	"github.com/stefanitsky/gachinator-telegram-bot/translations"
	tele "gopkg.in/telebot.v3"
)

var (
	application *app.App
)

func main() {
	application = app.NewApp()

	defer application.CloseAllConnections()

	startHandler := handlers.Handler{
		Endpoint:    "/start",
		HandlerFunc: handlers.StartHandler,
		SuccessHandleLog: handlers.SuccessHandleLog{
			Message:     "user invoked /start command",
			EventAction: "cmd-start",
		},
		ExtraContext: handlers.ExtraContext{"db": application.DB},
	}

	selectLangHandler := handlers.Handler{
		Endpoint:    "/lang",
		HandlerFunc: handlers.SelectLangHandler,
		SuccessHandleLog: handlers.SuccessHandleLog{
			Message:     "user invoked /lang command",
			EventAction: "cmd-lang",
		},
		ExtraContext: handlers.ExtraContext{"db": application.DB},
	}

	textHandler := handlers.Handler{
		Endpoint:    tele.OnText,
		HandlerFunc: handlers.OnTextHandler,
		SuccessHandleLog: handlers.SuccessHandleLog{
			Message:     "gachinator is used",
			EventAction: "gachinator",
		},
		ExtraContext: handlers.ExtraContext{"db": application.DB},
	}

	enChangeLangHandler := handlers.Handler{
		Endpoint:    &menu.EnBtn,
		HandlerFunc: handlers.ChangeLangHandler,
		SuccessHandleLog: handlers.SuccessHandleLog{
			Message:     "language changed to english",
			EventAction: "menu-lang-en",
		},
		ExtraContext: handlers.ExtraContext{"db": application.DB, "lc": translations.English},
	}
	ruChangeLangHandler := handlers.Handler{
		Endpoint:    &menu.RuBtn,
		HandlerFunc: handlers.ChangeLangHandler,
		SuccessHandleLog: handlers.SuccessHandleLog{
			Message:     "language changed to russian",
			EventAction: "menu-lang-ru",
		},
		ExtraContext: handlers.ExtraContext{"db": application.DB, "lc": translations.Russian},
	}

	application.AddHandler(startHandler)
	application.AddHandler(selectLangHandler)
	application.AddHandler(textHandler)
	application.AddHandler(enChangeLangHandler)
	application.AddHandler(ruChangeLangHandler)

	application.Start()
}

package app

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/stefanitsky/gachinator-telegram-bot/config"
	"github.com/stefanitsky/gachinator-telegram-bot/db"
	"github.com/stefanitsky/gachinator-telegram-bot/handlers"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

// App is the main app container
type App struct {
	cfg *config.Config
	DB  db.DB
	bot *tele.Bot
	Log *zap.Logger
}

// NewApp creates a new app
func NewApp() *App {
	app := App{}

	app.ConfigureLogging()
	app.LoadConfig()
	app.ConfigureDB()
	app.ConfigureBot()

	return &app
}

// ConfigureLogging configures the main logger
func (app *App) ConfigureLogging() {
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, os.Stdout, zap.DebugLevel)

	app.Log = zap.New(core)

	// Redirect default logger usage into the main logger (telebot, redis, panics etc.)
	zap.RedirectStdLog(app.Log)

	app.Log.Info("logging is configured")
}

// LoadConfig loads app's config
func (app *App) LoadConfig() {
	app.cfg = config.InitAndParse()

	app.Log.Info("config is loaded")
}

// ConfigureDB creates DB client
func (app *App) ConfigureDB() {
	app.DB = db.CreateRedisDB(&redis.Options{
		Addr: fmt.Sprintf("%v:%v", app.cfg.Redis.Host, app.cfg.Redis.Port),
		DB:   app.cfg.Redis.Db,
	})

	app.Log.Info("db is configured")
}

// ConfigureBot configures bot client
func (app *App) ConfigureBot() {
	botSettings := tele.Settings{
		Token: app.cfg.Bot.Token,
		Poller: &tele.Webhook{
			Listen: app.cfg.Bot.Webhook.Listen,
			Endpoint: &tele.WebhookEndpoint{
				PublicURL: app.cfg.Bot.Webhook.URL,
			},
		},
		OnError: app.HandlerError,
	}

	var err error
	app.bot, err = tele.NewBot(botSettings)
	if err != nil {
		app.Log.Fatal(err.Error())
	}

	app.Log.Info("bot is configured")
}

// CloseAllConnections closes all opened connections
func (app *App) CloseAllConnections() {
	app.DB.Close()
	app.bot.Close()
}

// AddHandler adds a new handler func to the bot
func (app *App) AddHandler(h handlers.Handler) {
	handlerWithLog := func(c tele.Context) error {
		user := c.Sender()

		err := h.Handle(c)

		if err != nil {
			app.Log.Info(err.Error(),
				zap.Int64("user.id", user.ID),
				zap.String("user.name", user.Username),
				zap.String("user.full_name", fmt.Sprintf("%v %v", user.FirstName, user.LastName)),
				zap.String("event.action", h.SuccessHandleLog.EventAction),
			)
		}

		app.Log.Info(h.SuccessHandleLog.Message,
			zap.Int64("user.id", user.ID),
			zap.String("user.name", user.Username),
			zap.String("user.full_name", fmt.Sprintf("%v %v", user.FirstName, user.LastName)),
			zap.String("event.action", h.SuccessHandleLog.EventAction),
		)

		return nil
	}

	app.bot.Handle(h.Endpoint, handlerWithLog)
}

// HandlerError replaces the default bot error handling and handles unhandled errors
func (app *App) HandlerError(err error, c tele.Context) {
	user := c.Sender()

	app.Log.Error(err.Error(),
		zap.Int64("user.id", user.ID),
		zap.String("user.name", user.Username),
		zap.String("user.full_name", fmt.Sprintf("%v %v", user.FirstName, user.LastName)),
		zap.String("event.action", "unknown"),
	)
}

// Start starts an app
func (app *App) Start() {
	app.bot.Start()
}

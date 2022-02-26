package main

import (
	"fmt"
	"os"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/stefanitsky/gachinator"
	config "github.com/stefanitsky/gachinator-telegram-bot/config"
	translations "github.com/stefanitsky/gachinator-telegram-bot/translations"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	tele "gopkg.in/telebot.v3"
)

// Telegram reply menu markup
var (
	langMenu = &tele.ReplyMarkup{}
	ruBtn    = langMenu.Text("🇷🇺")
	enBtn    = langMenu.Text("🇺🇸")
)

var (
	cfg         *config.Config
	botSettings tele.Settings
	rdb         *redis.Client
	ctx         = context.Background()
	bot         *tele.Bot
	log         *zap.Logger
)

var (
	defaultLangCode = translations.Russian
)

func init() {
	// Configure logging with zap
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, os.Stdout, zap.DebugLevel)
	log = zap.New(core)

	// Init and load config from environment variables
	cfg = config.InitAndParse()
	log.Info("config is loaded")

	// Create redis client
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%v:%v", cfg.Redis.Host, cfg.Redis.Port),
		DB:   cfg.Redis.Db,
	})
	log.Info("redis client is created")

	// Configure bot settings
	botSettings = tele.Settings{
		Token: cfg.Token,
		// TODO: webhook instead of polling
		Poller: &tele.LongPoller{Timeout: 1 * time.Second},
	}
	log.Info("bot settings are configured")

	// Create bot
	var err error
	bot, err = tele.NewBot(botSettings)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Info("bot is created")

	langMenu.Reply(
		langMenu.Row(ruBtn),
		langMenu.Row(enBtn),
	)
}

func main() {
	defer rdb.Close()
	defer bot.Close()

	bot.Handle("/start", startHandler)
	bot.Handle("/lang", selectLangHandler)
	bot.Handle(tele.OnText, onTextHandler)
	bot.Handle(&ruBtn, ruBtnHandler)
	bot.Handle(&enBtn, enBtnHandler)

	bot.Start()
}

// Start handler initiates first user communication with a bot
func startHandler(c tele.Context) error {
	user := c.Sender()
	defer logHandler("user invoked /start command", user, "command-start")

	langCode := getUserLangCode(user)

	return c.Send(translations.GetMessage(translations.Welcome, langCode))
}

// Select language handler returns language menu
func selectLangHandler(c tele.Context) error {
	user := c.Sender()
	defer logHandler("user invoked /lang command", user, "command-lang")

	langCode := getUserLangCode(user)

	return c.Send(translations.GetMessage(translations.SelectLangMenu, langCode), langMenu)
}

// Text handler does the main job, it "gachinates" a text
func onTextHandler(c tele.Context) error {
	user := c.Sender()
	defer logHandler("gachinator used", user, "gachinator")

	langCode := getUserLangCode(user)

	langConfig, err := gachinator.FindLangConfig(string(langCode))
	if err != nil {
		c.Send(err.Error())
		langConfig = &gachinator.RussianConfig
	}

	gachinated := gachinator.Gachinate([]byte(c.Text()), *langConfig)

	return c.Send(string(gachinated))
}

// Russian language button handler changes user default language to russian
func ruBtnHandler(c tele.Context) error {
	user := c.Sender()
	defer logHandler("language changed", user, "menu-lang-ru")

	if err := setUserLangCode(user, translations.Russian); err != nil {
		c.Reply(translations.GetErrorMessage(translations.Russian))
		return err
	}

	return c.Reply(translations.GetMessage(translations.SelectRussian, translations.Russian))
}

// English language button handler changes user default language to english
func enBtnHandler(c tele.Context) error {
	user := c.Sender()
	defer logHandler("language changed", user, "menu-lang-en")

	if err := setUserLangCode(user, translations.English); err != nil {
		c.Reply(translations.GetErrorMessage(translations.English))
		return err
	}

	return c.Reply(translations.GetMessage(translations.SelectEnglish, translations.English))
}

// Returns current user language code
func getUserLangCode(user *tele.User) translations.LangCode {
	langCode, err := rdb.Get(ctx, fmt.Sprint(user.ID)).Result()

	if err == redis.Nil {
		if err = rdb.Set(ctx, fmt.Sprint(user.ID), user.LanguageCode, 0).Err(); err != nil {
			logHandler(err.Error(), user, "get-user-language-code")
			return defaultLangCode
		} else {
			logHandler(err.Error(), user, "get-user-language-code")
			return defaultLangCode
		}
	} else if err != nil {
		logHandler(err.Error(), user, "get-user-language-code")
		return defaultLangCode
	}

	return translations.LangCode(langCode)
}

// Sets a new user language code
func setUserLangCode(user *tele.User, langCode translations.LangCode) error {
	_, err := rdb.Set(ctx, fmt.Sprint(user.ID), string(langCode), 0).Result()
	return err
}

// Logs info message about used bot handler
func logHandler(msg string, user *tele.User, eventAction string) {
	log.Info(msg,
		zap.Int64("user.id", user.ID),
		zap.String("user.name", user.Username),
		zap.String("user.full_name", fmt.Sprintf("%v %v", user.FirstName, user.LastName)),
		zap.String("event.action", eventAction),
	)
}

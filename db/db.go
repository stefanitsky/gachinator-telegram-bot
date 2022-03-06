package db

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/stefanitsky/gachinator-telegram-bot/translations"
)

type DB interface {
	SetUserLangCode(userId int64, lc translations.LangCode) error
	GetUserLangCode(userId int64) (translations.LangCode, error)
	Close()
}

func CreateRedisDB(options *redis.Options) DB {
	return &Redis{
		ctx:    context.Background(),
		client: redis.NewClient(options),
	}
}

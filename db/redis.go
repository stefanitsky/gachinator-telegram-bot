package db

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/stefanitsky/gachinator-telegram-bot/translations"
)

type Redis struct {
	ctx    context.Context
	client *redis.Client
}

func (r *Redis) SetUserLangCode(userId int64, lc translations.LangCode) error {
	return r.client.Set(r.ctx, fmt.Sprint(userId), string(lc), 0).Err()
}

func (r *Redis) GetUserLangCode(userId int64) (translations.LangCode, error) {
	langCode, err := r.client.Get(r.ctx, fmt.Sprint(userId)).Result()

	if err == redis.Nil {
		return "", fmt.Errorf("no value")
	} else if err != nil {
		return "", fmt.Errorf("unexpected err: %v", err)
	}

	return translations.LangCode(langCode), nil
}

func (r *Redis) Close() {
	r.client.Close()
}

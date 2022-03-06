package config

import (
	"log"

	"github.com/caarlos0/env/v6"
)

type RedisConfig struct {
	Host string `env:"HOST" envDefault:"127.0.0.1"`
	Port int    `env:"PORT" envDefault:"6379"`
	Db   int    `env:"DB" envDefault:"0"`
}

type WebhookConfig struct {
	Listen string `env:"LISTEN" envDefault:"127.0.0.1:8001"`
	URL    string `env:"URL,required"`
}
type BotConfig struct {
	Token   string        `env:"TOKEN,required"`
	Webhook WebhookConfig `envPrefix:"WEBHOOK_"`
}

type Config struct {
	Redis RedisConfig `envPrefix:"REDIS_"`
	Bot   BotConfig   `envPrefix:"BOT_"`
}

func (c *Config) Parse() {
	if err := env.Parse(c); err != nil {
		log.Fatal(err)
	}
}

func InitAndParse() *Config {
	cfg := Config{}
	cfg.Parse()
	return &cfg
}

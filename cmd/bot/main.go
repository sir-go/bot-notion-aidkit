package main

import (
	"flag"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	zlog "github.com/rs/zerolog/log"
)

func main() {
	fCfgPath := flag.String("c", "conf.toml",
		"path to conf file")
	flag.Parse()

	zlog.Info().Msg("init config")
	cfg, err := LoadConfig(*fCfgPath)
	if err != nil {
		panic(err)
	}

	zlog.Info().Msg("init tg bot")
	bot, err := tgbotapi.NewBotAPI(cfg.TgAPI.Token)
	if err != nil {
		panic(err)
	}

	//bot.Debug = true
	zlog.Info().Msg("init warehouse")
	var wh = NewWarehouse(cfg.Wh.UpdateInterval.Duration)

	if err = wh.update(cfg.NotionAPI, true); err != nil {
		zlog.Error().Err(err).Msg("get data")
		panic(err)
	}

	zlog.Info().Str("account", bot.Self.UserName)

	upd := tgbotapi.NewUpdate(0)
	upd.Timeout = 10
	updatesChan, err := bot.GetUpdatesChan(upd)

	if err != nil {
		zlog.Error().Err(err).Msg("get updates channel")
		panic(err)
	}

	for upd := range updatesChan {
		if err = UpdatesProcessing(bot, upd, wh, cfg.NotionAPI); err != nil {
			zlog.Warn().Err(err).Msg("")
		}
	}
}

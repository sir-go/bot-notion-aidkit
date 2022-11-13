package main

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	zlog "github.com/rs/zerolog/log"
)

func SendAnswer(bot *tgbotapi.BotAPI, upd tgbotapi.Update,
	msg string, mode ...string) {
	msgCfg := tgbotapi.NewMessage(upd.Message.Chat.ID, msg)
	if len(mode) > 0 {
		msgCfg.ParseMode = mode[0]
	}
	msgCfg.ReplyToMessageID = upd.Message.MessageID
	_, err := bot.Send(msgCfg)
	if err != nil {
		zlog.Error().Err(err).Msg("send bot answer")
	}
}

func UpdatesProcessing(bot *tgbotapi.BotAPI, upd tgbotapi.Update,
	wh *Warehouse, nCfg CfgNotionAPI) error {

	if upd.Message == nil || (upd.Message.Text == "" && upd.Message.Command() == "") {
		zlog.Debug().Msg("empty message")
		return nil
	}

	switch upd.Message.Command() {
	case "update":
		if err := wh.update(nCfg, true); err != nil {
			SendAnswer(bot, upd, "ошибка обновления данных")
			return err
		}
		SendAnswer(bot, upd, fmt.Sprintf(
			"данные обновлены, получено записей: %d", len(wh.items)))
		return nil
	case "start":
		SendAnswer(bot, upd, "💊 Аптечка")
		return nil
	}

	if err := wh.update(nCfg, false); err != nil {
		SendAnswer(bot, upd, "ошибка обновления данных")
		return err
	}

	foundItems := wh.Query(upd.Message.Text)
	if len(foundItems) < 1 {
		SendAnswer(bot, upd, "ничего не найдено 🧐")
		return nil
	}

	for _, r := range foundItems {
		SendAnswer(bot, upd, r.AsHTML(), "HTML")
	}
	return nil
}

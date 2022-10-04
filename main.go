package main

import (
	"fmt"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func SendAnswer(bot *tgbotapi.BotAPI, upd tgbotapi.Update, msg string, mode ...string) {
	msgCfg := tgbotapi.NewMessage(upd.Message.Chat.ID, msg)
	if 0 < len(mode) {
		msgCfg.ParseMode = mode[0]
	}
	msgCfg.ReplyToMessageID = upd.Message.MessageID
	_, err := bot.Send(msgCfg)
	ehSkip(err)
}

func main() {
	bot, err := tgbotapi.NewBotAPI(CFG.TgAPI.Token)
	eh(err)

	//bot.Debug = true

	store, err := getData()
	eh(err)
	dataExp := time.Now().Add(CFG.NotionAPI.UpdateInterval.Duration)

	LOG.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10
	updates, err := bot.GetUpdatesChan(u)
	eh(err)

	for update := range updates {
		if update.Message == nil || (update.Message.Text == "" && update.Message.Command() == "") {
			continue
		}

		switch update.Message.Command() {
		case "update":
			store, err = getData()
			if err != nil {
				ehSkip(err)
				SendAnswer(bot, update, "ошибка обновления данных")
				continue
			}
			SendAnswer(bot, update, fmt.Sprintf("данные обновлены, получено записей: %d", len(store)))
			continue
		case "start":
			SendAnswer(bot, update, "💊 Аптечка")
			continue
		}

		if time.Now().After(dataExp) {
			store, err = getData()
			if err != nil {
				ehSkip(err)
				SendAnswer(bot, update, "ошибка обновления данных")
				continue
			}
			dataExp = time.Now().Add(CFG.NotionAPI.UpdateInterval.Duration)
		}

		filteredStore := store.Filter(update.Message.Text)
		if 1 > len(filteredStore) {
			SendAnswer(bot, update, "ничего не найдено 🧐")
			continue
		}

		for _, r := range filteredStore {
			SendAnswer(bot, update, r.AsHTML(), "HTML")
		}
	}
}

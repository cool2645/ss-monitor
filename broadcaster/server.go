package broadcaster

import (
	"time"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yanzay/log"
	"strconv"
)

var subscribedChats []int64

func ServeTelegram(apiKey string)  {
	log.Infof("Started serve telegram %s", time.Now())
	bot, err := tg.NewBotAPI(apiKey)
	if err != nil {
		log.Fatal(err)
	}
	u := tg.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		m := update.Message
		if m.IsCommand() && m.Command() == "start" {
			replyMessage(start(m), bot, m)
		}
	}

}

func start(m *tg.Message) string {
	subscribedChats = append(subscribedChats, m.Chat.ID)
	return strconv.FormatInt(m.Chat.ID, 10)
}

func replyMessage(text string, bot *tg.BotAPI, req *tg.Message) {
	msg := tg.NewMessage(req.Chat.ID, text)
	msg.ReplyToMessageID = req.MessageID
	bot.Send(msg)
}
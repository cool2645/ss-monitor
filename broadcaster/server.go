package broadcaster

import (
	"time"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yanzay/log"
	"github.com/cool2645/ss-monitor/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sync"
	"fmt"
)

var subscribedChats []int64
var mux *sync.RWMutex
var Ch chan string

func ServeTelegram(db *gorm.DB, apiKey string, ch chan string) {
	log.Infof("Reading subscribed chats from database %s", time.Now())
	subscribedChats, err := model.ListSubscribers(db)
	if err != nil {
		log.Fatal(err)
	}
	log.Warningf("%v %s", subscribedChats, time.Now())
	log.Infof("Started serve telegram %s", time.Now())
	bot, err := tg.NewBotAPI(apiKey)
	if err != nil {
		log.Fatal(err)
	}
	go pushMessage(bot, ch)
	u := tg.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		m := update.Message
		if m.IsCommand() {
			switch m.Command() {
			case "start":
				replyMessage(start(db, m), bot, m)
			case "stop":
				replyMessage(stop(db, m), bot, m)
			}
		}
	}
}

func pushMessage(bot *tg.BotAPI, ch chan string) {
	var msg string
	for {
		msg = <-ch
		mux.RLock()
		for _, v := range subscribedChats {
			msg := tg.NewMessage(v, msg)
			bot.Send(msg)
		}
		mux.RUnlock()
	}
}

func start(db *gorm.DB, m *tg.Message) string {
	mux.Lock()
	defer mux.Unlock()
	subscribedChats = append(subscribedChats, m.Chat.ID)
	_, err := model.SaveSubscriber(db, m.Chat.ID)
	if err != nil {
		log.Fatal(err)
	}
	return "You have set up subscription for this chat, Yay!"
}

func stop(db *gorm.DB, m *tg.Message) string {
	mux.Lock()
	defer mux.Unlock()
	var newSubscribedChats []int64
	for _, v := range subscribedChats {
		if v != m.Chat.ID {
			newSubscribedChats = append(newSubscribedChats, v)
		}
	}
	subscribedChats = newSubscribedChats
	err := model.RemoveSubscriber(db, m.Chat.ID)
	if err != nil {
		log.Fatal(err)
	}
	return "You'll no longer receive messages from this bot."
}

func replyMessage(text string, bot *tg.BotAPI, req *tg.Message) {
	msg := tg.NewMessage(req.Chat.ID, text)
	msg.ReplyToMessageID = req.MessageID
	bot.Send(msg)
}

func Broadcast(msg string, worker string, class string) {
	msgf := fmt.Sprintf("%s\n%s(%s)", msg, worker, class)
	Ch <- msgf
	return
}

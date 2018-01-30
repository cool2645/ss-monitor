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
	. "github.com/cool2645/ss-monitor/config"
)

var subscribedChats = make(map[int64]int64)
var mux sync.RWMutex
var ch = make(chan string)
var ManagerChan = make(chan int64)
var ManagerNodeChan = make(chan int64)
var bot *tg.BotAPI

func ServeTelegram(db *gorm.DB, apiKey string) {
	log.Infof("Reading subscribed chats from database %s", time.Now())
	subscribers, err := model.ListSubscribers(db)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range subscribers {
		subscribedChats[v] = v
	}
	log.Warningf("%v %s", subscribedChats, time.Now())
	log.Infof("Started serve telegram %s", time.Now())
	bot, err = tg.NewBotAPI(apiKey)
	if err != nil {
		log.Fatal(err)
	}
	go pushMessage(ch)
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
				ReplyMessage(start(db, m), "", "broadcaster", m.Chat.ID)
			case "stop":
				ReplyMessage(stop(db, m), "", "broadcaster", m.Chat.ID)
			case "ping":
				ManagerChan <- m.Chat.ID
			case "status":
				ManagerNodeChan <- m.Chat.ID
			}
		}
	}
}

func pushMessage(c chan string) {
	var m string
	for {
		m = <-c
		mux.RLock()
		for _, v := range subscribedChats {
			msg := tg.NewMessage(v, m)
			msg.ParseMode = "Markdown"
			bot.Send(msg)
		}
		mux.RUnlock()
	}
}

func start(db *gorm.DB, m *tg.Message) string {
	mux.Lock()
	defer mux.Unlock()
	subscribedChats[m.Chat.ID] = m.Chat.ID
	_, err := model.SaveSubscriber(db, m.Chat.ID)
	if err != nil {
		log.Fatal(err)
	}
	return "You have set up subscription for this chat, Yay!"
}

func stop(db *gorm.DB, m *tg.Message) string {
	mux.Lock()
	defer mux.Unlock()
	delete(subscribedChats, m.Chat.ID)
	err := model.RemoveSubscriber(db, m.Chat.ID)
	if err != nil {
		log.Fatal(err)
	}
	return "You'll no longer receive messages from this bot."
}

func ReplyMessage(text string, worker string, class string, reqChatID int64) {
	textf := formatMessage(text, worker, class)
	msg := tg.NewMessage(reqChatID, textf)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func formatMessage(msg string, worker string, class string) (msgf string) {
	if class == "broadcaster" {
		if v, ok := GlobCfg.FRIENDLY_NAME[class]; ok {
			msgf = fmt.Sprintf("%s\n_By %s_", msg, v)
		} else {
			msgf = fmt.Sprintf("%s\n_By %s_", msg, "broadcaster")
		}
		return
	}
	if len(msg) > 0 && msg[len(msg)-1] != '\n' {
		if v, ok := GlobCfg.FRIENDLY_NAME[class]; ok {
			msgf = fmt.Sprintf("%s\n_By %s(%s)_", msg, worker, v)
		} else {
			msgf = fmt.Sprintf("%s\n_By %s(%s)_", msg, worker, class)
		}
	} else {
		if v, ok := GlobCfg.FRIENDLY_NAME[class]; ok {
			msgf = fmt.Sprintf("%s_By %s(%s)_", msg, worker, v)
		} else {
			msgf = fmt.Sprintf("%s_By %s(%s)_", msg, worker, class)
		}
	}
	return
}

func Broadcast(msg string, worker string, class string) {
	msgf := formatMessage(msg, worker, class)
	ch <- msgf
	return
}

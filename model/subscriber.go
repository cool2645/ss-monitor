package model

import (
	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
)

type Subscriber struct {
	ID        uint  `gorm:"AUTO_INCREMENT"`
	ChatID    int64 `gorm:"unique"`
	CreatedAt time.Time
}

func SaveSubscriber(db *gorm.DB, chatID int64) (newSubscriber Subscriber, err error) {
	var subscriber Subscriber
	subscriber.ChatID = chatID
	newSubscriber, err = CreateSubscriber(db, subscriber)
	return
}

func CreateSubscriber(db *gorm.DB, subscriber Subscriber) (newSubscriber Subscriber, err error) {
	err = db.Create(&subscriber).Error
	if err != nil {
		err = errors.Wrap(err, "SaveSubscriber")
		return
	}
	newSubscriber = subscriber
	return
}

func RemoveSubscriber(db *gorm.DB, chatID int64) (err error) {
	err = db.Delete(Subscriber{}, "chat_id = ?", chatID).Error
	if err != nil {
		err = errors.Wrap(err, "RemoveSubscriber")
		return
	}
	return
}
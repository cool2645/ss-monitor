package model

import "time"

type Subscriber struct {
	ID        uint  `gorm:"AUTO_INCREMENT"`
	ChatID    int64 `gorm:"unique"`
	CreatedAt time.Time
}

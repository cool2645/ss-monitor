package model

import "time"

type Task struct {
	ID         uint   `gorm:"AUTO_INCREMENT"`
	CallbackID uint   `gorm:"index"`
	Class      string `gorm:"not null;index"`
	NodeID     uint   `gorm:"index"`
	Node       Node
	IPVer      uint   `gorm:"index"`
	State      string
	Worker     string
	Log        string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	ServerName string
	SsJson     string
}

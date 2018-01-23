package model

import "time"

type Heartbeat struct {
	ID        uint      `gorm:"AUTO_INCREMENT"`
	Class     string    `gorm:"not null;index"`
	IPVer     uint      `gorm:"index"`
	Name      string
	CreatedAt time.Time `gorm:"index"`
}

func ServicesStatus() {

}

func SyncServiceStatus() {

}
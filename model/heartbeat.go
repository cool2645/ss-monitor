package model

import (
	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
)

type Heartbeat struct {
	ID        uint      `gorm:"AUTO_INCREMENT"`
	Class     string    `gorm:"not null;index"`
	IPVer     uint      `gorm:"index"`
	Name      string
	CreatedAt time.Time `gorm:"index"`
}

func WorkersStatus() {

}

func SyncWorkerStatus() {

}

func SaveHeartbeat(db *gorm.DB, class string, ipVer uint, name string) (newHeartbeat Heartbeat, err error) {
	var heartbeat Heartbeat
	heartbeat.Class = class
	heartbeat.IPVer = ipVer
	heartbeat.Name = name
	newHeartbeat, err = CreateHeartbeat(db, heartbeat)
	return
}

func CreateHeartbeat(db *gorm.DB, heartbeat Heartbeat) (newHeartbeat Heartbeat, err error) {
	err = db.Create(&heartbeat).Error
	if err != nil {
		err = errors.Wrap(err, "CreateHeartbeat")
		return
	}
	newHeartbeat = heartbeat
	return
}
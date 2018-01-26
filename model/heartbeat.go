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
	Time      int64
	CreatedAt time.Time `gorm:"index"`
}

func WorkersStatus() {

}

func SyncWorkerStatus() {

}

func SaveHeartbeat(db *gorm.DB, class string, ipVer uint, name string) (newHeartbeat Heartbeat, err error) {
	var count int
	var heartbeat Heartbeat
	err = db.Model(&Heartbeat{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		err = errors.Wrap(err, "SaveHeartbeat: CountHeartbeat")
		return
	}
	if count == 0 {
		heartbeat.Class = class
		heartbeat.IPVer = ipVer
		heartbeat.Time = time.Now().Unix()
		heartbeat.Name = name
		err = db.Create(&heartbeat).Error
		if err != nil {
			err = errors.Wrap(err, "SaveHeartbeat: CreateHeartbeat")
			return
		}
		newHeartbeat = heartbeat
		return
	} else {
		err = db.Where("name = ?", name).First(&heartbeat).Error
		if err != nil {
			err = errors.Wrap(err, "SaveHeartbeat: QueryHeartbeat")
			return
		}
		heartbeat.Class = class
		heartbeat.IPVer = ipVer
		heartbeat.Time = time.Now().Unix()
		err = db.Model(&heartbeat).Updates(heartbeat).Error
		if err != nil {
			err = errors.Wrap(err, "SaveHeartbeat: UpdateHeartbeat")
			return
		}
	}
	return
}
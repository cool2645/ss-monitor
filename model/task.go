package model

import (
	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
	"sync"
)

type Task struct {
	ID         uint   `gorm:"AUTO_INCREMENT"`
	CallbackID uint   `gorm:"index"`
	Class      string `gorm:"not null;index"`
	NodeID     uint   `gorm:"index"`
	Node       Node
	IPVer      uint   `gorm:"index"`
	State      string
	Worker     string
	Log        string `sql:"type:text;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	ServerName string
	SsJson     string
}

var mux sync.Mutex

func CreateTask(db *gorm.DB, task Task) (newTask Task, err error) {
	// Default Value
	task.State = "Queuing"

	err = db.Create(&task).Error
	if err != nil {
		err = errors.Wrap(err, "CreateTask")
		return
	}
	newTask = task
	return
}

func GetTasks(db *gorm.DB, class string, state string, ipVer string, callbackID string, order string, page uint) (tasks []Task, err error) {
	noLog := "id, callback_id, class, node_id, ip_ver, state, worker, created_at, updated_at, server_name, ss_json"
	switch class {
	case "tester":
		err = db.Where("callback_id like ?", callbackID).Where("class like ?", class).Where("ip_ver like ?", ipVer).Where("state like ?", state).
			Order("id " + order).Offset((page - 1) * 10).Limit(10).
			Select(noLog).Preload("Node").Find(&tasks).Error
	default:
		err = db.Where("callback_id like ?", callbackID).Where("class like ?", class).Where("state like ?", state).
			Order("id " + order).Offset((page - 1) * 10).Limit(10).
			Select(noLog).Preload("Node").Find(&tasks).Error
	}
	if err != nil {
		err = errors.Wrap(err, "GetTasks")
		return
	}
	return
}

func GetTask(db *gorm.DB, id uint) (task Task, err error) {
	err = db.Where("id = ?", id).Preload("Node").Find(&task).Error
	if err != nil {
		err = errors.Wrap(err, "GetTask")
		return
	}
	return
}

func AssignTask(db *gorm.DB, id uint, worker string) (err error) {

	mux.Lock()
	defer mux.Unlock()

	var task Task
	err = db.Where("id = ?", id).Find(&task).Error
	if err != nil {
		err = errors.Wrap(err, "AssignTask: Find task")
		return
	}
	if task.State == "Queuing" {
		task.Worker = worker
		task.State = "Starting"
		err = db.Model(&task).Updates(task).Error
		if err != nil {
			err = errors.Wrap(err, "AssignTask: Update task")
			return
		}
		return
	}
	err = errors.New("Not queuing task")
	err = errors.Wrap(err, "AssignTask: Check task status")
	return
}

func UpdateTaskStatus(db *gorm.DB, id uint, worker string, state string, log string) (err error) {
	var task Task
	err = db.Where("id = ?", id).Find(&task).Error
	if err != nil {
		err = errors.Wrap(err, "SyncTaskStatus: Find task")
		return
	}
	if task.Worker == worker {
		task.State = state
		task.Log = log
		err = db.Model(&task).Updates(task).Error
		if err != nil {
			err = errors.Wrap(err, "SyncTaskStatus: Update task")
			return
		}
		return
	}
	err = errors.New("Not assigned worker")
	err = errors.Wrap(err, "SyncTaskStatus: Check task status")
	return
}

func ResetTask(db *gorm.DB, id uint) (err error) {
	var task Task
	err = db.Where("id = ?", id).Find(&task).Error
	if err != nil {
		err = errors.Wrap(err, "ResetTask: Find task")
		return
	}
	if task.State != "Queuing" {
		task.State = "Queuing"
		task.Log = ""
		task.Worker = ""
		err = db.Save(task).Error
		if err != nil {
			err = errors.Wrap(err, "ResetTask: Update task")
			return
		}
	}
	return
}

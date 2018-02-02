package model

import (
	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
	"sync"
)

type Task struct {
	ID         uint      `gorm:"AUTO_INCREMENT"`
	CallbackID uint      `gorm:"index"`
	Class      string    `gorm:"not null;index"`
	NodeID     uint      `gorm:"index"`
	Node       Node
	IPVer      uint      `gorm:"index"`
	State      string
	Worker     string
	Result     string    `sql:"type:text;"`
	Log        string    `sql:"type:text;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time `gorm:"index"`
	ServerName string
	SsJson     string
}

var mux sync.Mutex

func CreateTask(db *gorm.DB, task Task) (newTask Task, err error) {

	task.State = "Queuing"

	err = db.Create(&task).Error
	if err != nil {
		err = errors.Wrap(err, "CreateTask")
		return
	}
	newTask = task
	return
}

func GetTasks(db *gorm.DB, class string, state string, ipVer string, order string, page uint, perPage uint) (tasks []Task, total uint, err error) {
	noLog := "id, callback_id, class, node_id, ip_ver, state, result, worker, created_at, updated_at, server_name, ss_json"
	if perPage == 0 {
		perPage = 10
	}
	err = db.Where("class like ?", class).Where("ip_ver like ?", ipVer).Where("state like ?", state).
		Order("updated_at " + order).Offset((page - 1) * perPage).Limit(perPage).
		Select(noLog).Preload("Node").Find(&tasks).Count(&total).Error
	if err != nil {
		err = errors.Wrap(err, "GetTasks")
		return
	}
	return
}

func GetTasksByCallbackID(db *gorm.DB, callbackID uint) (tasks []Task, err error) {
	noLog := "id, callback_id, class, node_id, ip_ver, state, result, worker, created_at, updated_at, server_name, ss_json"
	err = db.Where("callback_id = ?", callbackID).
		Order("id asc").Select(noLog).Preload("Node").Find(&tasks).Error
	if err != nil {
		err = errors.Wrap(err, "GetTasksByCallbackID")
		return
	}
	return
}

func GetTaskByNode(db *gorm.DB, nodeID uint, class string, ipVer string) (task Task, err error) {
	noLog := "id, callback_id, class, node_id, ip_ver, state, result, worker, created_at, updated_at, server_name, ss_json"
	err = db.Where("node_id = ?", nodeID).Where("class like ?", class).
		Where("state in (?)", []string{"Passing", "Failing", "Shiny☆"}).
		Order("updated_at desc").Select(noLog).First(&task).Error
	if err != nil {
		err = errors.Wrap(err, "GetTaskByNode")
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

func UpdateTaskStatus(db *gorm.DB, id uint, worker string, state string, result string, log string) (err error) {
	var task Task
	err = db.Where("id = ?", id).Find(&task).Error
	if err != nil {
		err = errors.Wrap(err, "SyncTaskStatus: Find task")
		return
	}
	if task.Worker == worker {
		task.State = state
		task.Result = result
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
		task.Result = ""
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

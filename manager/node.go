package manager

import (
	"github.com/cool2645/ss-monitor/model"
	. "github.com/cool2645/ss-monitor/config"
	"github.com/yanzay/log"
	"github.com/cool2645/ss-monitor/broadcaster"
	"github.com/pkg/errors"
	"fmt"
	"sync"
)

type node struct {
	Name       string
	IsCleaning bool
	Status     map[string]model.Task
}

var nodes map[uint]node
var nodeMux sync.RWMutex

func InitNodes() {
	nodeMux.Lock()
	defer nodeMux.Unlock()
	nodes = make(map[uint]node)
	ns, err := model.GetNodes(model.Db)
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, v := range ns {
		nodes[v.ID] = node{Name: v.Name, IsCleaning: v.IsCleaning, Status: make(map[string]model.Task)}
		if v.EnableWatching {
			t, err := model.GetTaskByNode(model.Db, v.ID, "watcher", "%")
			if err != nil {
				log.Error(err)
				nodes[v.ID].Status["CN"] = model.Task{}
			} else {
				nodes[v.ID].Status["CN"] = t
			}
		}
		if v.EnableIPv4Testing {
			t, err := model.GetTaskByNode(model.Db, v.ID, "tester", "4")
			if err != nil {
				log.Error(err)
				nodes[v.ID].Status["SS"] = model.Task{}
			} else {
				nodes[v.ID].Status["SS"] = t
			}
		}
		if v.EnableIPv6Testing {
			t, err := model.GetTaskByNode(model.Db, v.ID, "tester", "6")
			if err != nil {
				log.Error(err)
				nodes[v.ID].Status["SS-IPv6"] = model.Task{}
			} else {
				nodes[v.ID].Status["SS-IPv6"] = t
			}
		}
	}
}

func GetNodeStatus() map[uint]node {
	nodeMux.RLock()
	defer nodeMux.RUnlock()
	return nodes
}

func reportNodeStatus(ch chan int64) {
	for {
		reqChatID := <-ch
		var msg string
		nodeMux.RLock()
		for _, node := range nodes {
			for k, s := range node.Status {
				if node.IsCleaning {
					msg += fmt.Sprintf("%s 🔶 Cleaning\n", node.Name)
					continue
				}
				switch s.State {
				case "Passing":
					msg += fmt.Sprintf("%s-%s 🔵 Last tested at %s\n", node.Name, k, s.UpdatedAt.Format("15:04, Jan 2"))
				case "Failing":
					msg += fmt.Sprintf("%s-%s 🔴 Last tested at %s\n", node.Name, k, s.UpdatedAt.Format("15:04, Jan 2"))
				}
			}
		}
		nodeMux.RUnlock()
		broadcaster.ReplyMessage(msg, GlobCfg.MANAGER_NAME, "manager", reqChatID)
	}
}

func TaskCallback(taskID uint, worker string) (err error) {
	nodeMux.Lock()
	defer nodeMux.Unlock()
	task, err := model.GetTask(model.Db, taskID)
	if err != nil {
		err = errors.Wrap(err, "TaskCallback")
		return
	}
	if task.Worker != worker {
		err = errors.New("Not assigned worker")
		err = errors.Wrap(err, "TaskCallback")
	}
	// Update node status
	switch task.Class {
	case "watcher":
		nodes[task.NodeID].Status["CN"] = task
	case "tester":
		if task.IPVer == 6 {
			nodes[task.NodeID].Status["SS-IPv6"] = task
		} else {
			nodes[task.NodeID].Status["SS"] = task
		}
	}
	// Generate Report
	brotherTasks, err := model.GetTasksByCallbackID(model.Db, task.CallbackID)
	if err != nil {
		log.Error(err)
		msg := fmt.Sprintf("❗️ Error occured when task #%d callback with callback id #%d: %s", task.ID, task.CallbackID, err.Error())
		broadcaster.Broadcast(msg, GlobCfg.MANAGER_NAME, "manager")
		return
	}
	if len(brotherTasks) != 0 {
		var passedTasks, failedTasks, pendingTasks []model.Task
		for _, v := range brotherTasks {
			switch v.State {
			case "Passing":
				passedTasks = append(passedTasks, v)
			case "Shiny☆":
				passedTasks = append(passedTasks, v)
			case "Failing":
				failedTasks = append(failedTasks, v)
			default:
				pendingTasks = append(pendingTasks, v)
			}
		}
		if len(pendingTasks) == 0 {
			var msg string
			fatherTask, err := model.GetTask(model.Db, task.CallbackID)
			if err == nil {
				if len(failedTasks) != 0 {
					msg += fmt.Sprintf("*🔴 task #%d (Created at %s)*\n", fatherTask.ID, fatherTask.CreatedAt.Format("15:04"))
				} else {
					msg += fmt.Sprintf("*🔵 task #%d (Created at %s)*\n", fatherTask.ID, fatherTask.CreatedAt.Format("15:04"))
				}
			}
			for _, v := range failedTasks {
				switch v.Class {
				case "watcher":
					msg += fmt.Sprintf("    🔴 task #%d Watch %s: %s", v.ID, v.Node.Name, v.State)
				case "tester":
					if v.IPVer == 6 {
						msg += fmt.Sprintf("    🔴 task #%d Test %s with IPv6: %s", v.ID, v.Node.Name, v.State)
					} else {
						msg += fmt.Sprintf("    🔴 task #%d Test %s: %s", v.ID, v.Node.Name, v.State)
					}
				}
			}
			for _, v := range passedTasks {
				switch v.Class {
				case "watcher":
					msg += fmt.Sprintf("    🔵 task #%d Watch %s: %s", v.ID, v.Node.Name, v.State)
				case "tester":
					if v.IPVer == 6 {
						msg += fmt.Sprintf("    🔵 task #%d Test %s with IPv6: %s", v.ID, v.Node.Name, v.State)
					} else {
						msg += fmt.Sprintf("    🔵 task #%d Test %s: %s", v.ID, v.Node.Name, v.State)
					}
				}
			}
			broadcaster.Broadcast(msg, GlobCfg.MANAGER_NAME, "manager")
		}
	}
	return
}

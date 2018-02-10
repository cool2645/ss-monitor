package manager

import (
	"github.com/cool2645/ss-monitor/model"
	. "github.com/cool2645/ss-monitor/config"
	"github.com/yanzay/log"
	"github.com/cool2645/ss-monitor/broadcaster"
	"fmt"
	"github.com/pkg/errors"
	"time"
)

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
	if task.Class == "cleaner" && task.Result == "Shiny☆" {
		model.ResetNode(model.Db, task.NodeID)
	}
	// Update node status
	if _, ok := nodes[task.NodeID]; ok {
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
	}
	// Generate Report
	if task.CallbackID == 0 {
		return
	}
	brotherTasks, err := model.GetAllTasksByCallbackID(model.Db, task.CallbackID)
	if err != nil {
		log.Error(err)
		msg := fmt.Sprintf("❗️ Error occured when task #%d callback with callback id #%d: %s\n", task.ID, task.CallbackID, err.Error())
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
					model.UpdateTaskStatus(model.Db, fatherTask.ID, GlobCfg.MANAGER_NAME, "Failing", "", "")
					msg += fmt.Sprintf("*🔴 task #%d (Created at %s)*\n", fatherTask.ID, fatherTask.CreatedAt.Format("15:04"))
				} else {
					model.UpdateTaskStatus(model.Db, fatherTask.ID, GlobCfg.MANAGER_NAME, "Passing", "", "")
					msg += fmt.Sprintf("*🔵 task #%d (Created at %s)*\n", fatherTask.ID, fatherTask.CreatedAt.Format("15:04"))
				}
			}
			var toCleanNodeIDs []uint
			for _, v := range failedTasks {
				switch v.Class {
				case "watcher":
					toCleanNodeIDs = append(toCleanNodeIDs, v.NodeID)
					msg += fmt.Sprintf("    🔴 task #%d Watch %s: %s\n", v.ID, v.Node.Name, v.State)
				case "tester":
					if v.IPVer == 6 {
						msg += fmt.Sprintf("    🔴 task #%d Test %s with IPv6: %s\n", v.ID, v.Node.Name, v.State)
					} else {
						msg += fmt.Sprintf("    🔴 task #%d Test %s: %s\n", v.ID, v.Node.Name, v.State)
					}
				}
			}
			for _, v := range passedTasks {
				switch v.Class {
				case "watcher":
					msg += fmt.Sprintf("    🔵 task #%d Watch %s: %s\n", v.ID, v.Node.Name, v.State)
				case "tester":
					if v.IPVer == 6 {
						msg += fmt.Sprintf("    🔵 task #%d Test %s with IPv6: %s\n", v.ID, v.Node.Name, v.State)
					} else {
						msg += fmt.Sprintf("    🔵 task #%d Test %s: %s\n", v.ID, v.Node.Name, v.State)
					}
				}
			}
			broadcaster.Broadcast(msg, GlobCfg.MANAGER_NAME, "manager")
			for _, nodeID := range toCleanNodeIDs {
				node, err := model.GetNode(model.Db, nodeID)
				if err != nil {
					log.Error(err)
					continue
				}
				if node.EnableCleaning && !node.IsCleaning {
					model.SetNodeCleaning(model.Db, node.ID)
					task := model.Task{Class: "cleaner", NodeID: node.ID, IPVer: 4}
					task, err = model.CreateTask(model.Db, task)
					if err != nil {
						log.Error(err)
						continue
					}
					msg := fmt.Sprintf("🔶 Attempt to clean node %s by task #%s\n", node.Name, task.ID)
					broadcaster.Broadcast(msg, GlobCfg.MANAGER_NAME, "manager")
				}
			}
		}
	}
	return
}

func createTask() {
	ns, err := model.GetNodes(model.Db)
	if err != nil {
		log.Fatal(err)
		return
	}
	var tasks []model.Task
	for _, v := range ns {
		if !v.IsCleaning {
			if v.EnableWatching && len(workers["watcher"]) > 0 {
				task := model.Task{Class: "watcher", NodeID: v.ID, IPVer: 4}
				tasks = append(tasks, task)
			}
			if v.EnableIPv4Testing && len(workers["ipv4tester"]) > 0 {
				task := model.Task{Class: "tester", NodeID: v.ID, IPVer: 4}
				tasks = append(tasks, task)
			}
			if v.EnableIPv6Testing && len(workers["ipv6tester"]) > 0 {
				task := model.Task{Class: "tester", NodeID: v.ID, IPVer: 6}
				tasks = append(tasks, task)
			}
		}
	}
	log.Debugf("Creating %v tasks", len(tasks))
	if len(tasks) > 0 {
		task := model.Task{Class: "manager"}
		task, err = model.CreateTask(model.Db, task)
		if err != nil {
			log.Fatal(err)
			return
		}
		err = model.AssignTask(model.Db, task.ID, GlobCfg.MANAGER_NAME)
		if err != nil {
			log.Fatal(err)
			return
		}
		for _, v := range tasks {
			v.CallbackID = task.ID
			_, err = model.CreateTask(model.Db, v)
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func taskSchdLoop()  {
	for GlobCfg.MANAGER_SCHEDULE {
		createTask()
		time.Sleep(time.Minute * time.Duration(GlobCfg.MANAGER_INTERVAL))
	}
}
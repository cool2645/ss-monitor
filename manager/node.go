package manager

import (
	"github.com/cool2645/ss-monitor/model"
	. "github.com/cool2645/ss-monitor/config"
	"github.com/yanzay/log"
	"github.com/cool2645/ss-monitor/broadcaster"
	"fmt"
)

type node struct {
	Name       string
	IsCleaning bool
	Status     map[string]model.Task
}

var nodes map[uint]node

func InitNodes() {
	nodes = make(map[uint]node)
	ns, err := model.GetNodes(model.Db)
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, v := range ns {
		nodes[v.ID] = node{Name: v.Name, IsCleaning: v.IsCleaning, Status: make(map[string]model.Task)}
		if v.EnableWatching {
			t, err := model.GetLastFinishedTask(model.Db, v.ID, "watcher", "%")
			if err != nil {
				log.Error(err)
				nodes[v.ID].Status["CN"] = model.Task{}
			} else {
				nodes[v.ID].Status["CN"] = t
			}
		}
		if v.EnableIPv4Testing {
			t, err := model.GetLastFinishedTask(model.Db, v.ID, "tester", "4")
			if err != nil {
				log.Error(err)
				nodes[v.ID].Status["SS"] = model.Task{}
			} else {
				nodes[v.ID].Status["SS"] = t
			}
		}
		if v.EnableIPv6Testing {
			t, err := model.GetLastFinishedTask(model.Db, v.ID, "tester", "6")
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
	return nodes
}

func reportNodeStatus(ch chan int64) {
	for {
		reqChatID := <-ch
		var msg string
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
		broadcaster.ReplyMessage(msg, GlobCfg.MANAGER_NAME, "manager", reqChatID)
	}
}

package manager

import (
	"github.com/cool2645/ss-monitor/model"
	. "github.com/cool2645/ss-monitor/config"
	"time"
	"github.com/pkg/errors"
	"sync"
	"fmt"
	"github.com/cool2645/ss-monitor/broadcaster"
)

var mux sync.Mutex
var workers = make(map[string]map[string]model.Heartbeat)

func Init() {
	workers["watcher"] = make(map[string]model.Heartbeat)
	workers["ipv4tester"] = make(map[string]model.Heartbeat)
	workers["ipv6tester"] = make(map[string]model.Heartbeat)
	workers["cleaner"] = make(map[string]model.Heartbeat)
	heartbeats, err := model.GetHeartbeats(model.Db, time.Now().Unix()-GlobCfg.HEARTBEAT_ALARM, "watcher", 10)
	if err != nil {
		panic(err)
		return
	}
	for _, v := range heartbeats {
		workers["watcher"][v.Name] = v
	}
	heartbeats, err = model.GetHeartbeats(model.Db, time.Now().Unix()-GlobCfg.HEARTBEAT_ALARM, "tester", 4)
	if err != nil {
		panic(err)
		return
	}
	for _, v := range heartbeats {
		workers["ipv4tester"][v.Name] = v
	}
	heartbeats, err = model.GetHeartbeats(model.Db, time.Now().Unix()-GlobCfg.HEARTBEAT_ALARM, "tester", 6)
	if err != nil {
		panic(err)
		return
	}
	for _, v := range heartbeats {
		workers["ipv6tester"][v.Name] = v
	}
	heartbeats, err = model.GetHeartbeats(model.Db, time.Now().Unix()-GlobCfg.HEARTBEAT_ALARM, "cleaner", 10)
	if err != nil {
		panic(err)
		return
	}
	for _, v := range heartbeats {
		workers["cleaner"][v.Name] = v
	}
	go monitorWorkers()
	go reportWorkerStatus(broadcaster.ManagerChan)
	InitNodes()
	go reportNodeStatus(broadcaster.ManagerNodeChan)
}

func monitorWorkers() {
	for {
		checkWorkerStatus()
		time.Sleep(time.Second * time.Duration(GlobCfg.MONITOR_INTERVAL))
	}
}

func GetWorkerStatus() (map[string]map[string]model.Heartbeat) {
	return workers
}

func reportWorkerStatus(ch chan int64) {
	for {
		reqChatID := <-ch
		var msg string
		workers := workers
		for wck, ws := range workers {
			if len(ws) == 0 {
				msg += fmt.Sprintf("*🔴 %d %s%s currently online*\n", len(ws), wck, singleOrPlural(len(ws)))
			} else {
				msg += fmt.Sprintf("*🔵 %d %s%s currently online*\n", len(ws), wck, singleOrPlural(len(ws)))
			}
			for _, v := range ws {
				msg += fmt.Sprintf("    🔵 %s: Last heartbeat at %s\n", v.Name, time.Unix(v.Time, 0).Format("15:04"))
			}
		}
		broadcaster.ReplyMessage(msg, GlobCfg.MANAGER_NAME, "manager", reqChatID)
	}
}

func UpdateWorkerStatus(class string, ipVer uint, name string) (heartbeat model.Heartbeat, err error) {
	var workerClassKey []string
	switch class {
	case "tester":
		if ipVer == 4 {
			workerClassKey = append(workerClassKey, "ipv4tester")
		} else if ipVer == 6 {
			workerClassKey = append(workerClassKey, "ipv6tester")
		} else {
			workerClassKey = append(workerClassKey, "ipv4tester")
			workerClassKey = append(workerClassKey, "ipv6tester")
		}
	default:
		workerClassKey = append(workerClassKey, class)
	}
	t := time.Now().Unix()
	var worker model.Heartbeat
	var newWorkerWCK []string
	for _, wck := range workerClassKey {
		if val, ok := workers[wck][name]; ok {
			worker = val
		} else {
			worker.Name = name
			newWorkerWCK = append(newWorkerWCK, wck)
		}
		worker.Class = class
		worker.IPVer = ipVer
		worker.Time = t
		workers[wck][name] = worker
	}
	heartbeat, err = model.SaveHeartbeat(model.Db, class, ipVer, name, t)
	if err != nil {
		err = errors.Wrap(err, "UpdateWorkerStatus")
		return
	}
	for _, wck := range newWorkerWCK {
		alert(make([]model.Heartbeat, 0), workers[wck], wck)
	}
	return
}

func checkWorkerStatus() {
	mux.Lock()
	defer mux.Unlock()
	for i, ws := range workers {
		var offlineList []model.Heartbeat
		for k, v := range ws {
			if v.Time < time.Now().Unix()-GlobCfg.HEARTBEAT_ALARM {
				offlineList = append(offlineList, v)
				delete(workers[i], k)
			}
		}
		if len(offlineList) != 0 {
			alert(offlineList, workers[i], i)
		}
	}
}

func singleOrPlural(cnt int) (string) {
	if cnt > 1 {
		return "s"
	} else {
		return ""
	}
}

func rootOrSingular(cnt int) (string) {
	if cnt > 1 {
		return ""
	} else {
		return "s"
	}
}

func alert(offlineWorkers []model.Heartbeat, onlineWorkers map[string]model.Heartbeat, wck string) {
	var msg string
	if len(offlineWorkers) > 0 {
		msg += fmt.Sprintf("*🔴 %d %s%s get%s offline*\n", len(offlineWorkers), wck,
			singleOrPlural(len(offlineWorkers)), rootOrSingular(len(offlineWorkers)))
	}
	for _, v := range offlineWorkers {
		msg += fmt.Sprintf("    🔴 %s: Last heartbeat at %s\n", v.Name, time.Unix(v.Time, 0).Format("15:04"))
	}
	if len(onlineWorkers) == 0 {
		msg += fmt.Sprintf("*🔴 %d %s%s currently online*\n", len(onlineWorkers), wck, singleOrPlural(len(onlineWorkers)))
	} else {
		msg += fmt.Sprintf("*🔵 %d %s%s currently online*\n", len(onlineWorkers), wck, singleOrPlural(len(onlineWorkers)))
	}
	for _, v := range onlineWorkers {
		msg += fmt.Sprintf("    🔵 %s: Last heartbeat at %s\n", v.Name, time.Unix(v.Time, 0).Format("15:04"))
	}
	broadcaster.Broadcast(msg, GlobCfg.MANAGER_NAME, "manager")
}

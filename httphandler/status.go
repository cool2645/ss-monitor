package httphandler

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/yanzay/log"
	"strconv"
	"github.com/cool2645/ss-monitor/manager"
)

func HandleHeartbeat(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if !authAccessKey(w, req) {
		return
	}
	req.ParseForm()
	if len(req.Form["class"]) != 1 {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Invalid worker class.",
		}
		responseJson(w, res, http.StatusBadRequest)
		return
	}
	class := req.Form["class"][0]
	var ipVer uint
	if len(req.Form["ip_ver"]) == 1 {
		ipVer64, err := strconv.ParseUint(req.Form["ip_ver"][0], 10, 32)
		if err != nil {
			log.Error(err)
		}
		ipVer = uint(ipVer64)
	}
	name := ps.ByName("name")
	_, err := manager.UpdateWorkerStatus(class, ipVer, name)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred handling heartbeat: " + err.Error(),
		}
		responseJson(w, res, http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"msg":    "success",
	}
	responseJson(w, res, http.StatusOK)
}

func GetWorkerStatus(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	workers := manager.ReportWorkerStatus()
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   workers,
	}
	responseJson(w, res, http.StatusOK)
}

package httphandler

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/cool2645/ss-monitor/model"
	"github.com/yanzay/log"
	"strconv"
)

func HandleHeartbeat(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	req.ParseForm()
	if (len(req.Form["class"]) != 1) {
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
	if (len(req.Form["ip_ver"]) == 1) {
		ipVer64, err := strconv.ParseUint(req.Form["ip_ver"][0], 10, 32)
		if err != nil {
			log.Error(err)
		}
		ipVer = uint(ipVer64)
	}
	name := ps.ByName("name")
	_, err := model.SaveHeartbeat(model.Db, class, ipVer, name)
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

package httphandler

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/yanzay/log"
	"encoding/json"
	"github.com/cool2645/ss-monitor/broadcaster"
)

func responseJson(w http.ResponseWriter, data map[string]interface{}) {
	resJson, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error occurred encoding response.", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resJson)
	return
}

func Broadcast(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	req.ParseForm()
	if (len(req.Form["msg"]) != 0 && len(req.Form["worker"]) != 0 && len(req.Form["class"]) != 0) {
		broadcaster.Broadcast(req.Form["msg"][0], req.Form["worker"][0], req.Form["class"][0])
		res := map[string]interface{}{
			"code":   200,
			"result": true,
			"msg":    "success",
		}
		responseJson(w, res)
		return
	} else {
		res := map[string]interface{}{
			"code":   200,
			"result": false,
			"msg":    "invalid args",
		}
		responseJson(w, res)
	}
}

func Pong(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	res := map[string]interface{}{
		"code":   200,
		"result": true,
		"msg":    "OK",
	}
	responseJson(w, res)
	return
}

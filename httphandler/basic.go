package httphandler

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/yanzay/log"
	"encoding/json"
	"github.com/cool2645/ss-monitor/broadcaster"
)

func responseJson(w http.ResponseWriter, data map[string]interface{}, httpStatusCode int) {
	resJson, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error occurred encoding response.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	w.Write(resJson)
	return
}

func Broadcast(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	req.ParseForm()
	if (len(req.Form["msg"]) != 0 && len(req.Form["worker"]) != 0 && len(req.Form["class"]) != 0) {
		broadcaster.Broadcast(req.Form["msg"][0], req.Form["worker"][0], req.Form["class"][0])
		res := map[string]interface{}{
			"code":   http.StatusOK,
			"result": true,
			"msg":    "success",
		}
		responseJson(w, res, http.StatusOK)
		return
	} else {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "invalid args",
		}
		responseJson(w, res, http.StatusBadRequest)
	}
}

func Pong(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"msg":    "OK",
	}
	responseJson(w, res, http.StatusOK)
	return
}

package httphandler

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/yanzay/log"
	"encoding/json"
)

func responseJson(w http.ResponseWriter, data map[string]interface{}) {
	res_json, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error occurred encoding response.", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res_json)
	return
}

func Pong(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	res := map[string]interface{}{
		"code": 200,
		"result": true,
		"msg":    "OK",
	}
	responseJson(w, res)
	return
}
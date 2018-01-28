package httphandler

import (
	"net/http"
	. "github.com/cool2645/ss-monitor/config"
)

func authAccessKey(w http.ResponseWriter, req *http.Request) (result bool) {
	if !checkAccessKey(req) {
		res := map[string]interface{}{
			"code": http.StatusUnauthorized,
			"result": false,
			"msg":    "Authorize failed.",
		}
		responseJson(w, res, http.StatusUnauthorized)
		result = false
		return
	}
	result = true
	return
}

func checkAccessKey(req *http.Request) (result bool) {
	result = req.Header.Get("X-Access-Key") == GlobCfg.API_KEY
	return
}

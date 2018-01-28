package httphandler

import (
	"net/http"
	. "github.com/cool2645/ss-monitor/config"
	"github.com/julienschmidt/httprouter"
	"github.com/astaxie/beego/session"
)

var globalSessions *session.Manager

func InitSession()  {
	globalSessions, _ = session.NewManager("memory", &session.ManagerConfig{CookieName: "SSMonitorSession", EnableSetCookie: true, Gclifetime: 3600})
	go globalSessions.GC()
}

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

func authAdmin(w http.ResponseWriter, req *http.Request) (result bool) {
	if !checkAdmin(w, req) {
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

func checkAdmin(w http.ResponseWriter, req *http.Request) (result bool) {
	sess, _ := globalSessions.SessionStart(w, req)
	defer sess.SessionRelease(w)
	priv := sess.Get("privilege");
	result = !(priv == nil || priv.(string) != "admin")
	return
}

func auth(w http.ResponseWriter, req *http.Request) (result bool) {
	if !checkAccessKey(req) && !checkAdmin(w, req) {
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

func Login(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	sess, _ := globalSessions.SessionStart(w, req)
	defer sess.SessionRelease(w)
	if username := sess.Get("username"); username != nil {
		res := map[string]interface{}{
			"code": http.StatusFound,
			"result": false,
			"msg":    "Already logged in as: " + username.(string),
		}
		responseJson(w, res, http.StatusFound)
	} else {
		req.ParseForm()
		if len(req.Form["username"]) != 1 {
			res := map[string]interface{}{
				"code": http.StatusBadRequest,
				"result": false,
				"msg":    "Invalid username.",
			}
			responseJson(w, res, http.StatusBadRequest)
			return
		}
		username = req.Form["username"][0]
		if len(req.Form["password"]) != 1 {
			res := map[string]interface{}{
				"code": http.StatusBadRequest,
				"result": false,
				"msg":    "Invalid password.",
			}
			responseJson(w, res, http.StatusBadRequest)
			return
		}
		password := req.Form["password"][0]
		for _, admin := range GlobCfg.ADMIN {
			if admin.Username == username && admin.Password == password {
				sess.Set("username", username)
				sess.Set("privilege", "admin")
				res := map[string]interface{}{
					"code": http.StatusOK,
					"result": true,
					"msg":    "Successfully logged in as: " + username.(string),
				}
				responseJson(w, res, http.StatusOK)
				return
			}
		}
		res := map[string]interface{}{
			"code": http.StatusUnauthorized,
			"result": false,
			"msg":    "No such user or password mismatch.",
		}
		responseJson(w, res, http.StatusUnauthorized)
		return
	}
}

func Logout(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	sess, _ := globalSessions.SessionStart(w, req)
	defer sess.SessionRelease(w)
	sess.Delete("username")
	sess.Delete("privilege")
	res := map[string]interface{}{
		"code": http.StatusOK,
		"result": true,
		"msg":    "Successfully logged out.",
	}
	responseJson(w, res, http.StatusOK)
}
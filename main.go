package main

import (
	"github.com/BurntSushi/toml"
	"net/http"
	"github.com/urfave/negroni"
	"strconv"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/cool2645/ss-monitor/broadcaster"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yanzay/log"
	"github.com/cool2645/ss-monitor/model"
	"github.com/cool2645/ss-monitor/httphandler"
	. "github.com/cool2645/ss-monitor/config"
	"github.com/cool2645/ss-monitor/manager"
)

var mux = httprouter.New()

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "app/index.html")
}

func main() {

	_, err := toml.DecodeFile("config.toml", &GlobCfg)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open("mysql", ParseDSN(GlobCfg))
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Database init done")
	defer db.Close()

	db.AutoMigrate(&model.Subscriber{}, &model.Heartbeat{}, &model.Node{}, &model.Task{})
	model.Db = db

	go broadcaster.ServeTelegram(model.Db, GlobCfg.TG_KEY)
	manager.Init()

	httphandler.InitSession()

	mux.GET("/api", httphandler.Pong)

	mux.GET("/api/status", httphandler.GetStatus)
	mux.GET("/api/status/worker", httphandler.GetWorkerStatus)
	mux.GET("/api/status/node", httphandler.GetNodeStatus)
	mux.POST("/api/status/worker/:name", httphandler.HandleHeartbeat)

	mux.GET("/api/task", httphandler.GetTasks)
	mux.GET("/api/task/:id", httphandler.GetTask)
	mux.GET("/api/task/:id/log", httphandler.GetTaskLog)

	mux.POST("/api/task", httphandler.NewTask)
	mux.PUT("/api/task/:id/assign", httphandler.AssignTask)
	mux.PUT("/api/task/:id", httphandler.SyncTaskStatus)
	mux.DELETE("/api/task/:id", httphandler.ResetTask)
	mux.POST("/api/task/:id/callback", httphandler.TaskCallback)

	mux.POST("/api/broadcast", httphandler.Broadcast)

	mux.POST("/api/auth", httphandler.Login)
	mux.DELETE("/api/auth", httphandler.Logout)

	mux.GET("/api/node", httphandler.GetNodes)
	mux.GET("/api/node/:id", httphandler.GetNode)
	mux.POST("/api/node", httphandler.NewNode)
	mux.PUT("/api/node/:id", httphandler.EditNode)
	mux.PUT("/api/node/:id/enable/:name", httphandler.SetNodeTaskEnable)
	mux.DELETE("/api/node/:id/enable/:name", httphandler.SetNodeTaskDisable)
	mux.DELETE("/api/node/:id", httphandler.DeleteNode)
	mux.DELETE("/api/node/:id/status/isCleaning", httphandler.ResetNode)

	//mux.ServeFiles("/static/*filepath", http.Dir("static"))

	mux.NotFound = http.HandlerFunc(NotFoundHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   GlobCfg.ALLOW_ORIGIN,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowCredentials: true,
		//AllowedHeaders: []string{""},
	})
	handler := c.Handler(mux)

	n := negroni.New()
	n.Use(negroni.NewStatic(http.Dir("app")))
	n.UseHandler(handler)

	http.ListenAndServe(":"+strconv.FormatInt(GlobCfg.PORT, 10), n)

}

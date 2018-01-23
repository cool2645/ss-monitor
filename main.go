package main

import (
	"github.com/BurntSushi/toml"
	"net/http"
	"github.com/urfave/negroni"
	"strconv"
	"github.com/julienschmidt/httprouter"
	"github.com/astaxie/beego/session"
	"github.com/rs/cors"
	"github.com/cool2645/ss-monitor/broadcaster"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yanzay/log"
	"github.com/cool2645/ss-monitor/model"
	"github.com/cool2645/ss-monitor/httphandler"
)

var mux = httprouter.New()
var globalSessions *session.Manager

func main() {

	_, err := toml.DecodeFile("config.toml", &GlobCfg)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open("mysql", parseDSN(GlobCfg))
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Database init done")
	defer db.Close()

	db.AutoMigrate(&model.Subscriber{}, &model.Heartbeat{}, &model.Node{}, &model.Task{})
	model.Db = db

	go broadcaster.ServeTelegram(model.Db, GlobCfg.TG_KEY, broadcaster.Ch)

	mux.GET("/api", httphandler.Pong)

	mux.GET("/api/status", httphandler.Pong)
	mux.GET("/api/status/worker", httphandler.Pong)
	mux.GET("/api/status/node", httphandler.Pong)
	mux.POST("/api/status/worker/:id", httphandler.Pong)

	mux.GET("/api/task", httphandler.Pong)
	mux.GET("/api/task/:id", httphandler.Pong)
	mux.GET("/api/task/:id/log", httphandler.Pong)
	// Needs middleware here
	mux.POST("/api/task", httphandler.Pong)
	mux.PUT("/api/task/:id/assign", httphandler.Pong)
	mux.PUT("/api/task/:id", httphandler.Pong)
	mux.DELETE("/api/task/:id", httphandler.Pong)

	mux.POST("/api/broadcast", httphandler.Broadcast)

	//mux.ServeFiles("/static/*filepath", http.Dir("static"))

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

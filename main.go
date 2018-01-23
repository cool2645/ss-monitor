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

	go broadcaster.ServeTelegram(GlobCfg.TG_KEY)

	mux.GET("/api", Pong)

	//mux.ServeFiles("/static/*filepath", http.Dir("static"))

	c := cors.New(cors.Options{
		AllowedOrigins: GlobCfg.ALLOW_ORIGIN,
		AllowedMethods: []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowCredentials: true,
		//AllowedHeaders: []string{""},
	})
	handler := c.Handler(mux)

	n := negroni.New()
	n.Use(negroni.NewStatic(http.Dir("app")))
	n.UseHandler(handler)

	http.ListenAndServe(":"+strconv.FormatInt(GlobCfg.PORT, 10), n)

}
package main

import (
	"fmt"
	"gin-example/models"
	"gin-example/pkg/setting"
	"gin-example/routers"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	// log.SetFormatter(&log.JSONFormatter{})

	setting.Setup()
	models.SetUp()

	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.ServerSetting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}

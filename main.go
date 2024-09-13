package main

import (
	"fmt"
	"gin-example/models"
	"gin-example/pkg/setting"
	"gin-example/routers"
	"net/http"
	"time"

	"github.com/go-co-op/gocron/v2"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	router := routers.InitRouter()

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal(err)
	}

  //清除tag
	jCleanTag, err := scheduler.NewJob(
		gocron.DurationJob(10*time.Second),
		gocron.NewTask(models.CleanAllTag),
	)
	if err != nil {
		log.Fatal(err)
	}
  log.Info("[cleaning tag]:", jCleanTag.ID())

  //清除文章
	jCleanArticle, err := scheduler.NewJob(
		gocron.DurationJob(10*time.Second),
		gocron.NewTask(models.CleanAllArticle),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("[cleaning article]:", jCleanArticle.ID())

	scheduler.Start()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}

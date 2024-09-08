package setting

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-ini/ini"
)

var (
	Cfg *ini.File

	RunMode string

	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	PageSize  int
	JwtSecret string
)

func init() {
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}

	LoadBase()
	LoadServer()
	LoadApp()
}

func LoadBase() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
	log.WithFields(log.Fields{
		"RunMode": RunMode,
	}).Info("base loaded")
}

func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}

	HTTPPort = sec.Key("HTTP_PORT").MustInt(8000)

	rt := sec.Key("READ_TIMEOUT").MustInt(60)
	ReadTimeout = time.Duration(rt) * time.Second

	wt := sec.Key("WRITE_TIMEOUT").MustInt(60)
	WriteTimeout = time.Duration(wt) * time.Second

	defer log.WithFields(log.Fields{
		"HTTPPort":     HTTPPort,
		"ReadTimeout":  ReadTimeout,
		"WriteTimeout": WriteTimeout,
	}).Info("server loaded")
}

func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}

	JwtSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
	PageSize = sec.Key("PAGE_SIZE").MustInt(10)

	defer log.WithFields(log.Fields{
		"JwtSecret": JwtSecret,
		"PageSize":  PageSize,
	}).Info("app loaded")
}

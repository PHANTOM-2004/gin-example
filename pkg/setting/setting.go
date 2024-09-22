package setting

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-ini/ini"
)

var cfg *ini.File

type App struct {
	JwtSecret       string
	PageSize        int
	RuntimeRootPath string

	ImagePrefixUrl string
	ImageSavePath  string
	ImageMaxSize   int64
	ImageAllowExts []string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

type Server struct {
	RunMode      string
	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var (
	ServerSetting   = &Server{}
	AppSetting      = &App{}
	DatabaseSetting = &Database{}
	RedisSetting    = &Redis{}
)

func Setup() {
	var err error
	// load the config file
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}

	{
		// app setting
		err = cfg.Section("app").MapTo(AppSetting)
		if err != nil {
			log.Fatal(err)
		}
		AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024 // MB to B
		// log
		log.WithField(
			"app setting", AppSetting,
		).Info("")
	}

	{
		// server setting
		err = cfg.Section("server").MapTo(ServerSetting)
		if err != nil {
			log.Fatal(err)
		}
		ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
		ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second

		// log
		log.WithField(
			"server setting", ServerSetting,
		).Info("")

	}

	{
		// db setting
		err = cfg.Section("database").MapTo(DatabaseSetting)
		if err != nil {
			log.Fatal(err)
		}

		// log
		log.WithField(
			"db setting", DatabaseSetting,
		).Info("")

	}

	{
		// redis setting
		err = cfg.Section("redis").MapTo(RedisSetting)
		if err != nil {
			log.Fatal(err)
		}
		RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second

		// log
		log.WithField("redis setting", RedisSetting).Info("")
	}
}

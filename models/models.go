package models

import (
	"fmt"
	"gin-example/pkg/setting"
	"time"

	"github.com/go-co-op/gocron/v2"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/soft_delete"
)

var db *gorm.DB

type Model struct {
	ID         int                   `gorm:"primary_key" json:"id"`
	CreatedOn  int                   `json:"created_on"`
	ModifiedOn int                   `json:"modified_on"`
	DeletedOn  soft_delete.DeletedAt `json:"deleted_on"`
}

/*
CREATE TABLE `blog_tag` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) DEFAULT '' COMMENT '标签名称',
  `created_on` int(10) unsigned DEFAULT '0' COMMENT '创建时间',
  `created_by` varchar(100) DEFAULT '' COMMENT '创建人',
  `modified_on` int(10) unsigned DEFAULT '0' COMMENT '修改时间',
  `modified_by` varchar(100) DEFAULT '' COMMENT '修改人',
  `deleted_on` int(10) unsigned DEFAULT '0',
  `state` tinyint(3) unsigned DEFAULT '1' COMMENT '状态 0为禁用、1为启用',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='文章标签管理';

 alter table blog_article add cover_image_url
 varchar(255) DEFAULT '' COMMENT '封面图片地址';

CREATE TABLE `blog_article` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `tag_id` int(10) unsigned DEFAULT '0' COMMENT '标签ID',
  `title` varchar(100) DEFAULT '' COMMENT '文章标题',
  `desc` varchar(255) DEFAULT '' COMMENT '简述',
  `content` text,
  `created_on` int(11) DEFAULT NULL,
  `created_by` varchar(100) DEFAULT '' COMMENT '创建人',
  `modified_on` int(10) unsigned DEFAULT '0' COMMENT '修改时间',
  `modified_by` varchar(255) DEFAULT '' COMMENT '修改人',
  `deleted_on` int(10) unsigned DEFAULT '0',
  `state` tinyint(3) unsigned DEFAULT '1' COMMENT '状态 0为禁用1为启用',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='文章管理';

CREATE TABLE `blog_auth` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) DEFAULT '' COMMENT '账号',
  `password` varchar(50) DEFAULT '' COMMENT '密码',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `blog`.`blog_auth` (`id`, `username`, `password`) VALUES (null, 'test', 'test123456');
*/

func cleanHook() {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal(err)
	}

	// 清除tag
	jCleanTag, err := scheduler.NewJob(
		gocron.DurationJob(120*time.Second),
		gocron.NewTask(CleanAllTag),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("[cleaning tag]:", jCleanTag.ID())

	// 清除文章
	jCleanArticle, err := scheduler.NewJob(
		gocron.DurationJob(120*time.Second),
		gocron.NewTask(CleanAllArticle),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("[cleaning article]:", jCleanArticle.ID())

	scheduler.Start()
}

func SetUp() {
	defer log.Info("model initialized")

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name,
	)
	log.Info(dsn, setting.DatabaseSetting.Type)

	var err error

	db, err = gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{
			// 设置默认的db table handler
			NamingStrategy: schema.NamingStrategy{
				// table name prefix, table for `User` would be `t_users`
				TablePrefix: setting.DatabaseSetting.TablePrefix,
				// use singular table name, table for `User` would be `user` with this option enabled
				SingularTable: true,
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	sqldb, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqldb.SetMaxIdleConns(10)
	sqldb.SetMaxOpenConns(100)

	// 设置回调函数
	// https://github.com/go-gorm/gorm/blob/master/callbacks/callbacks.go
	// 注意这里是before, 而不是after, after都插入完毕了
	db.Callback().Create().Before("gorm:create").Register("my_plug:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Before("gorm:update").Replace("my_plug:update_time_stamp", updateTimeStampForUpdateCallback)

	cleanHook()
}

func CloseDB() {
	sqldb, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqldb.Close()
}

// 我们需要callback方法, 比如创建时间这种, 我们不可能为所有的
// 类都一个一个编写一个hook去更新时间
// 这里的callback是为了create和modify进行设计的
func updateTimeStampForCreateCallback(db *gorm.DB) {
	fCreatedOn := db.Statement.Schema.LookUpField("CreatedOn")
	fModifiedOn := db.Statement.Schema.LookUpField("ModifiedOn")
	if fCreatedOn == nil && fModifiedOn == nil {
		// 不存在这两个标签直接返回
		log.Debug("not found CreatedOn and ModifiedOn")
		return
	}

	nowTime := time.Now().Unix()

	if fCreatedOn != nil {

		// 设置创建时间
		_, isZero := fCreatedOn.ValueOf(db.Statement.Context, db.Statement.ReflectValue)
		log.Debug("Find CreatedOn, isZero=", isZero)
		if isZero {
			// 设置为当前时间
			err := fCreatedOn.Set(db.Statement.Context, db.Statement.ReflectValue, nowTime)
			if err != nil {
				log.Fatal(err)
			}
			log.Debug("set time:", nowTime)
		}
	}

	if fModifiedOn != nil {
		// 设置修改时间
		_, isZero := fModifiedOn.ValueOf(db.Statement.Context, db.Statement.ReflectValue)
		log.Debug("Find ModifiedOn, isZero=", isZero)
		if isZero {
			db.Statement.SetColumn("ModifiedOn", nowTime)
			log.Debug("set time:", nowTime)
		}
	}
}

func updateTimeStampForUpdateCallback(db *gorm.DB) {
	// 如果找不到会返回nil
	f := db.Statement.Schema.LookUpField("ModifiedOn")
	if f == nil {
		// 不存在这个那么就无需更新
		log.Debug("modified_on callback: no need to update")
		return
	}

	// 存在这一栏的话就需要更新
	// Get value from field
	_, isZero := f.ValueOf(db.Statement.Context, db.Statement.ReflectValue)
	log.Debug("Find ModifiedOn, isZero=", isZero)

	if !isZero {
		// 如果不是零值, 那么就不需要更改
		return
	}

	// Set value to field
	nowTime := time.Now().Unix()
	err := f.Set(db.Statement.Context, db.Statement.ReflectValue, nowTime)
	if err != nil {
		log.Fatal(err)
	}
}

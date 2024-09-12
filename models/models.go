package models

import (
	"fmt"
	"gin-example/pkg/setting"
	"time"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
}

func init() {
	defer log.Info("model initialized")

	var (
		err         error
		dbType      string
		dbName      string
		user        string
		password    string
		host        string
		tablePrefix string
	)

	sec, err := setting.Cfg.GetSection("database")
	if err != nil {
		log.Fatal(2, "Fail to get section 'database': %v", err)
	}

	// 从配置文件中读取数据库信息
	dbType = sec.Key("TYPE").String()
	dbName = sec.Key("NAME").String()
	user = sec.Key("USER").String()
	password = sec.Key("PASSWORD").String()
	host = sec.Key("HOST").String()
	tablePrefix = sec.Key("TABLE_PREFIX").String()

	db_info := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user,
		password,
		host,
		dbName)
	log.Info(db_info)

	db, err = gorm.Open(dbType, db_info)
	if err != nil {
		log.Fatal(err)
	}

	// 设置默认的db table handler
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}

	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

  //设置回调函数
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
}

func CloseDB() {
	defer db.Close()
}

// 我们需要callback方法, 比如创建时间这种, 我们不可能为所有的
// 类都一个一个编写一个hook去更新时间
// 这里的callback是为了create和modify进行设计的
func updateTimeStampForCreateCallback(s *gorm.Scope) {
	if s.HasError() {
		log.Fatal("update/create callback: db error")
		return
	}
	/*
		for _, field := range scope.Fields() {
		    if field.Name == name || field.DBName == name {
		        return field, true
		    }
		    if field.DBName == dbName {
		        mostMatchedField = field
		    }
	*/

	nowTime := time.Now().Unix()
	createTimeField, ok := s.FieldByName("CreatedOn")
	/*
		注意这里的空值说的是什么意思
				func isBlank(value reflect.Value) bool {
					switch value.Kind() {
					case reflect.String:
						return value.Len() == 0
					case reflect.Bool:
						return !value.Bool()
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						return value.Int() == 0
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
						return value.Uint() == 0
					case reflect.Float32, reflect.Float64:
						return value.Float() == 0
					case reflect.Interface, reflect.Ptr:
						return value.IsNil()
					}

					return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
				}
	*/

	if ok && createTimeField.IsBlank {
		// 前提是有这个栏
		// 如果这个栏是blank的时候我么进行更新
		createTimeField.Set(nowTime)
	}

	modifyTimeField, ok := s.FieldByName("ModifiedOn")
	if ok && modifyTimeField.IsBlank {
		modifyTimeField.Set(nowTime)
	}
}

func updateTimeStampForUpdateCallback(s *gorm.Scope) {
	_, ok := s.Get("gorm:update_column")

	if !ok {
		s.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

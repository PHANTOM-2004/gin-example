package gredis

import (
	"context"
	"encoding/json"
	"gin-example/pkg/setting"
	"time"

	redis "github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var Rdb *redis.Client

func Setup() error {
	Rdb = redis.NewClient(&redis.Options{
		Addr:           setting.RedisSetting.Host,
		Password:       setting.RedisSetting.Password,
		MaxActiveConns: setting.RedisSetting.MaxActive,
		DB:             0,
	})

	return nil
}

func Set(key string, data interface{}, expiresTime int) (bool, error) {
	value, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	// TODO: should context always be background
	ctx := context.Background()

	// TODO: this should not be zero(no limit)
	err = Rdb.Set(ctx, key, value, time.Second*time.Duration(expiresTime)).Err()
	if err != nil {
		return false, err
	}

	return true, err
}

func Exists(key string) bool {
	ctx := context.Background()
	res, err := Rdb.Exists(ctx, key).Result()
	if err != nil {
		return false
	}

	return res != 0
}

func Get(key string) ([]byte, error) {
	ctx := context.Background()
	value, err := Rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return []byte(value), nil
}

func Delete(key string) (bool, error) {
	ctx := context.Background()
	res, err := Rdb.Del(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return res != 0, nil
}

func LikeDeletes(key string) error {
	ctx := context.Background()
	strs, err := Rdb.Keys(ctx, "*"+key+"*").Result()
	if err != nil {
		return err
	}

	for _, s := range strs {
		log.Debug("delete: ", s)
		_, err := Rdb.Del(ctx, s).Result()
		if err != nil {
			return err
		}
	}
	return nil
}

package redisdb

import (
	"common"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

const AllGoodsKey = "all_goods"

var RedisClient *redis.Client

type connectionParams struct {
	Addr     string
	Password string
	DB       int
}

func OpenConnection(connParams connectionParams) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     connParams.Addr,
		Password: connParams.Password,
		DB:       connParams.DB,
	})

	_, err := RedisClient.Ping().Result()
	if err != nil {
		logrus.Errorf("error pinging redis connection [%s]", err.Error())
		return err
	}

	return nil
}

func GetConnectionParams() (connectionParams, error) {
	host, err := common.GetEnvVar("REDIS_HOST")
	if err != nil {
		logrus.Errorf("error getting REDIS_HOST var from env [%s]", err.Error())
		return connectionParams{}, err
	}

	port, err := common.GetEnvVar("REDIS_PORT")
	if err != nil {
		logrus.Errorf("error getting REDIS_PORT var from env [%s]", err.Error())
		return connectionParams{}, err
	}

	return connectionParams{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "",
		DB:       0,
	}, nil
}

func Cache(key string, data interface{}) error {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		logrus.Errorf("error marshal data [%s]", err.Error())
		return err
	}

	return RedisClient.Set(key, dataJSON, time.Minute).Err()
}

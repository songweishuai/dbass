package myRedis

import (
	"errors"
	"github.com/gomodule/redigo/redis"
)

func GetRedisInfo(key string) (string, error) {
	if key == "" {
		return "", nil
	}

	//con, err := GetRedisInstance()
	pool := GetPool()
	if pool == nil {
		err := errors.New("pool is nil")
		return "", err
	}
	con := pool.Get()
	if con == nil {
		err := errors.New("con got from pool is nil")
		return "", err
	}

	exist, err := con.Do("exists", key)
	if err != nil {
		return "", err
	}
	if exist != int64(1) {
		err = errors.New("key do not exist")
		return "", err
	}

	str, err := con.Do("get", key)
	if err != nil {
		return "", err
	}
	ret, err := redis.String(str, err)
	if err != nil {
		return "", err
	}

	return ret, nil
}

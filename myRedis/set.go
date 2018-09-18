package myRedis

import (
	"errors"
	"fmt"
)

func SetRedisInfo(key string, info string) error {
	//con, err := GetRedisInstance()
	//if err != nil {
	//	return err
	//}

	pool := GetPool()
	if pool == nil {
		err := errors.New("pool is nil")
		return err
	}
	con := pool.Get()
	if con == nil {
		err := errors.New("con got from pool is nil")
		return err
	}

	exist, err := con.Do("exists", key)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if exist == int64(1) {
		err = errors.New("key exist")
		return err
	}

	_, err = con.Do("set", key, info)
	if err != nil {
		return err
	}

	return nil
}

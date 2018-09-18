package myRedis

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

var (
	pool *redis.Pool
)

func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 120 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func init() {
	pool = newPool("127.0.0.1:6379")
}

func GetPool() *redis.Pool {
	if pool == nil {
		pool = newPool("127.0.0.1:6379")
	}
	return pool
}

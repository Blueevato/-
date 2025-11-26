package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

var pool *redis.Pool //全局变量

func initPool(address string, maxIdle, maxActive, idleTimeout int) {
	pool = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(idleTimeout),
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address)
		},
	}
}

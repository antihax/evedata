package redigohelper

import (
	"errors"
	"log"
	"os"
	"time"

	sentinel "github.com/FZambia/go-sentinel"
	"github.com/garyburd/redigo/redis"
)

func ConnectRedisProdPool() *redis.Pool {
	return connectRedisPool(
		[]string{"sentinel1.evedata:26379", "sentinel2.evedata:26379", "sentinel3.evedata:26379"},
		os.Getenv("REDIS_PASSWORD"),
		"evedata",
		true,
	)
}

func ConnectRedisTestPool() *redis.Pool {
	redis := connectRedisPool(
		[]string{"127.0.0.1:6379"},
		"",
		"",
		false,
	)
	c := redis.Get()
	_, err := c.Do("FLUSHALL")
	if err != nil {
		panic(err)
	}
	err = c.Close()
	if err != nil {
		panic(err)
	}
	return redis
}
func connectRedisPool(addresses []string, password string, masterName string, sentinel bool) *redis.Pool {
	if sentinel {
		return newSentinelPool(addresses, masterName, password)
	} else {
		return newRedisPool(addresses[0], password)
	}
}
func newRedisPool(address string, password string) *redis.Pool {
	// Build the redis pool
	return &redis.Pool{
		MaxIdle:     200,
		MaxActive:   0,
		Wait:        false,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					log.Fatalln(err)
				}
			}
			return c, err
		},
	}
}

func newSentinelPool(addresses []string, masterName string, password string) *redis.Pool {
	sntnl := &sentinel.Sentinel{
		Addrs:      addresses,
		MasterName: masterName,
		Dial: func(addr string) (redis.Conn, error) {
			timeout := 500 * time.Millisecond
			c, err := redis.DialTimeout("tcp", addr, timeout, timeout, timeout)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			return c, nil
		},
	}

	return &redis.Pool{
		MaxIdle:     200,
		MaxActive:   0,
		Wait:        false,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			masterAddr, err := sntnl.MasterAddr()
			if err != nil {
				log.Println(err)
				return nil, err
			}
			c, err := redis.Dial("tcp", masterAddr)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					log.Fatalln(err)
					return nil, err
				}
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if !sentinel.TestRole(c, "master") {
				return errors.New("Role check failed")
			} else {
				return nil
			}
		},
	}
}

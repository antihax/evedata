package redigohelper

import (
	"errors"
	"log"
	"os"
	"time"

	sentinel "github.com/FZambia/go-sentinel"
	"github.com/garyburd/redigo/redis"
)

// Memory Store
func ConnectRedisProdPool() *redis.Pool {
	pool := connectRedisPool(
		[]string{"redis.storage.svc.cluster.local:6379"},
		os.Getenv("REDIS_PASSWORD"),
		"evedata",
		false,
	)
	test := pool.Get()
	_, err := test.Do("PING")
	if err != nil {
		log.Fatalln(err)
	}

	return pool
}

// Disk Store
func ConnectLedisProdPool() *redis.Pool {
	pool := connectRedisPool(
		[]string{"redis.storage.svc.cluster.local:6379"},
		os.Getenv("REDIS_PASSWORD"),
		"evedata",
		false,
	)
	test := pool.Get()
	_, err := test.Do("PING")
	if err != nil {
		log.Fatalln(err)
	}

	return pool
}

func ConnectLedisTestPool() *redis.Pool {
	redis := connectRedisPool(
		[]string{"127.0.0.1:6379"},
		"",
		"",
		false,
	)
	c := redis.Get()
	defer c.Close()

	_, err := c.Do("FLUSHALL")
	if err != nil {
		panic(err)
	}
	return redis
}
func ConnectRedisTestPool() *redis.Pool {
	redis := connectRedisPool(
		[]string{"127.0.0.1:6379"},
		"",
		"",
		false,
	)
	c := redis.Get()
	defer c.Close()

	_, err := c.Do("FLUSHALL")
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
		MaxIdle:     20,
		MaxActive:   120,
		Wait:        false,
		IdleTimeout: 20 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address,
				redis.DialReadTimeout(60*time.Second),
				redis.DialWriteTimeout(60*time.Second))
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
		MaxIdle:     20,
		MaxActive:   120,
		Wait:        false,
		IdleTimeout: 20 * time.Second,
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
				_, err := c.Do("PING")
				return err
			}
		},
	}
}

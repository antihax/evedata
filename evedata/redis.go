package evedata

import (
	"errors"
	"time"

	sentinel "github.com/FZambia/go-sentinel"
	"github.com/antihax/evedata/appContext"

	"github.com/garyburd/redigo/redis"
)

func setupRedis(ctx *appContext.AppContext) *redis.Pool {
	if ctx.Conf.Redis.Sentinel {
		return newSentinelPool(ctx.Conf.Redis.Addresses, ctx.Conf.Redis.MasterName, ctx.Conf.Redis.Password)
	} else {
		return newRedisPool(ctx.Conf.Redis.Address, ctx.Conf.Redis.Password)
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
					return nil, err
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
				return nil, err
			}
			c, err := redis.Dial("tcp", masterAddr)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
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

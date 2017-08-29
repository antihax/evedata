package redigohelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	pool := connectRedisPool([]string{"127.0.0.1:6379"}, "", "", false)
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("INFO")
	assert.Nil(t, err)
}

func TestNoConnect(t *testing.T) {
	pool := connectRedisPool([]string{"deadend:6379"}, "", "", false)
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("INFO")
	assert.NotNil(t, err)
}

func TestSentinelNoConnect(t *testing.T) {
	pool := connectRedisPool([]string{"deadend:6379"}, "", "", true)
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("INFO")
	assert.NotNil(t, err)
}

func TestRedisTestConnect(t *testing.T) {
	pool := ConnectRedisTestPool()
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("INFO")
	assert.Nil(t, err)
}

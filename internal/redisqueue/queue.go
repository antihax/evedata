package redisqueue

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/garyburd/redigo/redis"
)

// RedisQueue operation queue to CCP APIs
type RedisQueue struct {
	redisPool *redis.Pool
	key       string
}

// Work to be performed
type Work struct {
	Operation string      `json:"operation"`
	Parameter interface{} `json:"parameters"`
}

// NewRedisQueue creates a new work queue with an existing
// redigo pool and key name.
func NewRedisQueue(r *redis.Pool, key string) *RedisQueue {
	gob.Register([]interface{}{})
	return &RedisQueue{redisPool: r, key: key}
}

// Size returns number of elements in the queue
func (hq *RedisQueue) Size() (int, error) {
	r := hq.redisPool.Get()
	return redis.Int(r.Do("LLEN", hq.key))
}

// QueueWork adds work to the queue
func (hq *RedisQueue) QueueWork(work []Work) error {
	// Get a redis connection from the pool
	conn := hq.redisPool.Get()
	defer conn.Close()

	// Pipeline our work to the connection.
	for i := range work {
		var b bytes.Buffer
		enc := gob.NewEncoder(&b)

		err := enc.Encode(work[i])
		if err != nil {
			return err
		}
		if err := conn.Send("LPUSH", hq.key, b.Bytes()); err != nil {
			return err
		}
	}
	if err := conn.Flush(); err != nil {
		return err
	}

	return nil
}

// GetWork retreives up to `n` items from the queue
func (hq *RedisQueue) GetWork() (*Work, error) {
	// Get a redis connection from the pool
	conn := hq.redisPool.Get()
	defer conn.Close()

	var w Work
	v, err := conn.Do("BRPOP", hq.key, "5")
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, errors.New("Timed out")
	}

	results := v.([]interface{})

	b := bytes.NewBuffer(results[1].([]byte))
	dec := gob.NewDecoder(b)
	if err := dec.Decode(&w); err != nil {
		return nil, err
	}

	return &w, nil
}

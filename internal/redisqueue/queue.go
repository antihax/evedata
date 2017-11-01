package redisqueue

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"

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
	defer r.Close()
	return redis.Int(r.Do("SCARD", hq.key))
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
		if _, err := conn.Do("SADD", hq.key, b.Bytes()); err != nil {
			return err
		}
	}
	//	err := conn.Flush()
	return nil
}

// CheckWorkCompleted takes a key and checks if the ID has been completed to prevent duplicates
func (hq *RedisQueue) CheckWorkCompleted(key string, id int64) bool {
	conn := hq.redisPool.Get()
	defer conn.Close()
	found, err := redis.Bool(conn.Do("SISMEMBER", key, id))
	if err != nil {
		log.Println(err)
	}
	return found
}

// SetWorkCompleted takes a key and sets if the ID has been completed to prevent duplicates
func (hq *RedisQueue) SetWorkCompleted(key string, id int64) error {
	conn := hq.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SADD", key, id)
	return err
}

// CheckWorkExpired takes a key and checks if the ID has expired
func (hq *RedisQueue) CheckWorkExpired(key string, id int64) bool {
	conn := hq.redisPool.Get()
	defer conn.Close()
	found, _ := redis.Bool(conn.Do("GET", fmt.Sprintf("%s:%d", key, id)))
	return found
}

// SetWorkExpire takes a key and sets if the ID has failed to prevent multiple failed
func (hq *RedisQueue) SetWorkExpire(key string, id int64, seconds int) error {
	conn := hq.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SETEX", fmt.Sprintf("%s:%d", key, id), seconds, true)
	return err
}

// GetWork retreives up to `n` items from the queue
func (hq *RedisQueue) GetWork() (*Work, error) {
	// Get a redis connection from the pool
	conn := hq.redisPool.Get()
	defer conn.Close()

	var (
		w   Work
		v   interface{}
		err error
	)

	// Block until we get data.
	for {
		v, err = conn.Do("SPOP", hq.key)
		if err != nil || v == nil {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		break
	}

	// Decode the data back into its structure
	b := bytes.NewBuffer(v.([]byte))
	dec := gob.NewDecoder(b)
	if err := dec.Decode(&w); err != nil {
		return nil, err
	}

	return &w, nil
}

package redisqueue

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/antihax/evedata/internal/gobcoder"

	"github.com/garyburd/redigo/redis"
)

const (
	Priority_Low    = iota
	Priority_Normal = iota
	Priority_High   = iota
	Priority_Urgent = iota
)

// RedisQueue operation queue to CCP APIs
type RedisQueue struct {
	redisPool   *redis.Pool
	key         string
	queueScript *redis.Script
}

// Work to be performed
type Work struct {
	Operation string      `json:"operation"`
	Parameter interface{} `json:"parameters"`
}

// NewRedisQueue creates a new work queue with an existing
// redigo pool and key name.
func NewRedisQueue(r *redis.Pool, key string) *RedisQueue {

	rq := &RedisQueue{redisPool: r, key: key, queueScript: redis.NewScript(0, priorityQueueScript)}
	conn := r.Get()
	defer conn.Close()

	err := rq.queueScript.Load(conn)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Loaded script %s\n", rq.queueScript.Hash())

	return rq
}

// Size returns number of elements in the queue
func (hq *RedisQueue) Size() (int, error) {
	r := hq.redisPool.Get()
	defer r.Close()
	return redis.Int(hq.queueScript.Do(r, "count", hq.key))
}

// ScaledSize  returns number of elements in the queue scaled by priority
func (hq *RedisQueue) ScaledSize() (int, error) {
	r := hq.redisPool.Get()
	defer r.Close()

	low, err := redis.Int(hq.queueScript.Do(r, "count", hq.key, Priority_Low, Priority_Low))
	if err != nil {
		return 0, err
	}
	norm, err := redis.Int(hq.queueScript.Do(r, "count", hq.key, Priority_Normal, Priority_Normal))
	if err != nil {
		return 0, err
	}
	high, err := redis.Int(hq.queueScript.Do(r, "count", hq.key, Priority_High, Priority_High))
	if err != nil {
		return 0, err
	}
	urg, err := redis.Int(hq.queueScript.Do(r, "count", hq.key, Priority_Urgent, Priority_Urgent))
	if err != nil {
		return 0, err
	}

	return low + (norm * 64) + (high * 1024) + (urg * 8192), nil
}

// QueueWork adds work to the queue
func (hq *RedisQueue) QueueWork(work []Work, priority int) error {
	// Get a redis connection from the pool
	conn := hq.redisPool.Get()
	defer conn.Close()

	// Pipeline our work to the connection.
	for i := range work {
		b, err := gobcoder.GobEncoder(work[i])
		if err != nil {
			return err
		}
		if _, err := hq.queueScript.Do(conn, "push", hq.key, b, priority); err != nil {
			return err
		}
	}

	return nil
}

// GetWork retreives items from the queue
func (hq *RedisQueue) GetWork() (*Work, error) {
	// Get a redis connection from the pool
	conn := hq.redisPool.Get()
	defer conn.Close()

	var (
		w   Work
		v   interface{}
		err error
	)

	// Poll until we get data.
	for {
		v, err = hq.queueScript.Do(conn, "pop", hq.key)
		if err != nil || v == nil {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		break
	}

	b, ok := v.([]byte)
	if !ok {
		return nil, errors.New("empty result")
	}

	// Decode the data back into its structure
	if err := gobcoder.GobDecoder(b, &w); err != nil {
		return nil, err
	}

	return &w, nil
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

// CheckWorkCompletedInBulk takes a key and checks if the list of IDs has completed
func (hq *RedisQueue) CheckWorkCompletedInBulk(key string, id []int64) ([]bool, error) {
	conn := hq.redisPool.Get()
	defer conn.Close()
	var found []bool

	conn.Send("MULTI")
	for _, i := range id {
		conn.Send("SISMEMBER", key, i)
	}

	res, err := conn.Do("EXEC")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for _, r := range res.([]interface{}) {
		b, _ := redis.Bool(r, nil)
		found = append(found, b)
	}

	return found, nil
}

// SetWorkCompleted takes a key and sets if the ID has been completed to prevent duplicates
func (hq *RedisQueue) SetWorkCompleted(key string, id int64) error {
	conn := hq.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SADD", key, id)
	return err
}

// SetWorkCompleted takes a key and sets if the ID has been completed to prevent duplicates
func (hq *RedisQueue) SetWorkCompletedInBulk(key string, ids []int64) error {
	conn := hq.redisPool.Get()
	defer conn.Close()

	for _, id := range ids {
		err := conn.Send("SADD", key, id)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	err := conn.Flush()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// CheckWorkExpired takes a key and checks if the ID has expired
func (hq *RedisQueue) CheckWorkExpired(key string, id int64) bool {
	conn := hq.redisPool.Get()
	defer conn.Close()
	found, _ := redis.Bool(conn.Do("GET", fmt.Sprintf("%s:%d", key, id)))
	return found
}

// CheckWorkExpiredInBulk takes a key and checks if the list of IDs has expired
func (hq *RedisQueue) CheckWorkExpiredInBulk(key string, id []int64) ([]bool, error) {
	conn := hq.redisPool.Get()
	defer conn.Close()
	var found []bool

	conn.Send("MULTI")
	for _, i := range id {
		conn.Send("GET", fmt.Sprintf("%s:%d", key, i))
	}

	res, err := conn.Do("EXEC")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for _, r := range res.([]interface{}) {
		b, _ := redis.Bool(r, nil)
		found = append(found, b)
	}

	return found, nil
}

// SetWorkExpire takes a key and sets if the ID has failed to prevent multiple failed
func (hq *RedisQueue) SetWorkExpire(key string, id int64, seconds int) error {
	conn := hq.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SETEX", fmt.Sprintf("%s:%d", key, id), seconds, true)
	return err
}

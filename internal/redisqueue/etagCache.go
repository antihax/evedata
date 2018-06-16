package redisqueue

import (
	"log"

	"github.com/garyburd/redigo/redis"
)

// KnownETag takes an etag and checks if it is known.
func (hq *RedisQueue) KnownETag(key string, etag string) bool {
	conn := hq.redisPool.Get()
	defer conn.Close()
	found, err := redis.Bool(conn.Do("SISMEMBER", key, etag))
	if err != nil {
		log.Println(err)
	}
	return found
}

// SetETagKnown takes a etag and sets if the tag has been seen
func (hq *RedisQueue) SetETagKnown(key string, etag string) error {
	conn := hq.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SADD", key, etag)
	return err
}

// SetETagKnownInBulk takes many etags and sets the known state
func (hq *RedisQueue) SetETagKnownInBulk(key string, etags []string) error {
	conn := hq.redisPool.Get()
	defer conn.Close()

	for _, etag := range etags {
		err := conn.Send("SADD", key, etag)
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

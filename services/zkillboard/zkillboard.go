package zkillboard

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/antihax/evedata/internal/apicache"
	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/garyburd/redigo/redis"
)

// NewZKillboard sucks down killmails from zkillboard redisq and API.
type ZKillboard struct {
	stop     chan bool
	wg       *sync.WaitGroup
	outQueue *redisqueue.RedisQueue
	redis    *redis.Pool
	http     *http.Client
}

// NewZKillboard Service.
func NewZKillboard(redis *redis.Pool) *ZKillboard {
	// Setup a new hammer
	s := &ZKillboard{
		stop: make(chan bool),
		wg:   &sync.WaitGroup{},
		outQueue: redisqueue.NewRedisQueue(
			redis,
			"evedata-hammer",
		),
		http:  apicache.CreateHTTPClientCache(),
		redis: redis,
	}
	return s
}

// Close the service
func (s *ZKillboard) Close() {
	close(s.stop)
	s.wg.Wait()
}

// Run the service
func (s *ZKillboard) Run() {
	// Run the API consumer in a loop
	go s.apiConsumer()

	for {
		select {
		case <-s.stop:
			return
		default:
			// Pull ReqisQ
			err := s.redisQ()
			if err != nil {
				log.Println(err)

				// Back off
				time.Sleep(time.Second * 5)
			}
		}
	}
}

// getJSON into a struct
func (s *ZKillboard) getJSON(url string, v interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("user-agent", "EVEData.org ZKill Consumer (nom nom)")
	r, err := s.http.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(v)
}

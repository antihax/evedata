package artifice

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type killmail struct {
	ID   int32
	Hash string
}

var zkillChan chan killmail

// zkillboardPost posts killmails to zkillboard from zkillChan
func (s *Artifice) zkillboardPost() {
	// Create the channel for feeding kills
	zkillChan = make(chan killmail, 100)

	// Need a http client for keep-alive/http2
	httpClient := http.Client{}

	// Throttle requests
	throttle := time.Tick(time.Second / 9)

	for {
		// pop a killmail off the channel
		k := <-zkillChan
		if k.ID < 100 {
			continue
		}
		mail := fmt.Sprintf("https://zkillboard.com/crestmail/%d/%s/", k.ID, k.Hash)

		r, _ := http.NewRequest("GET", mail, nil)
		r.Header.Add("Content-Type", "text/text")
		r.Header.Set("User-Agent", "EVEData.org - from croakroach with love.")
		resp, _ := httpClient.Do(r)

		log.Printf("Posted to Zkillboard %s %s\n", resp.Status, mail)

		// Don't hammer zkillboard
		<-throttle
	}
}

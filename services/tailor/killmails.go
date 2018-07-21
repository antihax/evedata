package tailor

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/evedata/internal/gobcoder"

	nsq "github.com/nsqio/go-nsq"
)

func (s *Tailor) killmailHandler(message *nsq.Message) error {
	killmail := datapackages.Killmail{}
	if err := gobcoder.GobDecoder(message.Body, &killmail); err != nil {
		log.Println(err)
		return err
	}
	mail, _ := json.Marshal(killmail.Kill)
	ioutil.WriteFile("./json/"+killmail.Hash+".json", mail, 0644)
	return nil
}

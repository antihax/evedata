package conservator

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"log"

	"github.com/antihax/evedata/internal/botservice"
	"github.com/antihax/evedata/internal/botservice/discordservice"
	"github.com/antihax/evedata/internal/botservice/tsservice"
)

type ChannelOptions struct {
	Killmail struct {
		IgnoreHighSec   bool `json:"ignoreHighsec,omitempty"`
		IgnoreLowSec    bool `json:"ignoreLowsec,omitempty"`
		IgnoreNullSec   bool `json:"ignoreNullsec,omitempty"`
		IgnoreWorthless bool `json:"ignoreWorthless,omitempty"`
		War             bool `json:"war,omitempty"`
	} `json:"killmail,omitempty"`
}

type Service struct {
	Server         botservice.BotService
	BotServiceID   int32  `db:"botServiceID"`
	EntityID       int32  `db:"entityID"`
	Address        string `db:"address"`
	Authentication string `db:"authentication"`
	Type           string `db:"type"`
	Services       string `db:"services"`
}

type Channel struct {
	BotServiceID int32          `db:"botServiceID"`
	ChannelID    string         `db:"channelID"`
	Services     string         `db:"services"`
	OptionsJSON  string         `db:"options"`
	Options      ChannelOptions `db:"-"`
}

type Share struct {
	BotServiceID     int32  `db:"botServiceID"`
	TokenCharacterID int32  `db:"tokenCharacterID"`
	EntityID         int32  `db:"entityID"`
	Types            string `db:"types"`
	Packed           string `db:"packed"`
}

// Load our bot services
func (s *Conservator) loadServices() error {
	// Mark what we touch so we can erase any missing items
	touched := make(map[int32]bool)

	services, err := s.getServices()
	if err != nil {
		return err
	}

	// [TODO] Slack, murmur
	for _, service := range services {
		var n botservice.BotService

		// Mark what we touch
		touched[service.BotServiceID] = true

		switch service.Type {
		case "discord":
			n = discordservice.NewDiscordService(s.discord, service.Address)
		case "ts3":
			auth := strings.Split(service.Authentication, ":")
			n, err = tsservice.NewTSService(service.Address, auth[0], auth[1])
			if err != nil {
				log.Println(err)
			}
		default:
			return errors.New("unknown service type")
		}

		// Store the server with the service information and put into our map
		service.Server = n
		s.services.Store(service.BotServiceID, service)
	}

	// Delete any items that were removed
	s.services.Range(func(ki, vi interface{}) bool {
		k := ki.(int32)
		v := vi.(Service)

		if !touched[v.BotServiceID] {
			s.services.Delete(k)
		}
		return true
	})

	return nil
}

// loadChannels will load any bot channels
// running a second time will add any new items and delete removed.
func (s *Conservator) loadChannels() error {
	// Mark what we touch so we can erase any missing items
	touched := make(map[string]bool)

	// Load channels
	channels, err := s.getChannels()
	if err != nil {
		return err
	}

	// Add channels to map
	for _, channel := range channels {
		// Mark what we touch
		touched[channel.ChannelID] = true
		json.Unmarshal([]byte(channel.OptionsJSON), &channel.Options)

		s.channels.Store(channel.ChannelID, channel)
	}

	// Delete any items that were removed
	s.channels.Range(func(ki, vi interface{}) bool {
		k := ki.(string)
		v := vi.(Channel)

		if !touched[v.ChannelID] {
			s.services.Delete(k)
		}
		return true
	})

	return nil
}

// loadShares will load any characters sharing data with entities
// running a second time will add any new items and delete removed.
func (s *Conservator) loadShares() error {
	// Mark what we touch so we can erase any missing items
	touched := make(map[string][]string)

	// Load character to entity shares
	shares, err := s.getShares()
	if err != nil {
		return err
	}

	// Add new shares for services
	for _, share := range shares {
		// Add all the shares
		services := strings.Split(share.Types, ",")

		// unpack channel data
		channelData := unpackChannelData(share.Packed)

		// Touch the entry
		touched[fmt.Sprint(share.TokenCharacterID, share.EntityID)] = services

		// Add the services
		for _, service := range services {
			// Should we be sending data to a channel for this service?
			found := false
			for _, channel := range channelData {
				if inSlice(service, channel.Services) {
					found = true
					break
				}
			}
			if !found {
				continue
			}

			s.notificationLock[service].Lock()
			a := s.notifications[service][share.TokenCharacterID]

			// Don't duplicate
			found = false
			for i := range a {
				if a[i].TokenCharacterID == share.TokenCharacterID && a[i].BotServiceID == share.BotServiceID {
					found = true
					break
				}
			}
			if !found {
				a = append(a, share)
				s.notifications[service][share.TokenCharacterID] = a
			}
			s.notificationLock[service].Unlock()
		}
	}

	// Delete any items that were removed
	for _, service := range NOTIFICATION_TYPES {
		s.notificationLock[service].Lock()
		for charID := range s.notifications[service] {
			for j, share := range s.notifications[service][charID] {
				entry, ok := touched[fmt.Sprint(charID, share.EntityID)]
				if !ok || !inSlice(service, entry) {
					s.notifications[service][charID] = append(s.notifications[service][charID][:j], s.notifications[service][charID][j+1:]...)
					if len(s.notifications[service][charID]) == 0 {
						delete(s.notifications[service], charID)
					}
				}
			}
		}
		s.notificationLock[service].Unlock()
	}

	return nil
}

func (s *Conservator) getServices() ([]Service, error) {
	services := []Service{}
	err := s.db.Select(&services, "SELECT botServiceID, entityID, address, authentication, type, services FROM evedata.botServices")
	if err != nil {
		return nil, err
	}
	return services, nil
}

func (s *Conservator) getChannels() ([]Channel, error) {
	channels := []Channel{}
	err := s.db.Select(&channels, "SELECT botServiceID, channelID, services, options FROM evedata.botChannels")
	if err != nil {
		return nil, err
	}
	return channels, nil
}

func (s *Conservator) updateData() {
	throttle := time.Tick(time.Second * 30)
	for {
		s.updateWars()
		if err := s.loadServices(); err != nil {
			log.Println(err)
		}
		if err := s.loadChannels(); err != nil {
			log.Println(err)
		}
		if err := s.loadShares(); err != nil {
			log.Println(err)
		}
		<-throttle
	}
}

func (s *Conservator) updateWars() {
	s.services.Range(func(ki, vi interface{}) bool {
		service := vi.(Service)
		warlist := []atWarWith{}
		err := s.db.Select(&warlist, "CALL evedata.atWarWith(?);", service.EntityID)
		if err != nil {
			log.Println(err)
		}

		for _, war := range warlist {
			if s.warsMap[service.EntityID] == nil {
				s.warsMap[service.EntityID] = &sync.Map{}
			}
			s.warsMap[service.EntityID].Store(war.ID, war)
		}
		return true
	})
}

type ChannelData struct {
	ChannelID string
	Services  []string
}

// unpackChannelData from our getShares query
func unpackChannelData(packed string) []ChannelData {
	c := []ChannelData{}
	channels := strings.Split(packed, ":")
	for _, channel := range channels {
		info := strings.Split(channel, "|")
		c = append(c, ChannelData{
			ChannelID: info[0],
			Services:  strings.Split(info[1], ","),
		})
	}
	return c
}

// getShares uses an ulgy hack to deal with comparing sets.
// Compare the sets localy after depacking
// channel1|services,for,channel:channel2|services,for,channel
func (s *Conservator) getShares() ([]Share, error) {
	shares := []Share{}
	err := s.db.Select(&shares, `
		SELECT B.botServiceID, tokenCharacterID, S.entityID, types, group_concat(concat(channelID, "|", C.services) SEPARATOR ":") as packed
		FROM evedata.sharing S
		INNER JOIN evedata.botServices B ON B.entityID = S.entityID
        INNER JOIN evedata.botChannels C ON C.botServiceID = B.botServiceID 
        GROUP BY botServiceID, tokenCharacterID, entityID`)
	if err != nil {
		return nil, err
	}
	return shares, nil
}

func (s *Conservator) getSolarSystems() (map[int32]float32, error) {
	type system struct {
		SolarSystemID int32   `db:"solarSystemID"`
		Security      float32 `db:"security"`
	}
	systems := []system{}
	err := s.db.Select(&systems, "SELECT solarSystemID, security FROM mapSolarSystems")
	if err != nil {
		return nil, err
	}

	solarSystems := make(map[int32]float32)
	for _, s := range systems {
		solarSystems[s.SolarSystemID] = s.Security
	}
	return solarSystems, nil
}

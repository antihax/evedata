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
		FactionWar      bool `json:"factionWar,omitempty"`
		SendAll         bool `json:"sendAll,omitempty"`
	} `json:"killmail,omitempty"`
}

type Service struct {
	Server         botservice.BotService `json:"-,omitempty"`
	BotServiceID   int32                 `db:"botServiceID" json:"botServiceID,omitempty"`
	Name           string                `db:"name" json:"name,omitempty"`
	EntityID       int32                 `db:"entityID" json:"entityID,omitempty"`
	EntityName     string                `db:"entityName" json:"entityName,omitempty"`
	EntityType     string                `db:"entityType" json:"entityType,omitempty"`
	Address        string                `db:"address" json:"address,omitempty" `
	Authentication string                `db:"authentication,omitempty"`
	Type           string                `db:"type" json:"type,omitempty"`
	Services       string                `db:"services" json:"services,omitempty"`
	OptionsJSON    string                `db:"options" json:"-,omitempty"`
}

type Channel struct {
	BotServiceID int32          `db:"botServiceID" json:"botServiceID,omitempty"`
	ChannelID    string         `db:"channelID"  json:"channelID,omitempty"`
	ChannelName  string         `db:"channelName"  json:"channelName,omitempty"`
	Services     string         `db:"services" json:"services,omitempty"`
	OptionsJSON  string         `db:"options" json:"-"`
	Options      ChannelOptions `db:"-" json:"options,omitempty"`
}

type Share struct {
	BotServiceID       int32  `db:"botServiceID" json:"botServiceID,omitempty"`
	CharacterID        int32  `db:"characterID" json:"characterID,omitempty"`
	TokenCharacterID   int32  `db:"tokenCharacterID" json:"tokenCharacterID,omitempty"`
	TokenCharacterName string `db:"tokenCharacterName" json:"tokenCharacterName,omitempty"`
	EntityID           int32  `db:"entityID" json:"entityID,omitempty"`
	EntityName         string `db:"entityName" json:"entityName,omitempty"`
	EntityType         string `db:"entityType" json:"entityType,omitempty"`
	Type               string `db:"type" json:"type,omitempty"`
	Types              string `db:"types" json:"types,omitempty"`
	Packed             string `db:"packed" json:"packed,omitempty"`
	Ignored            int32  `db:"ignored" json:"ignored,omitempty"`
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

		// Remove anything we didn't find
		if !touched[v.BotServiceID] {
			s.services.Delete(k)
		} else {
			// Update the server name while we are here
			serverName, err := v.Server.GetName()
			if err != nil {
				log.Println(err)
			} else if serverName != v.Name {
				err = s.updateServerName(v.BotServiceID, serverName)
				log.Println(err)
			}

			channels, err := v.Server.GetChannels()
			if err != nil {
				log.Println(err)
			} else {
				for _, ch := range channels {

					// Get the channel
					ci, ok := s.channels.Load(ch.ID)
					if ok {
						c := ci.(Channel)
						err = s.updateChannelName(v.BotServiceID, c.ChannelID, ch.Name)
						if err != nil {
							log.Println(err)
						}
					}
				}
			}

			roles, err := v.Server.GetRoles()
			if err != nil {
				log.Println(err)
			} else {

				err = s.updateRoles(v.BotServiceID, roles)
				if err != nil {
					log.Println(err)
				}

			}

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
			s.channels.Delete(k)
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
	err := s.db.Select(&services, "SELECT botServiceID, name, entityID, address, authentication, type, services FROM evedata.botServices")
	if err != nil {
		return nil, err
	}
	return services, nil
}

func (s *Conservator) updateServerName(b int32, name string) error {
	_, err := s.db.Exec("UPDATE evedata.botServices SET name = ? WHERE botServiceID = ?", name, b)
	if err != nil {
		return err
	}
	return nil
}

func (s *Conservator) updateChannelName(b int32, c, name string) error {
	_, err := s.db.Exec("UPDATE evedata.botChannels SET channelName = ? WHERE botServiceID = ? AND channelID = ?", name, b, c)
	if err != nil {
		return err
	}
	return nil
}

func (s *Conservator) getChannels() ([]Channel, error) {
	channels := []Channel{}
	err := s.db.Select(&channels, "SELECT botServiceID, channelID, services, options FROM evedata.botChannels")
	if err != nil {
		return nil, err
	}
	return channels, nil
}

func (s *Conservator) updateRoles(b int32, roles []botservice.Name) error {
	for _, r := range roles {
		_, err := s.db.Exec(`
		INSERT INTO evedata.botRoles 
			(botServiceID, roleID, roleName) VALUES(?,?,?) 
			ON DUPLICATE KEY UPDATE roleName = VALUES(roleName)`, b, r.ID, r.Name)
		if err != nil {
			return err
		}
	}
	return nil
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

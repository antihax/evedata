package conservator

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"log"

	"github.com/antihax/evedata/internal/botservice"
	"github.com/antihax/evedata/internal/botservice/discordservice"
)

type ServiceOptions struct {
	Auth struct {
		Members       string `json:"members,omitempty"`
		PlusFive      string `json:"plusFive,omitempty"`
		PlusTen       string `json:"plusTen,omitempty"`
		Militia       string `json:"militia,omitempty"`
		AlliedMilitia string `json:"alliedMilitia,omitempty"`
	} `json:"auth,omitempty"`
}

type ServiceTypes struct {
	Auth bool `json:"auth,omitempty"` // Authentication
}

func (c *ServiceTypes) GetServices() string {
	v := reflect.ValueOf(c).Elem()
	typeOf := v.Type()

	values := []string{}

	for i := 0; i < v.NumField(); i++ {
		b := v.Field(i).Interface().(bool)
		if b {
			values = append(values, strings.ToLower(typeOf.Field(i).Name))
		}
	}
	return strings.Join(values, ",")
}

type ChannelOptions struct {
	Killmail struct {
		IgnoreHighSec    bool `json:"ignoreHighsec,omitempty"`
		IgnoreLowSec     bool `json:"ignoreLowsec,omitempty"`
		IgnoreNullSec    bool `json:"ignoreNullsec,omitempty"`
		IgnoreWorthless  bool `json:"ignoreWorthless,omitempty"`
		War              bool `json:"war,omitempty"`
		FactionWar       bool `json:"factionWar,omitempty"`
		SendAll          bool `json:"sendAll,omitempty"`
		SendAllAbyssalT4 bool `json:"sendAllAbyssalT4,omitempty"`
	} `json:"killmail,omitempty"`
}

type ChannelTypes struct {
	War         bool `json:"war,omitempty"`         // war notifications
	Locator     bool `json:"locator,omitempty"`     // locator agent responses
	Kill        bool `json:"kill,omitempty"`        // killmails
	Structure   bool `json:"structure,omitempty"`   // structure notifications
	Application bool `json:"application,omitempty"` // applications to corp
}

func (c *ChannelTypes) GetServices() string {
	v := reflect.ValueOf(c).Elem()
	typeOf := v.Type()

	values := []string{}

	for i := 0; i < v.NumField(); i++ {
		b := v.Field(i).Interface().(bool)
		if b {
			values = append(values, strings.ToLower(typeOf.Field(i).Name))
		}
	}
	return strings.Join(values, ",")
}

type Service struct {
	Server         botservice.Integration `json:"-,omitempty"`
	IntegrationID  int32                  `db:"integrationID" json:"integrationID,omitempty"`
	FactionID      int32                  `db:"factionID" json:"factionID,omitempty"`
	Name           string                 `db:"name" json:"name,omitempty"`
	EntityID       int32                  `db:"entityID" json:"entityID,omitempty"`
	EntityName     string                 `db:"entityName" json:"entityName,omitempty"`
	EntityType     string                 `db:"entityType" json:"entityType,omitempty"`
	Address        string                 `db:"address" json:"address,omitempty" `
	Authentication string                 `db:"authentication,omitempty"`
	Type           string                 `db:"type" json:"type,omitempty"`
	Services       string                 `db:"services" json:"services,omitempty"`
	OptionsJSON    string                 `db:"options" json:"-"`
	Options        ServiceOptions         `db:"-" json:"options,omitempty"`
}

func (s *Service) checkRemoveRoles(memberID, roleToRemove string, memberRoles []string) error {
	if inSlice(roleToRemove, memberRoles) {
		return s.Server.RemoveRole(memberID, roleToRemove)
	}
	return nil
}

func (s *Service) checkAddRoles(memberID, roleToAdd string, memberRoles []string) error {
	if !inSlice(roleToAdd, memberRoles) {
		return s.Server.AddRole(memberID, roleToAdd)
	}
	return nil
}

type Channel struct {
	IntegrationID int32          `db:"integrationID" json:"integrationID,omitempty"`
	ChannelID     string         `db:"channelID"  json:"channelID,omitempty"`
	ChannelName   string         `db:"channelName"  json:"channelName,omitempty"`
	Services      string         `db:"services" json:"services,omitempty"`
	OptionsJSON   string         `db:"options" json:"-"`
	Options       ChannelOptions `db:"-" json:"options,omitempty"`
}

type Share struct {
	IntegrationID      int32  `db:"integrationID" json:"integrationID,omitempty"`
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
		var n botservice.Integration

		// Mark what we touch
		touched[service.IntegrationID] = true

		switch service.Type {
		case "discord":
			n = discordservice.NewDiscordService(s.discord, service.Address)
		default:
			return errors.New("unknown service type")
		}

		// Explode our options into the struct
		json.Unmarshal([]byte(service.OptionsJSON), &service.Options)

		// Store the server with the service information and put into our map
		service.Server = n
		s.services.Store(service.IntegrationID, service)
	}

	// Delete any items that were removed
	s.services.Range(func(ki, vi interface{}) bool {
		k := ki.(int32)
		v := vi.(Service)

		// Remove anything we didn't find
		if !touched[v.IntegrationID] {
			s.services.Delete(k)
		} else {
			// Update the server name while we are here
			serverName, err := v.Server.GetName()
			if err != nil {
				log.Println(err)
			} else if serverName != v.Name {
				err = s.updateServerName(v.IntegrationID, serverName)
				log.Println(err)
			}

			channels, err := v.Server.GetChannels()
			if err != nil {
				log.Println(err)
			} else {
				for _, ch := range channels {
					c, err := s.getChannel(ch.ID)
					if err != nil {
						continue
					}
					err = s.updateChannelName(v.IntegrationID, c.ChannelID, ch.Name)
					if err != nil {
						log.Println(err)
						continue
					}
				}
			}

			roles, err := v.Server.GetRoles()
			if err != nil {
				log.Println(err)
			} else {

				err = s.updateRoles(v.IntegrationID, roles)
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
				if a[i].TokenCharacterID == share.TokenCharacterID && a[i].IntegrationID == share.IntegrationID {
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
	err := s.db.Select(&services, "SELECT integrationID, name, entityID, address, authentication, type, services, options, factionID FROM evedata.integrations")
	if err != nil {
		return nil, err
	}
	return services, nil
}

func (s *Conservator) updateServerName(b int32, name string) error {
	_, err := s.db.Exec("UPDATE evedata.integrations SET name = ? WHERE integrationID = ?", name, b)
	if err != nil {
		return err
	}
	return nil
}

func (s *Conservator) updateChannelName(b int32, c, name string) error {
	_, err := s.db.Exec("UPDATE evedata.integrationChannels SET channelName = ? WHERE integrationID = ? AND channelID = ?", name, b, c)
	if err != nil {
		return err
	}
	return nil
}

func (s *Conservator) getChannels() ([]Channel, error) {
	channels := []Channel{}
	err := s.db.Select(&channels, "SELECT integrationID, channelID, services, options FROM evedata.integrationChannels")
	if err != nil {
		return nil, err
	}
	return channels, nil
}

func (s *Conservator) updateRoles(b int32, roles []botservice.Name) error {
	for _, r := range roles {
		_, err := s.db.Exec(`
		INSERT INTO evedata.integrationRoles 
			(integrationID, roleID, roleName) VALUES(?,?,?) 
			ON DUPLICATE KEY UPDATE roleName = VALUES(roleName)`, b, r.ID, r.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Conservator) updateData() {
	throttle := time.Tick(time.Second * 60)
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
		s.checkAllUsers()
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
		SELECT DISTINCT B.integrationID, tokenCharacterID, S.entityID, types, group_concat(concat(channelID, "|", C.services) SEPARATOR ":") as packed
		FROM evedata.sharing S
		INNER JOIN evedata.integrations B ON B.entityID = S.entityID
        INNER JOIN evedata.integrationChannels C ON C.integrationID = B.integrationID 
        GROUP BY integrationID, tokenCharacterID, entityID`)
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

package conservator

import (
	"errors"
	"fmt"
)

func (s *Conservator) getService(botServiceID int32) (*Service, error) {
	// Get the service
	si, ok := s.services.Load(botServiceID)
	if !ok {
		return nil, errors.New(fmt.Sprintf("missing botServiceID %d", botServiceID))
	}
	service, ok := si.(Service)
	if !ok {
		return nil, errors.New(fmt.Sprintf("missing botServiceID %d", botServiceID))
	}
	return &service, nil
}

func (s *Conservator) getServiceByAddress(address string) (*Service, error) {
	var service *Service
	s.services.Range(func(ki, vi interface{}) bool {
		v := vi.(Service)
		if v.Address == address {
			service = &v
			return false
		}
		return true
	})
	if service == nil {
		return service, errors.New("server not found")
	}
	return service, nil
}

func (s *Conservator) getChannel(channelID string) (*Channel, error) {
	// Get the service
	si, ok := s.channels.Load(channelID)
	if !ok {
		return nil, errors.New(fmt.Sprintf("missing channelID %s", channelID))
	}
	channel, ok := si.(Channel)
	if !ok {
		return nil, errors.New(fmt.Sprintf("missing channelID %s", channelID))
	}
	return &channel, nil
}

package conservator

import (
	"errors"
	"fmt"
)

func (s *Conservator) getService(integrationID int32) (*Service, error) {
	// Get the service
	si, ok := s.services.Load(integrationID)
	if !ok {
		return nil, fmt.Errorf("missing integrationID %d", integrationID)
	}
	service, ok := si.(Service)
	if !ok {
		return nil, fmt.Errorf("missing integrationID %d", integrationID)
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
		return nil, fmt.Errorf("missing channelID %s", channelID)
	}
	channel, ok := si.(Channel)
	if !ok {
		return nil, fmt.Errorf("missing channelID %s", channelID)
	}
	return &channel, nil
}

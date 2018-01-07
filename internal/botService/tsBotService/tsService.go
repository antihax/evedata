package tsBotService

import (
	"github.com/Darfk/ts3"
)

// AuthService provides access to a discord session
type TSService struct {
	session *ts3.Client
}

// NewTSService sends a message to a discord channel ID
func NewTSService(address, user, pass string) (*TSService, error) {
	conn, err := ts3.NewClient(address)
	if err != nil {
		return nil, err
	}

	_, err = conn.Exec(ts3.Login(user, pass))
	if err != nil {
		return nil, err
	}

	return &TSService{conn}, nil
}

// GetServerList gets the available TS3 virtual servers
func (c *TSService) UseServer(serverID int) error {
	_, err := c.session.Exec(ts3.Use(serverID))
	if err != nil {
		return err
	}

	return nil
}

// GetServerList gets the available TS3 virtual servers
func (c *TSService) GetServerList() ([]map[string]string, error) {
	res, err := c.session.Exec(ts3.Command{
		Command: "serverlist",
	})
	if err != nil {
		return nil, err
	}

	return res.Params, nil
}

// SendMessageToChannel sends a message to a discord channel ID
func (c *TSService) SendMessageToChannel(channel, message string) error {
	return nil
}

// SendMessageToUser sends a message to a discord user ID
func (c *TSService) SendMessageToUser(user, message string) error {
	return nil
}

// KickUser kicks a discord user ID
func (c *TSService) KickUser(user, message string) error {
	return nil
}

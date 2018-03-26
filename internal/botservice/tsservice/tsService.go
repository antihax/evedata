package tsservice

import (
	"strconv"

	"github.com/antihax/evedata/internal/botservice"

	ts3 "github.com/multiplay/go-ts3"
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

	err = conn.Login(user, pass)
	if err != nil {
		return nil, err
	}

	return &TSService{conn}, nil
}

// GetServerList gets the available TS3 virtual servers
func (c *TSService) UseServer(serverID int) error {
	err := c.session.Use(serverID)
	if err != nil {
		return err
	}

	return nil
}

// GetServerList gets the available TS3 virtual servers
func (c *TSService) GetServerList() (map[string]string, error) {
	list, err := c.session.Server.List()
	if err != nil {
		return nil, err
	}

	servers := make(map[string]string)
	for _, server := range list {
		servers[strconv.Itoa(server.ID)] = server.Name
	}

	return servers, nil
}

// GetChannelList gets the available TS3 channels
func (c *TSService) GetChannelList() (map[string]string, error) {
	list, err := c.session.Server.ChannelList()
	if err != nil {
		return nil, err
	}

	channels := make(map[string]string)
	for _, channel := range list {
		channels[strconv.Itoa(channel.ID)] = channel.ChannelName
	}

	return channels, nil
}

// SendMessageToChannel sends a message to a channel ID
func (c *TSService) SendMessageToChannel(channel, message string) error {
	_, err := c.session.ExecCmd(ts3.NewCmd("sendtextmessage").WithArgs(
		ts3.NewArg("targetmode", "2"),
		ts3.NewArg("target", channel),
		ts3.NewArg("msg", message),
	))
	return err
}

// SendMessageToUser sends a message to a discord user ID
func (c *TSService) SendMessageToUser(user, message string) error {
	_, err := c.session.ExecCmd(ts3.NewCmd("sendtextmessage").WithArgs(
		ts3.NewArg("targetmode", "1"),
		ts3.NewArg("target", user),
		ts3.NewArg("msg", message),
	))
	return err
}

// SendMessageToUser sends a message to a discord user ID
func (c *TSService) SendMessageToServer(user, message string) error {
	_, err := c.session.ExecCmd(ts3.NewCmd("sendtextmessage").WithArgs(
		ts3.NewArg("targetmode", "3"),
		ts3.NewArg("target", user),
		ts3.NewArg("msg", message),
	))
	return err
}

// KickUser kicks a discord user ID
func (c *TSService) KickUser(user, message string) error {
	return nil
}

// Get the server name
func (c *TSService) GetName() (string, error) {
	server := ts3.Server{}
	_, err := c.session.ExecCmd(ts3.NewCmd("serverinfo").WithResponse(&server))
	return server.Name, err
}

// IMPLIMENT
func (c *TSService) GetChannels() ([]botservice.Name, error) {

	return nil, nil
}

// IMPLIMENT
func (c *TSService) GetRoles() ([]botservice.Name, error) {

	return nil, nil
}

// IMPLIMENT
func (c *TSService) GetMembers() ([]botservice.Name, error) {

	return nil, nil
}

// IMPLIMENT
func (c *TSService) RemoveRole(user, role string) error {

	return nil
}

// IMPLIMENT
func (c *TSService) AddRole(user, role string) error {

	return nil
}

// AddUser
func (c *TSService) AddUser(auth, user, name string) error {
	return nil
}

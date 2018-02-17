package discordservice

import "github.com/bwmarrin/discordgo"

// BotService provides access to a discord session
// Discordgo handles rate throttling
type DiscordService struct {
	session  *discordgo.Session
	serverID string
}

// NewDiscordService sends a message to a discord channel ID
func NewDiscordService(session *discordgo.Session, serverID string) DiscordService {
	return DiscordService{session, serverID}
}

// SendMessageToChannel sends a message to a discord channel ID
func (c DiscordService) SendMessageToChannel(channel, message string) error {
	_, err := c.session.ChannelMessageSend(channel, message)
	return err
}

// SendMessageToUser sends a message to a discord user ID
func (c DiscordService) SendMessageToUser(user, message string) error {
	_, err := c.session.ChannelMessageSend(user, message)
	return err
}

// KickUser kicks a discord user ID
func (c DiscordService) KickUser(user, message string) error {
	_, err := c.session.ChannelMessageSend(user, message)
	return err
}

// Get the server name
func (c DiscordService) GetName() (string, error) {
	g, err := c.session.Guild(c.serverID)
	if err != nil {
		return "", err
	}
	return g.Name, nil
}

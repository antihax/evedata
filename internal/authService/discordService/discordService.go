package discordService

import "github.com/bwmarrin/discordgo"

// AuthService provides access to a discord session
type AuthService struct {
	session *discordgo.Session
}

// NewDiscordService sends a message to a discord channel ID
func NewDiscordService(token string) (*AuthService, error) {
	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	return &AuthService{bot}, nil
}

// SendMessageToChannel sends a message to a discord channel ID
func (c *AuthService) SendMessageToChannel(channel, message string) error {
	_, err := c.session.ChannelMessageSend(channel, message)
	return err
}

// SendMessageToUser sends a message to a discord user ID
func (c *AuthService) SendMessageToUser(user, message string) error {
	_, err := c.session.ChannelMessageSend(user, message)
	return err
}

// KickUser kicks a discord user ID
func (c *AuthService) KickUser(user, message string) error {
	_, err := c.session.ChannelMessageSend(user, message)
	return err
}

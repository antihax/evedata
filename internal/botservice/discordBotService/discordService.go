package discordbotservice

import "github.com/bwmarrin/discordgo"

// BotService provides access to a discord session
// Discordgo handles rate throttling
type DiscordService struct {
	session *discordgo.Session
}

// NewDiscordService sends a message to a discord channel ID
func NewDiscordService(token string) (*DiscordService, error) {
	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	return &DiscordService{bot}, nil
}

// SendMessageToChannel sends a message to a discord channel ID
func (c *DiscordService) SendMessageToChannel(channel, message string) error {
	_, err := c.session.ChannelMessageSend(channel, message)
	return err
}

// SendMessageToUser sends a message to a discord user ID
func (c *DiscordService) SendMessageToUser(user, message string) error {
	_, err := c.session.ChannelMessageSend(user, message)
	return err
}

// KickUser kicks a discord user ID
func (c *DiscordService) KickUser(user, message string) error {
	_, err := c.session.ChannelMessageSend(user, message)
	return err
}

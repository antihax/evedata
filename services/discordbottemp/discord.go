package discordbottemp

import (
	"errors"
	"log"

	"github.com/bwmarrin/discordgo"
)

// This is all one massive hack until we get microservices going for this.
// currently locked into stuff i need personally.

var dg *discordgo.Session

// GoDiscordBot runs a temporary hack of a bot as I line out functionality.
func (s *DiscordBot) connectToDiscord() error {
	var err error
	log.Printf("DiscordBot: Starting \n")
	dg, err = discordgo.New("Bot " + s.discordToken)
	if err != nil {
		return err
	}

	err = dg.Open()
	if err != nil {
		return err
	}
	return nil
}

// [TODO] Temporary Hack... test feasibility
func sendKillMessage(message string) error {
	if dg == nil {
		return errors.New("Not Connected")
	}
	_, err := dg.ChannelMessageSend("369208842443292675", message)
	return err
}

// [TODO] Temporary Hack... test feasibility
func sendNotificationMessage(message string) error {
	if dg == nil {
		return errors.New("Not Connected")
	}
	_, err := dg.ChannelMessageSend("369620236019695616", message)
	return err
}

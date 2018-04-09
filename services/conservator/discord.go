package conservator

import (
	"github.com/bwmarrin/discordgo"
)

var dg *discordgo.Session

func (c *Conservator) setupDiscord() (*discordgo.Session, error) {
	d, err := discordgo.New("Bot " + c.discordToken)
	if err != nil {
		return nil, err
	}
	d.AddHandler(c.discordMessageCreate)
	d.AddHandler(c.discordGuildCreate)
	return d, nil
}

func (c *Conservator) discordNewMember(s *discordgo.Session, e *discordgo.GuildMemberAdd) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if e.Member.User.ID == s.State.User.ID {
		return
	}
	c.handleNewMember(e.Member.User.ID, e.Member.Nick, e.GuildID)
}

func (c *Conservator) discordMessageCreate(s *discordgo.Session, e *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if e.Author.ID == s.State.User.ID {
		return
	}

}

func (c *Conservator) discordGuildCreate(s *discordgo.Session, e *discordgo.GuildCreate) {
	if e.Guild.Unavailable {
		return
	}

	for _, channel := range e.Guild.Channels {
		if channel.ID == e.Guild.ID {
			_, _ = s.ChannelMessageSend(channel.ID, "Discord Integration Complete - https://www.evedata.org")
			return
		}
	}
}

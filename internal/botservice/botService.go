package botservice

// AuthService provides access to an authenticated service
type BotService interface {
	SendMessageToChannel(message string) error
	SendMessageToUser(user, message string) error
	KickUser(user, message string) error
	GetChannelList() (map[string]string, error) // id: name
}

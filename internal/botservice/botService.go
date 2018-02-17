package botservice

// AuthService provides access to an authenticated service
type BotService interface {
	SendMessageToChannel(channel, message string) error
	SendMessageToUser(user, message string) error
	KickUser(user, message string) error
	GetName() (string, error)
}

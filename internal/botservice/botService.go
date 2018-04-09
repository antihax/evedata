package botservice

// Name returns internal ID, name, and roles from endpoints
type Name struct {
	ID    string
	Name  string
	Roles []string
}

// Integration provides access to an authenticated service
type Integration interface {
	SendMessageToChannel(channel, message string) error
	SendMessageToUser(user, message string) error
	KickUser(user, message string) error
	GetName() (string, error)
	GetChannels() ([]Name, error)
	GetRoles() ([]Name, error)
	GetMembers() ([]Name, error)
	RemoveRole(user, role string) error
	AddRole(user, role string) error
	AddUser(auth, user, name string) error
}

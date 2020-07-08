package chat

// Command formats
const (
	Connect    = "/connect>"
	Message    = "/message>"
	Disconnect = "/disconnect>"
	Users      = "/users>"
)

// Message Formats
const (
	UserListFormat = Users + "%s\n"
)

// View names
const (
	MessageView = "messages"
	InputView   = "input"
	UsersView   = "users"
	NameView    = "name"
)

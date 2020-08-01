package chat

//go:generate protoc --go_out=plugins=grpc:. chat.proto

// Message Type
const (
	UserList         = 1
	BroadcastMessage = 2
)

// View names
const (
	MessageView = "messages"
	InputView   = "input"
	UsersView   = "users"
	NameView    = "name"
)

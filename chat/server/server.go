package server

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/arunvm/chat_app/chat"
)

// Store keeps track of connected users
type Store struct {
	Mutex sync.Mutex
	Users map[string]connection
}

type connection struct {
	stream chat.Broadcast_ConnectServer
	err    chan error
}

// New returns an instance of store
func New() *Store {
	s := &Store{
		Users: make(map[string]connection),
	}
	return s
}

// Connect opens a new stream with the client and closes it when the function returns an error or nil
func (s *Store) Connect(user *chat.User, stream chat.Broadcast_ConnectServer) error {
	s.Mutex.Lock()
	s.Users[user.GetName()] = connection{
		stream: stream,
		err:    make(chan error),
	}
	s.Mutex.Unlock()

	s.broadcastOnlineUsers()

	// Steps to do before returning from connect function
	// - Closes the grpc stream by passing error/nil to the Connect function
	// - Removes the user from the map
	// - Retrieves the new list of users and broadcasts to all users that have the stream open
	err := <-s.Users[user.GetName()].err
	s.disconnect(user.GetName())
	s.broadcastOnlineUsers()

	return err
}

func (s *Store) BroadcastMessage(ctx context.Context, message *chat.Message) (*chat.Empty, error) {
	s.Mutex.Lock()
	for name, user := range s.Users {
		err := user.stream.Send(message)
		if err != nil {
			log.Printf("Error when sending message, Error=%v", err)
			s.Users[name].err <- err
		}
	}
	s.Mutex.Unlock()

	return &chat.Empty{}, nil
}

func (s *Store) Disconnect(ctx context.Context, user *chat.User) (*chat.Empty, error) {
	s.Users[user.GetName()].err <- nil
	return &chat.Empty{}, nil
}

func (s *Store) getUserList() []string {
	var users []string

	s.Mutex.Lock()
	for user := range s.Users {
		users = append(users, user)
	}
	s.Mutex.Unlock()

	return users
}

func (s *Store) broadcastOnlineUsers() {
	message := &chat.Message{
		Type:    chat.UserList,
		Message: strings.Join(s.getUserList(), " "),
	}

	s.BroadcastMessage(context.Background(), message)

	return
}

func (s *Store) disconnect(name string) {
	s.Mutex.Lock()
	delete(s.Users, name)
	s.Mutex.Unlock()
	return
}

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

func (s *Store) Connect(user *chat.User, stream chat.Broadcast_ConnectServer) error {
	s.Mutex.Lock()
	s.Users[user.GetName()] = connection{
		stream: stream,
		err:    make(chan error),
	}
	s.Mutex.Unlock()

	err := s.broadcastOnlineUsers()
	if err != nil {
		log.Printf("Error when broadcasting users,Error=%v", err)
		s.Users[user.GetName()].err <- err
	}

	return <-s.Users[user.GetName()].err
}

func (s *Store) BroadcastMessage(ctx context.Context, message *chat.Message) (*chat.Empty, error) {
	s.Mutex.Lock()
	for _, user := range s.Users {
		err := user.stream.Send(message)
		if err != nil {
			log.Fatalf("Error when sending message, Error=%v", err)
		}
	}
	s.Mutex.Unlock()

	return &chat.Empty{}, nil
}

func (s *Store) Disconnect(ctx context.Context, user *chat.User) (*chat.Empty, error) {
	s.Mutex.Lock()
	delete(s.Users, user.GetName())
	s.Mutex.Unlock()

	s.broadcastOnlineUsers()
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

func (s *Store) broadcastOnlineUsers() error {
	message := &chat.Message{
		Type:    chat.UserList,
		Message: strings.Join(s.getUserList(), " "),
	}

	_, err := s.BroadcastMessage(context.Background(), message)
	if err != nil {
		log.Printf("Error broadcasting online users,Error=%v", err)
		return err
	}

	return nil
}

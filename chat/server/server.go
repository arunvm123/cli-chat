package server

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

// Store keeps track of connected users
type Store struct {
	Mutex sync.Mutex
	Users map[string]*net.Conn
}

// New returns an instance of store
func New() *Store {
	s := &Store{
		Users: make(map[string]*net.Conn),
	}
	return s
}

// Connect adds a new user and their connection to the store
func (s *Store) connect(name string, conn *net.Conn) {
	s.Mutex.Lock()
	s.Users[name] = conn
	s.Mutex.Unlock()
	log.Println(s.Users)
}

// Handle handles all incoming messages
func (s *Store) Handle(conn *net.Conn) {
	reader := bufio.NewReader(*conn)

	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Error reading from connection, Error %v", err)
		}
		switch {
		case strings.HasPrefix(data, "/connect>"):
			name := strings.TrimSuffix(strings.Trim(data, "/connect>"), "\n")
			s.connect(name, conn)
			s.broadcast("/users>" + strings.Join(s.getUserList(), " ") + "\n")
		case strings.HasPrefix(data, "/message>"):
			log.Println(data)
			message := strings.Trim(data, "/message>")
			s.broadcast(message)
		case strings.HasPrefix(data, "/disconnect>"):
			name := strings.TrimSuffix(strings.Trim(data, "/disconnect>"), "\n")
			s.disconnect(name)
			s.broadcast("/users>" + strings.Join(s.getUserList(), " ") + "\n")
			break
		}
	}
}

// Broadcast sends the message to all the clients
func (s *Store) broadcast(message string) error {
	log.Println(message)
	s.Mutex.Lock()
	for name, c := range s.Users {
		writer := bufio.NewWriter(*c)
		_, err := writer.Write([]byte(message))
		if err != nil {
			log.Printf("Error writing to user %v, Error: %v", name, err)
			delete(s.Users, name)
		}

		writer.Flush()
	}

	s.Mutex.Unlock()
	return nil
}

// Disconnect removes the user and connection from the store
func (s *Store) disconnect(name string) {
	s.Mutex.Lock()
	delete(s.Users, name)
	s.Mutex.Unlock()
	log.Println(s.Users)
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

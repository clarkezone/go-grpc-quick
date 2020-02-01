package grpc_quick

import (
	"fmt"
	"log"

	grpc "google.golang.org/grpc"
)

// Client object
type Client struct {
	config     *ClientConf
	Connection *grpc.ClientConn
}

// CreateClient is a factory for servers
func CreateClient() *Client {
	client := &Client{}
	client.config = &ClientConf{}
	client.config.getClientConfEnvironment()
	if client.config.ServerPort == 0 {
		fmt.Printf("Config not detected in environment, attempting YAML\n")
		client.config.getClientConf()
	}

	return client
}

// Connect start serving
func (s *Client) Connect() {
	var err error

	if s.config.UseTLS == false {
		s.Connection, err = createclient(s.config.TLSServerName, s.config.ServerPort)
	} else {
		s.Connection, err = createclientsecure(s.config.TLSServerName, s.config.ServerPort, s.config.KeyWord)
	}

	if err != nil {
		log.Fatalf("unable to create connection: %v\n", err)
	}
}

// Disconnect from remote server
func (s *Client) Disconnect() {
	err := s.Connection.Close()
	if err != nil {
		log.Fatal("Couldn't close connection")
	}
}

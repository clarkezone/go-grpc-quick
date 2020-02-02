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

// GetClientConfig attempts to retrieve configuration from the environment, then a YAML file
func GetClientConfig() *ClientConf {
	config := getClientConfEnvironment()
	if config == nil {
		fmt.Printf("Config not detected in environment, attempting YAML\n")
		config = getClientConf()
	}
	return config
}

// CreateEmptyClientConfig creates an empty config
func CreateEmptyClientConfig() {
	createEmptyClientConfig()
}

// CreateClient is a factory for servers
func CreateClient(c *ClientConf) *Client {
	if c == nil {
		panic("Invalid config")
	}
	client := &Client{}
	client.config = c

	return client
}

// Connect start serving
func (s *Client) Connect() {
	var err error

	if s.config.UseTLS == false {
		s.Connection, err = createclient(s.config)
	} else {
		s.Connection, err = createclientsecure(s.config)
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

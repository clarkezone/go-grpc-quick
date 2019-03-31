package grpc_quick

import (
	"io/ioutil"
	"log"

	grpc "google.golang.org/grpc"
	yaml "gopkg.in/yaml.v2"
)

// Client

type clientconf struct {
	ServerPort    int    `yaml:"serverport"`
	TLSServerName string `yaml:"tlsservername"`
	IsSecure      bool   `yaml:"issecure"`
	KeyWord       string `yaml:"keyword"`
}

func (c *clientconf) getClientConf() {
	yamlFile, err := ioutil.ReadFile("clientconfig.yaml")
	//TODO: create an empty one
	if err != nil {
		log.Fatalf("Please create a clientconfig.yaml file  #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Invalid clientconfig.yaml: %v", err)
	}
}

// Client object
type Client struct {
	config     *clientconf
	Connection *grpc.ClientConn
}

// CreateClient is a factory for servers
func CreateClient() *Client {
	client := &Client{}
	client.config = &clientconf{}
	client.config.getClientConf()
	return client
}

// Connect start serving
func (s *Client) Connect() {
	var err error

	if s.config.IsSecure == false {
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

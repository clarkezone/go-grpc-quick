package grpc_quick

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	grpc "google.golang.org/grpc"
	yaml "gopkg.in/yaml.v2"
)

// Client

type clientconf struct {
	ServerPort    int    `yaml:"serverport"`
	TLSServerName string `yaml:"tlsservername"`
	UseTLS        bool   `yaml:"useTLS"`
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

func (c *clientconf) getClientConfEnvironment() {
	i, err := strconv.ParseInt(os.Getenv("SERVERPORT"), 10, 32)
	if err == nil {
		c.ServerPort = int(i)
	}

	c.TLSServerName = os.Getenv("TLSSERVERNAME")

	b, err := strconv.ParseBool(os.Getenv("USETLS"))
	if err == nil {
		c.UseTLS = b
	}

	c.KeyWord = os.Getenv("KEYWORD")
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

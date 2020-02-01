package grpc_quick

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

type ClientConf struct {
	ServerPort    int    `yaml:"serverport"`
	TLSServerName string `yaml:"tlsservername"`
	UseTLS        bool   `yaml:"useTLS"`
	KeyWord       string `yaml:"keyword"`
}

func (c *ClientConf) getClientConf() {
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

func (c *ClientConf) getClientConfEnvironment() {
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

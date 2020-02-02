package grpc_quick

import (
	"io/ioutil"
	"os"
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

//ClientConf is a client configuraiton struct
type ClientConf struct {
	ServerPort      int    `yaml:"serverport"`
	TLSServerName   string `yaml:"tlsservername"`
	UseTLS          bool   `yaml:"useTLS"`
	PerCallSecurity bool   `yaml:"percallsecurity"`
	KeyWord         string `yaml:"keyword"`
}

func createEmptyClientConfig() bool {
	cnf := ClientConf{}
	data, err := yaml.Marshal(cnf)
	if err != nil {
		return false
	}
	ioutil.WriteFile("clientconfig.yaml", data, 0777)
	return true
}

func getClientConf() *ClientConf {
	c := &ClientConf{}
	yamlFile, err := ioutil.ReadFile("clientconfig.yaml")
	//TODO: create an empty one
	if err != nil {
		return nil
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil
	}
	return c
}

func getClientConfEnvironment() *ClientConf {
	c := &ClientConf{}
	i, err := strconv.ParseInt(os.Getenv("SERVERPORT"), 10, 32)
	if err == nil {
		c.ServerPort = int(i)
	} else {
		return nil
	}

	c.TLSServerName = os.Getenv("TLSSERVERNAME")

	b, err := strconv.ParseBool(os.Getenv("USETLS"))
	if err == nil {
		c.UseTLS = b
	} else {
		return nil
	}

	b, err = strconv.ParseBool(os.Getenv("PERCALLSECURITY"))
	if err == nil {
		c.PerCallSecurity = b
	} else {
		return nil
	}

	c.KeyWord = os.Getenv("KEYWORD")
	return c
}

package grpc_quick

import (
	"io/ioutil"
	"os"
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

//ServerConf is the server configuration struct
type ServerConf struct {
	ServerPort      int    `yaml:"serverport"`
	ServerCertPort  int    `yaml:"servercertport"`
	TLSServerName   string `yaml:"tlsservername"`
	UseTLS          bool   `yaml:"usetls"`
	PerCallSecurity bool   `yaml:"percallsecurity"`

	KeyWord string `yaml:"keyword"`
}

func createEmptyConfig() bool {
	cnf := ServerConf{}
	data, err := yaml.Marshal(cnf)
	if err != nil {
		return false
	}
	ioutil.WriteFile("serverconfig.yaml", data, 0777)
	return true
}

func getServerConf() *ServerConf {
	c := &ServerConf{}
	yamlFile, err := ioutil.ReadFile("serverconfig.yaml")
	if err != nil {
		return nil
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil
	}
	//TODO validate config
	return c
}

func getServerConfEnvironment() *ServerConf {
	c := &ServerConf{}
	i, err := strconv.ParseInt(os.Getenv("SERVERPORT"), 10, 32)
	if err == nil {
		c.ServerPort = int(i)
	} else {
		return nil
	}

	i, err = strconv.ParseInt(os.Getenv("SERVERCERTPORT"), 10, 32)
	if err == nil {
		c.ServerCertPort = int(i)
	} else {
		return nil
	}

	c.TLSServerName = os.Getenv("TLSSERVERNAME")
	if c.TLSServerName == "" {
		return nil
	}

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

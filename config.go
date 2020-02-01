package grpc_quick

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

func (c *conf) getServerConf() bool {
	yamlFile, err := ioutil.ReadFile("serverconfig.yaml")
	if err != nil {
		fmt.Printf("Please create a serverconfig.yaml file  #%v ", err)
		return false
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Printf("Invalid serverconfig.yaml: %v", err)
		return false
	}
	//TODO validate config
	return true
}

func (c *conf) getServerConfEnvironment() bool {
	i, err := strconv.ParseInt(os.Getenv("SERVERPORT"), 10, 32)
	if err == nil {
		c.ServerPort = int(i)
	} else {
		return false
	}

	i, err = strconv.ParseInt(os.Getenv("SERVERCERTPORT"), 10, 32)
	if err == nil {
		c.ServerCertPort = int(i)
	} else {
		return false
	}

	c.TLSServerName = os.Getenv("TLSSERVERNAME")
	if c.TLSServerName == "" {
		return false
	}

	b, err := strconv.ParseBool(os.Getenv("USETLS"))
	if err == nil {
		c.UseTLS = b
	} else {
		return false
	}

	b, err = strconv.ParseBool(os.Getenv("PERCALLSECURITY"))
	if err == nil {
		c.PerCallSecurity = b
	} else {
		return false
	}

	c.KeyWord = os.Getenv("KEYWORD")

	return true
}

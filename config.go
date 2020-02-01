package grpc_quick

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

type conf struct {
	ServerPort     int    `yaml:"serverport"`
	ServerCertPort int    `yaml:"servercertport"`
	TLSServerName  string `yaml:"tlsservername"`
	IsSecure       bool   `yaml:"issecure"`
	KeyWord        string `yaml:"keyword"`
}

func createEmptyConfig() {
	cnf := conf{}
	data, err := yaml.Marshal(cnf)
	if err != nil {
		log.Fatalf("File write failed  #%v ", err)
	}
	ioutil.WriteFile("serverconfig.yaml", data, 0777)
}

func getServerConf() (c *conf) {
	yamlFile, err := ioutil.ReadFile("serverconfig.yaml")
	//TODO: create an empty one
	if err != nil {
		//TODO validate config

		//if err != nil {
		//	log.Fatalf("Please create a serverconfig.yaml file  #%v ", err)
		//}
		return nil
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Invalid serverconfig.yaml: %v", err)
	}
	return c
}

func  getServerConfEnvironment() (*conf){
	c := &conf{}
	i, err := strconv.ParseInt(os.Getenv("SERVERPORT"), 10, 32)
	if err == nil {
		c.ServerPort = int(i)
	}

	i, err = strconv.ParseInt(os.Getenv("SERVERCERTPORT"), 10, 32)
	if err == nil {
		c.ServerCertPort = int(i)
	}

	c.TLSServerName = os.Getenv("TLSSERVERNAME")

	b, err := strconv.ParseBool(os.Getenv("ISSECURE"))
	if err == nil {
		c.IsSecure = b
	}

	c.KeyWord = os.Getenv("KEYWORD")
	return c
}

func GetConfiguration() (*conf){
	config := getServerConfEnvironment()
	if config.ServerPort == 0 {
		fmt.Printf("Config not detected in environment, attempting YAML\n")
		config = getServerConf()
		if config ==nil {
			createEmptyConfig()
		}
		return config
	}
	return config
}
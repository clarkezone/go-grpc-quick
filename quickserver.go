package grpc_quick

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"

	grpc "google.golang.org/grpc"
	yaml "gopkg.in/yaml.v2"
)

type conf struct {
	ServerPort     int    `yaml:"serverport"`
	ServerCertPort int    `yaml:"servercertport"`
	TLSServerName  string `yaml:"tlsservername"`
	IsSecure       bool   `yaml:"issecure"`
	KeyWord        string `yaml:"keyword"`
}

type registerCallback func(*grpc.Server)

func (c *conf) getServerConf() {
	yamlFile, err := ioutil.ReadFile("serverconfig.yaml")
	//TODO: create an empty one
	if err != nil {

		c = &conf{}
		i, err := strconv.ParseInt(os.Getenv("SERVERPORT"), 0, 10000)
		if err == nil {
			c.ServerPort = int(i)
		}

		i, err = strconv.ParseInt(os.Getenv("SERVERCERTPORT"), 0, 10000)
		if err == nil {
			c.ServerCertPort = int(i)
		}

		c.TLSServerName = os.Getenv("TLSSERVERNAME")

		b, err := strconv.ParseBool(os.Getenv("ISECURE"))
		if err == nil {
			c.IsSecure = b
		}

		c.KeyWord = os.Getenv("KEYWORD")

		//TODO validate config

		if err != nil {
			log.Fatalf("Please create a serverconfig.yaml file  #%v ", err)
		}
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Invalid serverconfig.yaml: %v", err)
	}
}

// Server object
type Server struct {
	config *conf
}

// CreateServer is a factory for servers
func CreateServer() *Server {
	serv := &Server{}
	serv.config = &conf{}
	serv.config.getServerConf()
	return serv
}

// Serve start serving
func (s *Server) Serve(regcb registerCallback) {
	if s.config.IsSecure {
		s.servegRPCAutoCert(s.config.TLSServerName, s.config.ServerPort, s.config.ServerCertPort, regcb)
	} else {
		servegRPC(s.config.TLSServerName, s.config.ServerPort, regcb)
	}
}

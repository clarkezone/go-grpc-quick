package grpc_quick

import (
	"io/ioutil"
	"log"

	grpc "google.golang.org/grpc"
	yaml "gopkg.in/yaml.v2"
)

type conf struct {
	ServerPort    int    `yaml:"serverport"`
	TlsServerName string `yaml:"tlsservername"`
}

type thing func(*grpc.Server)

func (c *conf) getConf() {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Please create a config.yaml file  #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Invalid config.yaml: %v", err)
	}
}

type Server struct {
	config *conf
}

func CreateServer() *Server {
	serv := &Server{}
	serv.config = &conf{}
	serv.config.getConf()
	return serv
}

func (s *Server) Serve(regcb thing) {
	servegRPC(s.config.TlsServerName, s.config.ServerPort, regcb)
}

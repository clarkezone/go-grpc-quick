package grpc_quick

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	grpc "google.golang.org/grpc"
	yaml "gopkg.in/yaml.v2"
)

type conf struct {
	ServerPort      int    `yaml:"serverport"`
	ServerCertPort  int    `yaml:"servercertport"`
	TLSServerName   string `yaml:"tlsservername"`
	UseTLS          bool   `yaml:"usetls"`
	PerCallSecurity bool   `yaml:"percallsecurity"`

	KeyWord string `yaml:"keyword"`
}

type registerCallback func(*grpc.Server)

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

// Server object
type Server struct {
	config *conf
}

// CreateServer is a factory for servers
func CreateServer() *Server {
	serv := &Server{}
	serv.config = &conf{}
	serv.config.getServerConfEnvironment()
	if serv.config.ServerPort == 0 {
		fmt.Printf("Config not detected in environment, attempting YAML\n")
		if !serv.config.getServerConf() {
			return nil
		}
	}
	return serv
}

// Serve start serving
func (s *Server) Serve(ctx context.Context, regcb registerCallback) {
	ctx, cancel := context.WithCancel(ctx)
	var srv *grpc.Server
	if s.config.UseTLS {
		srv = s.servegRPCAutoCert(s.config.TLSServerName, s.config.ServerPort, s.config.ServerCertPort, regcb, cancel, s.config.PerCallSecurity)
	} else {
		srv = servegRPC(s.config.TLSServerName, s.config.ServerPort, regcb, cancel)
	}
	<-ctx.Done()
	srv.GracefulStop()
	log.Println("Serve complete")
}

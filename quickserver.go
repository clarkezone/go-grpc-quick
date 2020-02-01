package grpc_quick

import (
	"context"
	"fmt"
	"log"

	grpc "google.golang.org/grpc"
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

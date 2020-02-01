package grpc_quick

import (
	"context"
	"fmt"
	"log"

	grpc "google.golang.org/grpc"
)

type registerCallback func(*grpc.Server)

// Server object
type Server struct {
	config *ServerConf
}

// GetConfig attempts to retrieve configuration from the environment, then a YAML file
func GetConfig() *ServerConf {
	config := getServerConfEnvironment()
	if config == nil {
		fmt.Printf("Config not detected in environment, attempting YAML\n")
		config = getServerConf()
	}
	return config
}

// CreateServer is a factory for servers
func CreateServer(c *ServerConf) *Server {
	if c == nil {
		panic("Invalid config")
	}
	serv := &Server{config: c}

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

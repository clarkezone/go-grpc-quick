package grpc_quick

import (
	"context"
	"log"

	grpc "google.golang.org/grpc"
)

type registerCallback func(*grpc.Server)



// Server object
type Server struct {
	config *conf
}

// CreateServer is a factory for servers
func CreateServer(config *conf) *Server {
	if config == nil {
		log.Fatal("Bad configuration")
	}
	serv := &Server{}
	serv.config = config
	
	return serv
}

// Serve start serving
func (s *Server) Serve(ctx context.Context, regcb registerCallback) {
	ctx, cancel := context.WithCancel(ctx)
	var srv *grpc.Server
	if s.config.IsSecure {
		srv = s.servegRPCAutoCert(s.config.TLSServerName, s.config.ServerPort, s.config.ServerCertPort, regcb, cancel)
	} else {
		srv = servegRPC(s.config.TLSServerName, s.config.ServerPort, regcb, cancel)
	}
	<-ctx.Done()
	srv.GracefulStop()
	log.Println("Serve complete")
}

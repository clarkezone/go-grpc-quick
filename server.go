package grpc_quick

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

func servegRPC(serverName string, serverPort int, cb registerCallback, cancel context.CancelFunc) *grpc.Server {
	fmt.Printf("Serving gRPC for endpoint %v on port %v\n", serverName, serverPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", serverPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	cb(grpcServer)
	go func() {
		grpcServer.Serve(lis)

		cancel()
	}()
	return grpcServer
}

func (be *Server) servegRPCAutoCert(serverName string, serverPort int, serverCertPort int, cb registerCallback, cancel context.CancelFunc) *grpc.Server {
	fmt.Printf("Serving gRPC AutoCert for endpoint %v on port %v\nwith certificate completion on port %v with keyword %v\n", serverName, serverPort, serverCertPort, be.config.KeyWord)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", serverPort))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer, err := be.listenWithAutoCert(serverName, serverCertPort)
	if err != nil {
		log.Fatalf("failed to listenwithautocert: %v", err)
	}

	cb(grpcServer)

	go func() {
		grpcServer.Serve(lis)

		cancel()
	}()

	return grpcServer
}

func (be *Server) listenWithAutoCert(serverName string, certport int) (*grpc.Server, error) {
	m := &autocert.Manager{
		Cache:      autocert.DirCache("//tlsdata"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(serverName),
	}

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%v", certport), m.HTTPHandler(nil))
		if err != nil {
			log.Fatalf("Error starting certificate completion thread: %v\n", err)
		}
	}()
	creds := credentials.NewTLS(&tls.Config{GetCertificate: m.GetCertificate})

	opts := []grpc.ServerOption{grpc.Creds(creds),
		grpc.UnaryInterceptor(be.unaryInterceptor),
		grpc.StreamInterceptor(be.streamInterceptor),
	}

	srv := grpc.NewServer(opts...)
	reflection.Register(srv)
	return srv, nil
}

func (be *Server) unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod == "/proto.EventStoreService/GetJWT" { //skip auth when requesting JWT

		return handler(ctx, req)
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		clientLogin := strings.Join(md["login"], "")

		if clientLogin != be.config.KeyWord {
			return nil, fmt.Errorf("bad creds")
		}

		return handler(ctx, req)
	}

	return nil, fmt.Errorf("missing credentials")
}

func (be *Server) streamInterceptor(req interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if md, ok := metadata.FromIncomingContext(ss.Context()); ok {
		clientLogin := strings.Join(md["login"], "")

		if clientLogin != be.config.KeyWord {
			return fmt.Errorf("bad creds")
		}

		return handler(req, ss)
	}
	return fmt.Errorf("missing credentials")
}

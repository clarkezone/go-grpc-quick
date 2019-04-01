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

func servegRPC(serverName string, serverPort int, cb registerCallback) {
	fmt.Printf("Serving gRPC for endpoint %v on port %v\n", serverName, serverPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", serverPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	cb(grpcServer)
	//helloServer := HelloServer{}
	//jamestestrpc.RegisterJamesTestServiceServer(grpcServer, &helloServer)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (be *Server) servegRPCAutoCert(serverName string, serverPort int, serverCertPort int, cb registerCallback) {
	fmt.Printf("Serving gRPC AutoCert for endpoint %v on port %v\n", serverName, serverPort)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", serverPort))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer, err := be.listenWithAutoCert(serverName, serverPort, serverCertPort)
	if err != nil {
		log.Fatalf("failed to listenwithautocert: %v", err)
	}

	cb(grpcServer)
	// helloServer := HelloServer{}
	// jamestestrpc.RegisterJamesTestServiceServer(grpcServer, &helloServer)
	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatalf("failed to serve grpc with autocert: %v", err)
	}
}

func (be *Server) listenWithAutoCert(serverName string, p int, certport int) (*grpc.Server, error) {
	m := &autocert.Manager{
		Cache:      autocert.DirCache("tls"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(serverName),
	}
	//todo: do we actually need to listen here to get autocert to work?  If yes, put port in config
	go http.ListenAndServe(fmt.Sprintf(":%v, certport"), m.HTTPHandler(nil))
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

		//ctx = context.WithValue(ctx, clientIDKey, clientID)
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

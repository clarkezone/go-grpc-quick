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

func (be *Server) servegRPC(serverName string, serverPort int, cb registerCallback, usePerCallSecurity bool, cancel context.CancelFunc) *grpc.Server {
	fmt.Printf("Serving gRPC for endpoint %v on port %v\n", serverName, serverPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", serverPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	if usePerCallSecurity {
		opts = []grpc.ServerOption{
			grpc.UnaryInterceptor(be.unaryInterceptor),
			grpc.StreamInterceptor(be.streamInterceptor),
		}
	}

	grpcServer := grpc.NewServer(opts...)
	cb(grpcServer)
	go func() {
		grpcServer.Serve(lis)

		cancel()
		log.Println("grpc Serve goroutine has exited, cancel called")
	}()
	return grpcServer
}

func (be *Server) servegRPCAutoCert(conf *ServerConf, cb registerCallback, cancel context.CancelFunc) *grpc.Server {
	fmt.Printf("Serving gRPC AutoCert for endpoint %v on port %v\nwith certificate completion on port %v with keyword %v\n", conf.TLSServerName, conf.ServerPort, conf.ServerCertPort, be.config.KeyWord)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", conf.ServerPort))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if conf.TLSCacheDir == "" {
		conf.TLSCacheDir = "."
	}

	grpcServer, err := be.listenWithAutoCert(conf.TLSServerName, conf.ServerCertPort, conf.PerCallSecurity, conf.TLSCacheDir)
	if err != nil {
		log.Fatalf("failed to listenwithautocert: %v", err)
	}

	cb(grpcServer)

	go func() {
		grpcServer.Serve(lis)

		cancel()
		log.Println("grpc Serve goroutine has exited, cancel called")
	}()

	return grpcServer
}

func (be *Server) listenWithAutoCert(serverName string, certport int, usePerCallSecurity bool, cacheDir string) (*grpc.Server, error) {
	m := &autocert.Manager{
		Cache:      autocert.DirCache(cacheDir),
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

	var opts []grpc.ServerOption

	if usePerCallSecurity {
		opts = []grpc.ServerOption{grpc.Creds(creds),
			grpc.UnaryInterceptor(be.unaryInterceptor),
			grpc.StreamInterceptor(be.streamInterceptor),
		}
	} else {
		opts = []grpc.ServerOption{grpc.Creds(creds)}
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

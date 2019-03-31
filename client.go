package grpc_quick

import (
	"context"
	"crypto/tls"
	"fmt"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func createclient(servername string, port int) (*grpc.ClientConn, error) {
	fmt.Println("Client")

	return grpc.Dial(fmt.Sprintf("%v:%v", servername, port), grpc.WithInsecure())
}

func createclientsecure(servername string, port int, keyword string) (*grpc.ClientConn, error) {
	fmt.Printf("Client Secure %v %v\n", servername, port)

	conf := &tls.Config{ServerName: servername}

	creds := credentials.NewTLS(conf)

	auth := authentication{Login: keyword}

	return grpc.Dial(fmt.Sprintf("%v:%v", servername, port), grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(&auth))
}

type authentication struct {
	Login string
}

func (a *authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"login": a.Login,
	}, nil
}

func (a *authentication) RequireTransportSecurity() bool {
	return true
}

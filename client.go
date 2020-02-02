package grpc_quick

import (
	"context"
	"crypto/tls"
	"fmt"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func createclient(confg *ClientConf) (*grpc.ClientConn, error) {
	fmt.Println("Client")

	if confg.PerCallSecurity {
		auth := authentication{Login: confg.KeyWord}
		return grpc.Dial(fmt.Sprintf("%v:%v", confg.TLSServerName, confg.ServerPort), grpc.WithInsecure(), grpc.WithPerRPCCredentials(&auth))
	}

	return grpc.Dial(fmt.Sprintf("%v:%v", confg.TLSServerName, confg.ServerPort), grpc.WithInsecure())
}

func createclientsecure(confg *ClientConf) (*grpc.ClientConn, error) {
	fmt.Printf("Client Secure %v %v with keyword %v\n", confg.TLSServerName, confg.ServerPort, confg.KeyWord)

	conf := &tls.Config{ServerName: confg.TLSServerName}

	creds := credentials.NewTLS(conf)

	auth := authentication{Login: confg.KeyWord}

	if confg.PerCallSecurity {
		return grpc.Dial(fmt.Sprintf("%v:%v", confg.TLSServerName, confg.ServerPort), grpc.WithTransportCredentials(creds),
			grpc.WithPerRPCCredentials(&auth))
	}

	return grpc.Dial(fmt.Sprintf("%v:%v", confg.TLSServerName, confg.ServerPort), grpc.WithTransportCredentials(creds))
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

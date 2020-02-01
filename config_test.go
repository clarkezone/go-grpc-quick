package grpc_quick

import (
	"os"
	"testing"
)

func TestSerSuccess(t *testing.T) {
	os.Setenv("SERVERPORT", "8443")
	os.Setenv("SERVERCERTPORT", "8080")
	config := getServerConfEnvironment()
	if config.ServerPort != 8443 {
		t.Fatalf("Port wasn't correct")
	}
}

func TestGetServerConf(t *testing.T) {
	
}

func TestCreateEmptyConf(t *testing.T) {

}
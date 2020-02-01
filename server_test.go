package grpc_quick

import (
	"os"
	"testing"
)

func TestSerSuccess(t *testing.T) {
	os.Setenv("SERVERPORT", "8443")
	os.Setenv("SERVERCERTPORT", "8080")
	os.Setenv("USETLS", "TRUE")
	os.Setenv("TLSSERVERNAME", "FOO")
	os.Setenv("PERCALLSECURITY", "FALSE")
	config := &conf{}
	if !config.getServerConfEnvironment() {
		t.Fatalf("Config wasn't complete")
	}
	if config.ServerPort != 8443 {
		t.Fatalf("Port wasn't correct")
	}
	if config.ServerCertPort != 8080 {
		t.Fatalf("Cert Port wasn't correct")
	}
	if config.UseTLS != true {
		t.Fatalf("UseTls not correct")
	}
	if config.TLSServerName != "FOO" {
		t.Fatalf("Server name wasn't correct")
	}
	if config.PerCallSecurity != false {
		t.Fatalf("per call security broken")
	}
	os.Setenv("SERVERPORT", "")
	os.Setenv("SERVERCERTPORT", "")
	os.Setenv("USETLS", "")
	os.Setenv("TLSSERVERNAME", "")
	os.Setenv("PERCALLSECURITY", "")
}

func TestDetectBadPort(t *testing.T) {
	config := &conf{}
	if config.getServerConfEnvironment() {
		t.Fatalf("Config not correctly rejectedf")
	}
}

func TestDetectBadCertPort(t *testing.T) {
	os.Setenv("SERVERPORT", "8443")
	config := &conf{}
	if config.getServerConfEnvironment() {
		t.Fatalf("Config not correctly rejectedf")
	}
	os.Setenv("SERVERPORT", "")
}

func TestDetectBadServerName(t *testing.T) {
	os.Setenv("SERVERPORT", "8443")
	os.Setenv("SERVERCERTPORT", "8443")
	config := &conf{}
	if config.getServerConfEnvironment() {
		t.Fatalf("Config not correctly rejectedf")
	}
	os.Setenv("SERVERPORT", "")
	os.Setenv("SERVERCERTPORT", "")
}

func TestDetectBadTLS(t *testing.T) {
	os.Setenv("SERVERPORT", "8443")
	os.Setenv("SERVERCERTPORT", "8443")
	os.Setenv("TLSSERVERNAME", "8443")
	config := &conf{}
	if config.getServerConfEnvironment() {
		t.Fatalf("Config not correctly rejectedf")
	}
	os.Setenv("SERVERPORT", "")
	os.Setenv("SERVERCERTPORT", "")
	os.Setenv("TLSSERVERNAME", "")
}

func TestDetectBadPerCallSecurity(t *testing.T) {
	os.Setenv("SERVERPORT", "8443")
	os.Setenv("SERVERCERTPORT", "8443")
	os.Setenv("TLSSERVERNAME", "8443")
	os.Setenv("USETLS", "TRUE")
	config := &conf{}
	if config.getServerConfEnvironment() {
		t.Fatalf("Config not correctly rejectedf")
	}
	os.Setenv("SERVERPORT", "")
	os.Setenv("SERVERCERTPORT", "")
	os.Setenv("TLSSERVERNAME", "")
	os.Setenv("USETLS", "")
}

func TestCreateServer(t *testing.T) {
	srv := CreateServer()

	if srv != nil {
		t.Fatalf("there should be no config because no YAML or environment")
	}
}

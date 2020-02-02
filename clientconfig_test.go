package grpc_quick

import (
	"os"
	"testing"
)

func TestClientSuccess(t *testing.T) {
	os.Setenv("SERVERPORT", "8443")
	os.Setenv("USETLS", "TRUE")
	os.Setenv("TLSSERVERNAME", "FOO")
	os.Setenv("PERCALLSECURITY", "TRUE")

	config := getClientConfEnvironment()
	if config == nil {
		t.Fatalf("Config wasn't complete")
	}
	if config.ServerPort != 8443 {
		t.Fatalf("Port wasn't correct")
	}
	if config.UseTLS != true {
		t.Fatalf("UseTls not correct")
	}
	if config.TLSServerName != "FOO" {
		t.Fatalf("Server name wasn't correct")
	}
	if config.PerCallSecurity != true {
		t.Fatalf("per call security broken")
	}
	os.Setenv("SERVERPORT", "")
	os.Setenv("USETLS", "")
	os.Setenv("TLSSERVERNAME", "")
	os.Setenv("PERCALLSECURITY", "")
}

func TestClientDetectBadPort(t *testing.T) {
	config := getClientConfEnvironment()
	if config != nil {
		t.Fatalf("Config not correctly rejectedf")
	}
}

func TestClientDetectBadCertPort(t *testing.T) {
	os.Setenv("SERVERPORT", "8443")
	config := getClientConfEnvironment()
	if config != nil {
		t.Fatalf("Config not correctly rejectedf")
	}
	os.Setenv("SERVERPORT", "")
}

func TestClientDetectBadServerName(t *testing.T) {
	os.Setenv("SERVERPORT", "8443")
	os.Setenv("SERVERCERTPORT", "8443")
	config := getClientConfEnvironment()
	if config != nil {
		t.Fatalf("Config not correctly rejectedf")
	}
	os.Setenv("SERVERPORT", "")
	os.Setenv("SERVERCERTPORT", "")
}

func TestClientDetectBadTLS(t *testing.T) {
	os.Setenv("SERVERPORT", "8443")
	os.Setenv("SERVERCERTPORT", "8443")
	os.Setenv("TLSSERVERNAME", "8443")
	config := getClientConfEnvironment()
	if config != nil {
		t.Fatalf("Config not correctly rejectedf")
	}
	os.Setenv("SERVERPORT", "")
	os.Setenv("SERVERCERTPORT", "")
	os.Setenv("TLSSERVERNAME", "")
}

func TestClientDetectBadPerCallSecurity(t *testing.T) {
	os.Setenv("SERVERPORT", "8443")
	os.Setenv("SERVERCERTPORT", "8443")
	os.Setenv("TLSSERVERNAME", "8443")
	os.Setenv("USETLS", "TRUE")
	config := getClientConfEnvironment()
	if config != nil {
		t.Fatalf("Config not correctly rejectedf")
	}
	os.Setenv("SERVERPORT", "")
	os.Setenv("SERVERCERTPORT", "")
	os.Setenv("TLSSERVERNAME", "")
	os.Setenv("USETLS", "")
}

package grpc_quick

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestSerSuccess(t *testing.T) {
	os.Setenv("SERVERPORT", "8443")
	os.Setenv("SERVERCERTPORT", "8080")
	os.Setenv("USETLS", "TRUE")
	os.Setenv("TLSSERVERNAME", "FOO")
	os.Setenv("PERCALLSECURITY", "FALSE")
	os.Setenv("TLSCACHEDIR", "BAR")

	config := getServerConfEnvironment()
	if config == nil {
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
	if config.TLSCacheDir != "BAR" {
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
	os.Setenv("TLSCACHEDIR", "")
}

func TestDetectBadPort(t *testing.T) {
	config := getServerConfEnvironment()
	if config != nil {
		t.Fatalf("Config not correctly rejectedf")
	}
}

func TestDetectBadCertPort(t *testing.T) {
	os.Setenv("SERVERPORT", "8443")
	config := getServerConfEnvironment()
	if config != nil {
		t.Fatalf("Config not correctly rejectedf")
	}
	os.Setenv("SERVERPORT", "")
}

func TestDetectBadServerName(t *testing.T) {
	os.Setenv("SERVERPORT", "8443")
	os.Setenv("SERVERCERTPORT", "8443")
	config := getServerConfEnvironment()
	if config != nil {
		t.Fatalf("Config not correctly rejectedf")
	}
	os.Setenv("SERVERPORT", "")
	os.Setenv("SERVERCERTPORT", "")
}

func TestDetectBadTLS(t *testing.T) {
	os.Setenv("SERVERPORT", "8443")
	os.Setenv("SERVERCERTPORT", "8443")
	os.Setenv("TLSSERVERNAME", "8443")
	config := getServerConfEnvironment()
	if config != nil {
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
	config := getServerConfEnvironment()
	if config != nil {
		t.Fatalf("Config not correctly rejectedf")
	}
	os.Setenv("SERVERPORT", "")
	os.Setenv("SERVERCERTPORT", "")
	os.Setenv("TLSSERVERNAME", "")
	os.Setenv("USETLS", "")
}

func TestGetConfig(t *testing.T) {
	yamlString := `
serverport: 8443
servercertport: 8080
tlsservername: iplayerdev.objectivepixel.io
usetls: true
percallsecurity: false
keyword: foobar
tlscachedir: "//tlsdata"
`
	bytes := []byte(yamlString)

	err := ioutil.WriteFile("serverconfig.yaml", bytes, 777)
	if err == nil {
		defer func() {
			os.Remove("serverconfig.yaml")
		}()
	}

	conf := GetServerConfig()

	if conf == nil {
		t.Fatal("didn't read config")
	}

	if conf.TLSCacheDir != "//tlsdata" {
		t.Fatalf("TLSCacheDir incorrect %v\n", conf.TLSCacheDir)
	}

}

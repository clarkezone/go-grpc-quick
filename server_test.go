package grpc_quick

import (
	"testing"
)

func TestCreateServer(t *testing.T) {
	srv := CreateServer()

	if srv != nil {
		t.Fatalf("there should be no config because no YAML or environment")
	}
}

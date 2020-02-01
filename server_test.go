package grpc_quick

import (
	"testing"
)

func TestCreateServer(t *testing.T) {
	c := &Conf{}
	srv := CreateServer(c)

	if srv == nil {
		t.Fatalf("there should be no config because no YAML or environment")
	}
}

func TestCreateServerFailNilConfig(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	_ = CreateServer(nil)
}

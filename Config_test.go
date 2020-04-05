package shm

import "testing"

func TestConfigGet1(t *testing.T) {
	x, err := GetConfig()
	if err != nil {
		t.Error(err.Error())
	}

	if x["idleTimeout"].(string) != "3h" {
		t.Error("not equal")
	}
}

func TestConfigGet2(t *testing.T) {
	x, err := GetConfig()
	if err != nil {
		t.Error(err.Error())
	}

	if x["idleTxxxxx"] != nil {
		t.Error("what is this?")
	}
}

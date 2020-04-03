package shm

import "testing"

func TestConfigGet1(t *testing.T) {
	x, err := ConfigGet("idleTimeout")
	if err != nil {
		t.Error(err.Error())
	}

	if x.(string) != "3h" {
		t.Error("not equal")
	}
}

func TestConfigGet2(t *testing.T) {
	x, err := ConfigGet("idleTxxxxx")
	if err != nil {
		t.Error(err.Error())
	}

	if x != nil {
		t.Error("what is this?")
	}
}

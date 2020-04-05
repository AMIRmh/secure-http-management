package shm

import (
	"testing"
	"time"
)

func TestRenewSession(t *testing.T) {
	err := CreateDataBase()
	if err != nil {
		t.Error(err.Error())
	}

	err = AddUserInfo(TEST_USER, "email", 3, []byte("pass"), []byte("salt"))
	if err != nil {
		t.Error(err.Error())
	}

	uinfo, err := GetUserInfo(TEST_USER)
	if err != nil {
		if err.Error() == ErrUserNotFound.ErrorMsg {
			t.Error(err.Error())
		} else {
			t.Error(err.Error())
		}
	}

	sm, err := NewSessionManager()
	if err != nil {
		t.Error(err.Error())
	}

	err = sm.RenewSession(uinfo.Username)
	if err != nil {
		t.Error(err.Error())
	}

	uinfo, err = GetUserInfo(TEST_USER)
	if err != nil {
		if err.Error() == ErrUserNotFound.ErrorMsg {
			t.Error(err.Error())
		} else {
			t.Error(err.Error())
		}
	}

	if uinfo.session.session == nil {
		t.Error("session is empty")
	}

	err = DelUserInfo(TEST_USER)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestNewSessionManager(t *testing.T) {
	sm, err := NewSessionManager()
	if err != nil {
		t.Error(err.Error())
	}
	config, _ := GetConfig()
	tt, err := time.ParseDuration(config["absoluteTimeout"].(string))
	if err != nil {
		t.Error(err.Error())
	}
	if sm.absoluteTimeout != tt {
		t.Error("wrong assignment")
	}
}

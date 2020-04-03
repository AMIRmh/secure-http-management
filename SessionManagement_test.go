package shm

import (
	"testing"
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

	err = RenewSession(uinfo.Username)
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

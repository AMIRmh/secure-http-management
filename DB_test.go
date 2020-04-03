package shm

import (
	"bytes"
	"testing"
)

func TestAddDelGetUpdateUserInfo(t *testing.T) {
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

	if !bytes.Equal(uinfo.Salt, []byte("salt")) ||
		!bytes.Equal(uinfo.Password, []byte("pass")) ||
		uinfo.Iterations != 3 || uinfo.Username != TEST_USER ||
		uinfo.Email != "email" {
		t.Error("wrong add user")
	}

	uinfo.Iterations = 1000

	err = UpdateUserInfo(&uinfo)
	if err != nil {
		t.Error(err)
	}

	uinfo, err = GetUserInfo(TEST_USER)
	if err != nil {
		if err.Error() == ErrUserNotFound.ErrorMsg {
			t.Error(err.Error())
		} else {
			t.Error(err.Error())
		}
	}

	if uinfo.Iterations != 1000 {
		t.Error("wrong update")
	}

	err = DelUserInfo(TEST_USER)
	if err != nil {
		t.Error(err.Error())
	}

}

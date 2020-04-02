package shm

import (
	"testing"
)

const TEST_USER = "test_user"

func TestRightPassword(t *testing.T) {
	pass, err := CreatePassword("password")
	if err != nil {
		t.Error(err.Error())
	}
	err = AddUserInfo(TEST_USER, "email", pass.Iterations, pass.Password, pass.Salt)
	if err != nil {
		t.Error(err.Error())
	}
	err = VerifyPassword("password", "user2")
	if err != nil {
		t.Error(err.Error())
	}
	err = DelUserInfo(TEST_USER)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestWrongPassword(t *testing.T) {
	pass, err := CreatePassword("password")
	if err != nil {
		t.Error(err.Error())
	}
	err = AddUserInfo(TEST_USER, "email", pass.Iterations, pass.Password, pass.Salt)
	if err != nil {
		t.Error(err.Error())
	}
	err = VerifyPassword("wrongpassword", "user2")
	if err != nil {
		if err.Error() != ErrPassWrong.ErrorMsg {
			t.Error(err.Error())
		}
	}
	err = DelUserInfo(TEST_USER)
	if err != nil {
		t.Error(err.Error())
	}
}

package shm

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"golang.org/x/crypto/pbkdf2"
)

var PEPPER = []byte{23, 44, 22, 1, 88, 45, 12, 55, 123, 244, 2, 145, 178, 23, 180,
	45, 55, 79, 68, 44, 126, 187, 139, 178, 198, 192, 21, 34, 65, 78, 22, 98}

const (
	ITERATIONS = 14000
	KEY_LEN    = 32
)

func VerifyPassword(inputPassword, username string) error {
	uinfo, err := GetUserInfo(username)
	if err != nil {
		return err
	}

	pass := pbkdf2.Key([]byte(inputPassword), append(PEPPER, uinfo.Salt...), uinfo.Iterations, KEY_LEN, sha256.New)
	if !bytes.Equal(pass, uinfo.Password) {
		return errors.New(ErrPassWrong.ErrorMsg)
	}
	return nil
}

type Password struct {
	Password   []byte
	Salt       []byte
	Iterations int
}

func CreatePassword(password string) (Password, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return Password{}, err
	}
	pass := pbkdf2.Key([]byte(password), append(PEPPER, salt...), ITERATIONS, KEY_LEN, sha256.New)
	return Password{Password: pass, Salt: salt, Iterations: ITERATIONS}, nil
}

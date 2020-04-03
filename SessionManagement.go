package shm

import (
	"crypto/rand"
	"time"
)

type Session struct {
	IdleTime     time.Time
	AbsoluteTime time.Time
	RenewalTime  time.Time
	session      []byte
}

func RenewSession(username string) error {

	uInfo, err := GetUserInfo(username)
	if err != nil {
		return err
	}

	uInfo.session.session = make([]byte, 128)
	_, err = rand.Read(uInfo.session.session)
	if err != nil {
		return err
	}

	uInfo.session.AbsoluteTime = time.Now()
	uInfo.session.IdleTime = time.Now()
	uInfo.session.RenewalTime = time.Now()

	err = UpdateUserInfo(&uInfo)
	if err != nil {
		return err
	}

	return nil
}

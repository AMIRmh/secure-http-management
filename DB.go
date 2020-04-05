package shm

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Userinfo struct {
	Username   string
	Email      string
	Password   []byte
	Salt       []byte
	Iterations int
	session    Session
}

func GetUserInfo(username string) (Userinfo, error) {
	db, err := sql.Open("sqlite3", "./db.sql")
	if err != nil {
		return Userinfo{}, err
	}
	defer db.Close()
	stmt, err := db.Prepare(`select * from userinfo`)
	if err != nil {
		return Userinfo{}, err
	}
	_, err = stmt.Exec()
	if err != nil {
		return Userinfo{}, err
	}
	row, err := db.Query(fmt.Sprintf("select * from userinfo where Username='%s'", username))
	if err != nil {
		return Userinfo{}, err
	}

	var found bool = false
	uInfo := Userinfo{}
	uInfo.session = Session{}
	for row.Next() {
		found = true
		var uid int
		var sessionBytes, sessionStructBytes []byte
		err = row.Scan(&uid, &uInfo.Username, &uInfo.Password,
			&uInfo.Email, &uInfo.Salt,
			&uInfo.Iterations, &sessionStructBytes, &sessionBytes)

		uInfo.session.username = uInfo.Username
		if sessionStructBytes != nil {
			ss := bytes.NewBuffer(sessionStructBytes)
			dec := gob.NewDecoder(ss)
			err := dec.Decode(&uInfo.session)
			uInfo.session.session = sessionBytes
			if err != nil {
				return Userinfo{}, err
			}
		}
		if err != nil {
			return Userinfo{}, err
		}
	}
	if !found {
		return Userinfo{}, errors.New(ErrUserNotFound.ErrorMsg)
	}
	return uInfo, nil
}

func AddUserInfo(username, email string, iterations int, password, salt []byte) error {
	db, err := sql.Open("sqlite3", "./db.sql")
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(`insert into userinfo (Username, Password, Email, Salt, Iterations) values (?,?,?,?,?)`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(username, password, email, salt, iterations)
	if err != nil {
		return err
	}

	return nil
}

func DelUserInfo(username string) error {
	db, err := sql.Open("sqlite3", "./db.sql")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(`delete from userinfo where username=?`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(username)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserInfo(uInfo *Userinfo) error {
	db, err := sql.Open("sqlite3", "./db.sql")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(`update userinfo set Email=?, Iterations=?, 
		Password=?, Salt=?, sessionStruct=?, session=?  where Username=?`)
	if err != nil {
		return err
	}
	var sessionBytes []byte
	if uInfo.session.session == nil {
		sessionBytes = nil
	} else {

		var sessionBuffer bytes.Buffer
		enc := gob.NewEncoder(&sessionBuffer)
		err = enc.Encode(uInfo.session)
		if err != nil {
			return err
		}
		sessionBytes = sessionBuffer.Bytes()
	}

	_, err = stmt.Exec(uInfo.Email, uInfo.Iterations, uInfo.Password, uInfo.Salt,
		sessionBytes, uInfo.session.session, uInfo.Username)
	if err != nil {
		return err
	}

	return nil
}

func CreateDataBase() error {
	db, err := sql.Open("sqlite3", "./db.sql")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(`CREATE TABLE "userinfo" (
	"uid" INTEGER PRIMARY KEY AUTOINCREMENT,
	"Username" VARCHAR(64) NOT NULL UNIQUE ,
	"Password" BLOB NOT NULL,
	"Email" VARCHAR(32) NOT NULL,
	"Salt" BLOB NOT NULL,
	"Iterations" INT NOT NULL,
	"sessionStruct" BLOB NULL,
	"session" BLOB NULL
	)`)
	if err == nil {
		_, err = stmt.Exec()
		if err != nil {
			return err
		}
	}

	return nil
}

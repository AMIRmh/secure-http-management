package shm

import (
	"database/sql"
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
	for row.Next() {
		found = true
		var uid int
		err = row.Scan(&uid, &uInfo.Username, &uInfo.Password, &uInfo.Email, &uInfo.Salt, &uInfo.Iterations)
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
	"Iterations" INT NOT NULL
	)`)
	if err == nil {
		_, err = stmt.Exec()
		if err != nil {
			return err
		}
	}

	return nil
}

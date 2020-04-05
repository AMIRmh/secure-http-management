# Secure HTTP Management

## What is this project

A simple documentation written in code for managing http 
sessions, registrations, managing
 credentials and ...
 
## How it works

this project uses a local sqlite3 
database called db.sql that stores users'
information like passwords, email,
session information and ...

To use all the features, the project
has to store your users' username
and password.

the passwords are encrypted with
[PBKDF2](https://en.wikipedia.org/wiki/PBKDF2)
. the passwords are stored with salt
and pepper.

## How to use

[Automatic Registration Prevention](#Automatic-Registration-Prevention)

[Secure Password](#Secure-Password)

### Automatic Registration Prevention

by adding ```shm.AddRegister``` with
```*http.request``` input signup
is limited to one signup in every
15 minutes by an IP.

```go
package main

import (
	"io/ioutil"
	"log"
	"net/http"
	shm "secure-http-management"
)

var (
	sm shm.SessionManager
)

func signupHandler(w http.ResponseWriter, r *http.Request) {
    //your logic to signup
    passwordStruct, err := shm.CreatePassword(password)
	checkErr(err)
	err = shm.AddUserInfo(username, email, passwordStruct.Iterations,
		passwordStruct.Password, passwordStruct.Salt)
	checkErr(err)

    sm.AddRegister(r)
}
func main() {
    sm, err := shm.NewSessionManager()
	checkErr(err)
    mux := http.NewServeMux()
	mux.Handle("/signup",
		sm.AutomaticRegistrationPrevention(http.HandlerFunc(signupHandler)))

	log.Fatal(http.ListenAndServeTLS(":4443",
		"test/server.crt", "test/server.key", mux))

}
```

### Secure Password

```shm.CreatePassword``` and 
```shm.VerifyPassword``` are
two functions to create and verify
the secure password.
these functions use local database.
```go
    passwordStruct, err := shm.CreatePassword(password)
    
    // if err == nil password is true
    err := shm.VerifyPassword(password, username)
```
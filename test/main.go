package main

import (
	"io/ioutil"
	"log"
	"net/http"
	shm "secure-http-management"
)

func loginHanlder(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	checkErr(err)

	username := r.FormValue("username")
	password := r.FormValue("password")

	err = shm.VerifyPassword(password, username)
	if err != nil {
		if err.Error() == shm.ErrPassWrong.ErrorMsg ||
			err.Error() == shm.ErrUserNotFound.ErrorMsg {
			w.WriteHeader(401)
			_, err = w.Write([]byte("username or password is wrong"))
			checkErr(err)
			return
		}
	}

}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	checkErr(err)

	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	passwordStruct, err := shm.CreatePassword(password)
	checkErr(err)
	err = shm.AddUserInfo(username, email, passwordStruct.Iterations,
		passwordStruct.Password, passwordStruct.Salt)
	checkErr(err)

	shm.AddRegister(r)
	content, _ := ioutil.ReadFile("test/login.html")
	_, err = w.Write(content)
	checkErr(err)
	w.WriteHeader(200)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	content, _ := ioutil.ReadFile("test/testPage.html")
	_, err := w.Write(content)
	checkErr(err)
	w.WriteHeader(200)
}

func main() {
	//http.HandleFunc("/login", loginHanlder)
	mux := http.NewServeMux()
	mux.HandleFunc("/", defaultHandler)
	mux.HandleFunc("/login", loginHanlder)
	mux.Handle("/signup",
		shm.AutomaticRegistrationPrevention(http.HandlerFunc(signupHandler)))

	log.Fatal(http.ListenAndServeTLS(":4443",
		"test/server.crt", "test/server.key", mux))

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

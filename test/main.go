package main

import (
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

func main() {
	//http.HandleFunc("/login", loginHanlder)
	mux := http.NewServeMux()
	mux.HandleFunc("/login", loginHanlder)

	log.Fatal(http.ListenAndServeTLS(":443",
		"server.crt", "server.key", mux))

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

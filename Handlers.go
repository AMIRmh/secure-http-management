package shm

import (
	"net/http"
	"time"
)

type registrationStruct struct {
	ip   string
	time time.Time
}

var (
	Registers []registrationStruct
)

func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func AddRegister(r *http.Request) {
	reg := registrationStruct{ip: getIP(r), time: time.Now()}
	Registers = append(Registers, reg)
}

func AutomaticRegistrationPrevention(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for index, reg := range Registers {
			if getIP(r) == reg.ip {
				if time.Now().Sub(reg.time).Minutes() < 15 {
					w.WriteHeader(406)
					return
				} else {
					next.ServeHTTP(w, r)
					Registers = append(Registers[:index],
						Registers[index+1:]...)
					break
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

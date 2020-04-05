package shm

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
	"time"
)

type registrationStruct struct {
	ip   string
	time time.Time
}

type SessionManager struct {
	registers         []registrationStruct
	activeSessions    map[string]Session
	anonymousSessions []string
	idleTimeout       time.Duration
	absoluteTimeout   time.Duration
	renewalTimeout    time.Duration
}

type Session struct {
	username     string
	IdleTime     time.Time
	AbsoluteTime time.Time
	RenewalTime  time.Time
	session      []byte
}

func NewSessionManager() (SessionManager, error) {
	config, err := GetConfig()
	if err != nil {
		return SessionManager{}, err
	}

	ss := SessionManager{}

	for key, val := range config {
		if strings.Contains(key, "Timeout") {
			valStr := val.(string)
			switch key {
			case "idleTimeout":
				ss.idleTimeout, err = time.ParseDuration(valStr)
			case "renewalTimeout":
				ss.renewalTimeout, err = time.ParseDuration(valStr)
			case "absoluteTimeout":
				ss.absoluteTimeout, err = time.ParseDuration(valStr)
			}
			if err != nil {
				return SessionManager{}, err
			}
		}
	}
	return ss, nil
}

func (sm *SessionManager) RenewSession(username string) error {

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

	sm.addOrUpdateSession(uInfo.session)
	return nil
}

func (sm *SessionManager) SessionManagement(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO set necessary headers
		token := r.Header.Get("token_id")
		tokenBytes, err := base64.StdEncoding.DecodeString(token)
		if err != nil || len(tokenBytes) != 128 {
			token, err = randomBase64()
			if err != nil {
				w.WriteHeader(501)
				return
			}
			sm.addAnonymousSession(token)
			next.ServeHTTP(w, r)

			// maybe should be changed in the future
			return
		}

	})

}

func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func (sm *SessionManager) AddRegister(r *http.Request) {
	reg := registrationStruct{ip: getIP(r), time: time.Now()}
	sm.registers = append(sm.registers, reg)
}

func (sm *SessionManager) AutomaticRegistrationPrevention(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for index, reg := range sm.registers {
			if getIP(r) == reg.ip {
				if time.Now().Sub(reg.time).Minutes() < 15 {
					w.WriteHeader(406)
					return
				} else {
					next.ServeHTTP(w, r)
					sm.registers = append(sm.registers[:index],
						sm.registers[index+1:]...)
					break
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (sm *SessionManager) addOrUpdateSession(session Session) {
	s := base64.StdEncoding.EncodeToString(session.session)
	sm.activeSessions[s] = session
}

func (sm *SessionManager) addAnonymousSession(token string) {
	sm.anonymousSessions = append(sm.anonymousSessions, token)
}

func (sm *SessionManager) invalidateSession() {

}

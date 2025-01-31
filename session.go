package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"
)

type sessionInfo struct {
	username  string
	startTime string
	endTime   string
	hash      string
}

var activeSessions map[string]sessionInfo

type sessionHandler struct {
	next http.Handler
}

type oldSessionHandler struct {
	authentication http.Handler
	home           http.Handler
}

func (os *oldSessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in old session handler")

	cookie, err := r.Cookie("session")
	usernameCookie, err1 := r.Cookie("username")
	if err != http.ErrNoCookie && err1 != http.ErrNoCookie {
		fmt.Println("got cookie!" + cookie.Value)
		session, okuser := activeSessions[usernameCookie.Value]
		if okuser != false {
			signature, _ := base64.StdEncoding.DecodeString(cookie.Value)
			ok := verifySignature([]byte(session.hash), signature, privateKeyGlobal)
			if ok == nil {
				fmt.Println("authenticated user!")
				fmt.Println("username: " + session.username + " startTime: " + session.startTime)
				os.home.ServeHTTP(w, r)
			}
		}
	}
	os.authentication.ServeHTTP(w, r)
}

func OldSessionHandler(auth, home http.Handler) *oldSessionHandler {
	return &oldSessionHandler{auth, home}
}

func (s *sessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in session handler")
	user := r.Context().Value(userKey)
	if user != nil {
		username := user.(string)
		fmt.Println("started new session")
		sessionId := getNewSessionId()
		signature, ok := generateSignature([]byte(sessionId), privateKeyGlobal)
		if ok != nil {
			log.Fatal("unable to generate signature")
		}
		fmt.Println("setting cookie" + string(signature))
		cookie := http.Cookie{Name: "session", Value: base64.StdEncoding.EncodeToString(signature)}
		cookie1 := http.Cookie{Name: "username", Value: username}
		http.SetCookie(w, &cookie)
		http.SetCookie(w, &cookie1)
		newSession := sessionInfo{username, time.Now().String(), "", getHash(sessionId)}
		activeSessions[username] = newSession
	}
	s.next.ServeHTTP(w, r)
}

func NewSessionHandler(next http.Handler) *sessionHandler {
	activeSessions = make(map[string]sessionInfo)
	return &sessionHandler{next}
}

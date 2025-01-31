package main

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const (
	userKey contextKey = "username"
)

type auth struct {
	next http.Handler
}

func (a *auth) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Println("got request")
	fmt.Println(r.URL.Path)
	fmt.Println(r.Method)

	if r.Method == http.MethodGet {
		fmt.Println("in get ")
		fmt.Println(r.URL.Path)
		if r.URL.Path == "/signup" {
			fmt.Println("trying to signup")
			http.ServeFile(w, r, "signup.html")
		} else if r.URL.Path == "/ws" || r.URL.Path == "/checkSession" {
			a.next.ServeHTTP(w, r)
		} else {
			http.ServeFile(w, r, "login.html")
		}
	} else if r.Method == http.MethodPost {
		var password, username string

		username = r.FormValue("username")
		password = r.FormValue("password")
		fmt.Println("got username: " + username)
		fmt.Println("got pass: " + password)
		//addNewPasswordDB(username, getHashFromPass(password))
		if r.URL.Path == "/signup" {
			fmt.Println("adding password")
			AddNewPassword(DbPtr, "hashes", username, getHashFromPass(password))
			ctx := context.WithValue(r.Context(), userKey, username)
			r = r.WithContext(ctx)
			a.next.ServeHTTP(w, r)
		} else {
			hash := GetUserNamePassword(DbPtr, "hashes", username)
			if compareAgainstHash(hash, password) {
				ctx := context.WithValue(r.Context(), userKey, username)
				r = r.WithContext(ctx)
				a.next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(401)
				http.ServeFile(w, r, "login.html")
			}
		}
	}

}

func NewAuth(next http.Handler) *auth {
	return &auth{next}
}

func getHashFromPass(password string) string {
	//automatically generates salt while hashing
	//second parameter is the adaptive cost factor; increase to slow it down
	//it wants the password as a byte slice, so convert using []byte()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 13)
	if err != nil {
		fmt.Printf("error generating bcrypt hash: %v\n", err)
		return ""
	}

	return string(hash)
}

func compareAgainstHash(hash, password string) bool {
	//compare a password against this hash
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	} else {
		return true
	}
}

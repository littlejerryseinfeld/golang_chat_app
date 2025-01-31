package main

import (
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

var privateKeyGlobal *rsa.PrivateKey

func serveHome(w http.ResponseWriter, r *http.Request) {

	/*
		if r.Method == "GET" {
			w.Write([]byte("Hello From Server"))
		}
	*/

	http.ServeFile(w, r, "home.html")

}

func serveWs(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal(err)
	}

	err = ws.WriteMessage(1, []byte("hello from websocket"))

	for {
		_, mess, _ := ws.ReadMessage()
		fmt.Println("got a message")
		fmt.Println(string(mess))
	}
}

var DbPtr *Db

func main() {

	privateKeyGlobal = generatePrivateKey()
	db := InitDb("/.", "newDb.db")
	defer db.closeDb()
	fmt.Printf("db ptr: %p \n", db)
	DbPtr = db
	AddTable(db, "hashes", []string{"username TEXT", "hash TEXT"})

	//AddNewPassword(db, "hashes", "test", "TESTHASHH")

	mux := http.NewServeMux()
	mux.HandleFunc("/", serveHome)
	mux.HandleFunc("/ws", serveWs)
	mux.HandleFunc("/checkSession", serveHome)
	//wrap entire mux with logger middleware
	wrappedMux := NewAuth(NewSessionHandler(mux))
	first := OldSessionHandler(wrappedMux, mux)

	err := http.ListenAndServe(":8080", first)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

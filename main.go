package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var ActiveClients = make(map[Clients]int) // map clients that are connected now
var db = InitDB("database/ChatDB")

//TODO Use the Sync.Mutex in every function where you use the map concurrently.
func main() {
	CreateUserTable()
	port := flag.Int("port", 80, "port to serve on")
	dir := flag.String("directory", "web/", "directory of web files")
	flag.Parse()
	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)
	http.Handle("/", fileHandler)
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/login", SetSessionID)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/signup", signUp)
	http.HandleFunc("/getUser", getUser)
	http.HandleFunc("/checkSession", checkSession)
	log.Printf("Running on port %d\n", *port)
	err := http.ListenAndServe(":80", nil)
	fmt.Println(err.Error())
}

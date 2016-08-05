package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type User struct {
	Username string
}

//TODO Use the Sync.Mutex in every function where you use the map concurrently.
var connections map[*websocket.Conn]bool
var Users = User{}

func main() {
	// command line flags
	port := flag.Int("port", 80, "port to serve on")
	dir := flag.String("directory", "web/", "directory of web files")
	flag.Parse()

	connections = make(map[*websocket.Conn]bool)

	// handle all requests by serving a file of the same name
	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)
	http.Handle("/", fileHandler)
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/getUser", getUsername)
	http.HandleFunc("/storeUser", storeUsername)
	log.Printf("Running on port %d\n", *port)

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	// this call blocks -- the progam runs here forever
	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
}

func getJSON(r *http.Request) map[string]interface{} {
	var data map[string]interface{}

	//	log.Printf("getJSON:\tBegin execution")
	if r.Body == nil {
		log.Printf("No Request Body")
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Printf("Error Decoding JSON")
	}
	defer r.Body.Close()
	return data
} //decode JSON

//Anything we write to w will be send to the client.  and r represents the request.
func storeUsername(w http.ResponseWriter, r *http.Request) {
	//log.Printf("Store User Handler")
	//defer log.Printf("done Store User Handler")
	data := getJSON(r)
	Users.Username = data["Username"].(string)
}
func getUsername(w http.ResponseWriter, r *http.Request) {
	//	log.Printf("Get User Handler")
	//defer log.Printf("done Get User Handler")
	json.NewEncoder(w).Encode(Users)
}
func sendAll(msg []byte) {
	for conn := range connections {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			delete(connections, conn)
			return
		}
	}
}
func wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("called")
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	connections[conn] = true
	for { //inf loop
		_, msg, err := conn.ReadMessage()
		if err != nil {
			delete(connections, conn)
			return
		}
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			fmt.Println("sent")//does this even do anything?
			return
		}
		log.Println(string(msg))
		sendAll(msg)
	}
}

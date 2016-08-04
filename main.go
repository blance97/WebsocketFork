package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)
type User struct {
	Username	string
}
//TODO Use the Sync.Mutex in every function where you use the map concurrently.
var connections map[*websocket.Conn]bool

func getJSON(r *http.Request) map[string]interface{} {
	var data map[string]interface{}

	log.Printf("getJSON:\tBegin execution")
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


func storeUsername(w http.ResponseWriter, r *http.Request){
	log.Printf("Store User Handler")
	defer log.Printf("done Store User Handler")
	data:= getJSON(r)
	//username = User{Username: data["Username"].(string)} //store in database later
	Users.Username = data["Username"].(string)
}
func getUsername(w http.ResponseWriter, r *http.Request){
	log.Prinf("Get User Handler")
	defer log.Printf("done Get User Handler")
	json.NewEncoder(w).Encode(Users)
}
func wsHandler(w http.ResponseWriter, r *http.Request) {
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
	for {//inf loop
		_, msg, err := conn.ReadMessage()
		if err != nil {
			delete(connections, conn)
			return
		}

		log.Println(string(msg))
	}
}
func main() {
	// command line flags
	port := flag.Int("port", 80, "port to serve on")
	dir := flag.String("directory", "web/", "directory of web files")
	flag.Parse()
	var Users = []User
	connections = make(map[*websocket.Conn]bool)

	// handle all requests by serving a file of the same name
	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)
	http.Handle("/", fileHandler)
	http.HandleFunc("/ws", wsHandler)

	log.Printf("Running on port %d\n", *port)

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	// this call blocks -- the progam runs here forever
	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
}

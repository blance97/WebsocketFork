package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)




// Client connection consists of the websocket and the client ip
type Clients struct {
    websocket *websocket.Conn
    IP string
}
var ActiveClients = make(map[Clients]int) // map clients that are connected now
//TODO Use the Sync.Mutex in every function where you use the map concurrently.
//var connections map[*websocket.Conn]bool

func main() {
	// command line flags
	port := flag.Int("port", 80, "port to serve on")
	dir := flag.String("directory", "web/", "directory of web files")
	flag.Parse()

	// handle all requests by serving a file of the same name
	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)
	http.Handle("/", fileHandler)
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/getUser", getUsername)
	http.HandleFunc("/storeUser", storeUsername)
	log.Printf("Running on port %d\n", *port)

	addr := fmt.Sprintf("127.0.0.1:%d", *port)

	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
}

func sendAll(ip string, msg []byte) {
	for conn,_ := range ActiveClients {
		//msg = append(msg,[]byte(ip)...)
		if err := conn.websocket.WriteMessage(websocket.TextMessage,  msg); err != nil {
			conn.websocket.Close()
		}
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Taken from gorilla's website
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	log.Println("Succesfully upgraded connection")

	socketClientIP:= r.RemoteAddr
	socketClient := Clients{conn, socketClientIP}
	ActiveClients[socketClient] = 0
	log.Println("Total clients live:", len(ActiveClients))


	for {
		// Blocks until a message is read
		_, msg, err := conn.ReadMessage()
		if err != nil {
		//delete(connections, conn)
			conn.Close()
			return
		}
		log.Println(string(msg))
		sendAll(socketClientIP,msg)
	}
}

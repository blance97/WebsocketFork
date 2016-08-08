package main

import (
	"encoding/json"
	"log"
	"net/http"
		"github.com/gorilla/websocket"
		"strings"
)

// Client connection consists of the websocket and the client ip
type Clients struct {
    websocket *websocket.Conn
    IP string
		Username string
		DateCreated int64
}
type Room struct{
	Members []Clients
	RoomName string
	Password string
}

/**
JSON Decoder
*/
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

/**
	Returns IP address of current user to the client so that it can validate whether or not a message belongs to them.
*/
func getIP(w http.ResponseWriter, r *http.Request) {
	socketClientIP := strings.Split(r.RemoteAddr, ":")
	q:=Clients{
		IP: socketClientIP[0],
	}
	json.NewEncoder(w).Encode(q)
}
/**
Function to attain the users info and then store it in to the database
*/
func getUserInfo(w http.ResponseWriter, r *http.Request) {
	//	log.Printf("Get User Handler")
	//defer log.Printf("done Get User Handler")
		socketClientIP := strings.Split(r.RemoteAddr, ":")
		data := getJSON(r)
		StoreUserInfo(socketClientIP[0], data["Username"].(string))

}
func CreateRoom(w http.ResponseWriter, r *http.Request){

}

/**
	websocket handler
*/
func wsHandler(w http.ResponseWriter, r *http.Request) {
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
	socketClient := Clients{conn, socketClientIP,"",0}// <--- Look into that
	ActiveClients[socketClient] = 0
	log.Println("Total clients live:", len(ActiveClients))


	for {
		// Blocks until a message is read
		_, msg, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			return
		}
		log.Println(string(msg))
		sendAll(socketClientIP,msg)
	}
}
func sendAll(ip string, msg []byte) {
	for conn,_ := range ActiveClients {
		//msg = append(msg,[]byte(ip)...)
		if err := conn.websocket.WriteMessage(websocket.TextMessage,  msg); err != nil {
			conn.websocket.Close()
		}
	}
}

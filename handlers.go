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
}
type NewUser struct{
	IP string
	Username string
	Password string
	DateCreated int64
}
type Room struct{
	Members []Clients
	RoomName string
	Password string
}

// func checkSessionHandler(w http.ResponseWriter, r *http.Request) {
// 	// Check if the user's session is still valid
// 	if r.Method == http.MethodGet {
// 		log.Printf("checkSessionHandler:\tBegin execution")
// 		err := checkSession(w, r)
// 		if err != nil {
// 			log.Printf("checkSessionHandler:\t%s", err.Error())
// 			w.WriteHeader(http.StatusUnauthorized)
// 			return
// 		}
// 		w.WriteHeader(http.StatusOK)
// 		return
// 	}
// 	w.WriteHeader(http.StatusBadRequest)
// 	return
// }
func checkSessionID(w http.ResponseWriter, r *http.Request){
	log.Println(r.URL.Path)
		socketClientIP := strings.Split(r.RemoteAddr, ":")
	cookie, err := r.Cookie("logged-in")
	if err==http.ErrNoCookie{
		cookie = &http.Cookie{
			Name: "logged-in",
			Value: "0",
		}
	}
	if(r.URL.Path == "/login"){
		log.Println("login")
		data:= getJSON(r)
		username:=data["Username"].(string)
		password:=data["Pass"].(string)
		 p,err:=getUserPassword(username);
		if err!=nil{
			log.Println("Error in getpassword ", err)
		}
		if password == p{
			log.Println("nigge")
			cookie = &http.Cookie{
				Name: "logged-in",
				Value: "1",
			}
			StoreUserInfo(socketClientIP[0], username,password)
		}
	}
	if r.URL.Path == "/logout"{
		cookie = &http.Cookie{
			Name: "logged-in",
			Value: "0",
		}
	}
	http.SetCookie(w,cookie)
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
func storeUserInfo(w http.ResponseWriter, r *http.Request) {
	//	log.Printf("Get User Handler")
	//defer log.Printf("done Get User Handler")
		socketClientIP := strings.Split(r.RemoteAddr, ":")
		data := getJSON(r)
		StoreUserInfo(socketClientIP[0], data["Username"].(string),data["Password"].(string))
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
			socketClientIP := strings.Split(r.RemoteAddr, ":")
	socketClient := Clients{conn, socketClientIP[0]}// <--- Look into that
	ActiveClients[socketClient] = 0
	log.Println("Total clients live:", len(ActiveClients))


	for {
		// Blocks until a message is read
		log.Println("loop")
		_, msg, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			return
		}
		 username,err2:=getUserInfo(socketClientIP[0])
		if err2!=nil{
			log.Println(err2)
			return
		}
		log.Println(username)
		msg = append(msg, "," + username...)
		log.Println(string(msg))
		sendAll(socketClientIP[0],msg)
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

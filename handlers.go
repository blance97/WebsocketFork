package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"time"
)

// Client connection consists of the websocket and the client ip
type Clients struct {
	websocket *websocket.Conn
	IP        string
}
type NewUser struct {
	IP          string
	Username    string
	Password    string
	SessionID 	string
	DateCreated int64
}
type Room struct {
	Members  []Clients
	RoomName string
	Password string
}
type Cookie struct {
	Name       string
	Value      string
	Path       string
	Domain     string
	Expires    time.Time
	RawExpires string
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int
	Secure   bool
	HttpOnly bool
	Raw      string
	Unparsed []string // Raw text of unparsed attribute-value pairs
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

func SetSessionID(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	socketClientIP := strings.Split(r.RemoteAddr, ":")
	cookie, err := r.Cookie("logged-in")
	username := r.FormValue("Username")
	password := r.FormValue("password")
	if err == http.ErrNoCookie {
		cookie = &http.Cookie{
			Name:  "logged-in",
			Value: "0",
		}
	}
	if r.URL.Path == "/login" {
		log.Println("login")
		p, err := getUserPassword(username)
		if err != nil {
			log.Println("Error in getpassword ", err)
		}
		if password == p {
			log.Println("nigge")
			expiration := time.Now().Add(365 * 24 * time.Hour)
			cookie = &http.Cookie{Name: "logged-in", Value: "1", Expires: expiration}
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/chat.html", 302)
			StoreUserInfo(socketClientIP[0], username, password, cookie.Value)
		} else {
			http.Redirect(w, r, "/", 302)
		}
	}
	/**
	TODO: GET Username from sessionID and then set that sessionid to zero.
	*/
	if r.URL.Path == "/logout" {
		cookie = &http.Cookie{
			Name:  "logged-in",
			Value: "0",
		}
		StoreUserInfo(socketClientIP[0], username, password, cookie.Value)
	}
	http.SetCookie(w, cookie)
}
/**
	 checks the SessionID
*/
func checkSession(w http.ResponseWriter, r *http.Request){
	cookie, _ := r.Cookie("logged-in")
	if(cookie.Value!="0"){
			w.WriteHeader(http.StatusOK)
			return
	}
		http.Error(w, "No Session", 403)
		http.Redirect(w, r, "/", 302)
		return
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
	q := Clients{
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
	StoreUserInfo(socketClientIP[0], data["Username"].(string), data["Password"].(string), "0")
}
func CreateRoom(w http.ResponseWriter, r *http.Request) {

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
	socketClient := Clients{conn, socketClientIP[0]} // <--- Look into that
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
		username, err2 := getUserInfo(socketClientIP[0])
		if err2 != nil {
			log.Println(err2)
			return
		}
		log.Println(username)
		msg = append(msg, ","+username...)
		log.Println(string(msg))
		sendAll(socketClientIP[0], msg)
	}
}
func sendAll(ip string, msg []byte) {
	for conn, _ := range ActiveClients {
		//msg = append(msg,[]byte(ip)...)
		if err := conn.websocket.WriteMessage(websocket.TextMessage, msg); err != nil {
			conn.websocket.Close()
		}
	}
}

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
type User struct {
	IP          string
	Username    string
	Password    string
	SessionID   string
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

var SessionToken string

func SetSessionID(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	socketClientIP := strings.Split(r.RemoteAddr, ":")
	cookie, err := r.Cookie("SessionToken")
	SessionToken = cookie.Value
	username := r.FormValue("Username")
	password := r.FormValue("password")
	if err == http.ErrNoCookie {
		cookie = &http.Cookie{
			Name:  "SessionToken",
			Value: "0",
		}
	}
	if r.URL.Path == "/login" {
		log.Println("login")
		redirectTarget := "/"
		p, err := getUserPassword(username)
		if err != nil {
			log.Println("Error in getpassword ", err)
			http.Redirect(w, r, "/login.html", 302)
		}

		var token string
		if password == p {
			log.Printf("Password1: %s Password2: %s", password, p)
			expiration := time.Now().Add(365 * 24 * time.Hour)
			for {
				token, _ = GenerateRandomString(64)
				if CheckValidSessionToken(token) {
					break
				}
			}
			cookie := &http.Cookie{Name: "SessionToken", Value: token, Expires: expiration}
			http.SetCookie(w, cookie)
			redirectTarget = "/chat.html"
			StoreUserInfo(socketClientIP[0], username, password, SessionToken)
		}
		http.Redirect(w, r, redirectTarget, 302)
		http.SetCookie(w, cookie)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	log.Println("Logout Hanlder")
	username, _ := getUsername(SessionToken)
	log.Println("Cleared Logout")
	cookie := &http.Cookie{
		Name:  "SessionToken",
		Value: "0",
	}
	http.SetCookie(w, cookie)
	storeNewSessionToken(cookie.Value, username)
}

/**
checks the SessionID
*/
func checkSession(w http.ResponseWriter, r *http.Request) {
	if SessionToken != "0" {
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Error(w, "No Session", 403)
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
Returns User so that it can validate whether or not a message belongs to them.
*/
func getUser(w http.ResponseWriter, r *http.Request) {
	username, _ := getUsername(SessionToken)
	q := User{
		Username: username,
	}
	json.NewEncoder(w).Encode(q)
}

/**
Function to attain the users info and then store it in to the database
*/
func signUp(w http.ResponseWriter, r *http.Request) {
	//	log.Printf("Get User Handler")
	//defer log.Printf("done Get User Handler")
	socketClientIP := strings.Split(r.RemoteAddr, ":")
	data := getJSON(r)
	StoreUserInfo(socketClientIP[0], data["Username"].(string), data["Pass"].(string), "0")
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
		cookie, err := r.Cookie("SessionToken")
		if err != nil {
			log.Println("Lolcde", err)
			return
		}
		SessionToken = cookie.Value
		username, _ := getUsername(SessionToken)
		log.Println("SessionToken: ", SessionToken)
		msg := []byte(username + ": ")
		_, msg2, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			return
		}

		msg = append(msg, msg2...)
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

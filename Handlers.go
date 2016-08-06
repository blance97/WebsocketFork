package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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
//	data := getJSON(r)
//	Users.Username = data["Username"].(string)
}
func getUsername(w http.ResponseWriter, r *http.Request) {
	//	log.Printf("Get User Handler")
	//defer log.Printf("done Get User Handler")
  	socketClientIP:= r.RemoteAddr
    c := Clients{
      IP: socketClientIP,
    }
	json.NewEncoder(w).Encode(c)
}

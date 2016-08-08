package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
  "time"
)

var dbMu sync.Mutex

/**
Initialize db tables with file path
@param String filepath
@Returns *sql.DB
*/
func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	log.Println(filepath)
	if err != nil {
		log.Print(err)
	}
	if db == nil {
		log.Print("db nil")
	}
	log.Println("Successfull opened db")
	return db
}

/**
Creates the user table with the following paramters:
IP: String
Name: string
*/
func CreateUserTable() {
	sql_table := `
	CREATE TABLE IF NOT EXISTS Users(
		IP TEXT PRIMARY KEY,
		Name TEXT,
    DateCreated
	);
	`
	dbMu.Lock()
	_, err := db.Exec(sql_table)
	dbMu.Unlock()
	if err != nil {
		log.Print(err)
	}
}

func StoreUserInfo(socketClientIP string, Username string) {
	sql_stmt := `
	INSERT OR REPLACE INTO Users(
		IP,
		Name,
    DateCreated
	)values(?, ?, ?)
	`
	 stmt, err := db.Prepare(sql_stmt)
    if err!=nil{
		log.Print(err)
	}
	c := Clients{
		IP:          socketClientIP,
		Username:    Username,
		DateCreated: time.Now().Unix(),
	}
	if _, err := stmt.Exec(c.IP, c.Username, c.DateCreated); err != nil {
		log.Println(err)
	}
}
func getUserInfo(socketClientIP string)(string,error){
	var ip string
		sql_stmt := "SELECT Name FROM Users WHERE IP = $1"
		if err:=db.QueryRow(sql_stmt, socketClientIP).Scan(&ip); err!=nil{
			return "",err
		}
		return ip,nil
}

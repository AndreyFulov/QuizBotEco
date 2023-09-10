package main

import (
	"os"
	"time"
)

var host = os.Getenv("HOST")
var port = os.Getenv("PORT")
var user = os.Getenv("USER")
var password = os.Getenv("PASSWORD")
var dbname = os.Getenv("DBNAME")
var sslmode = os.Getenv("SSLMODE")
var token = os.Getenv("TOKEN")

func main() {
	var db DataBase
	db.InitInfo(host,port,user,password,dbname,sslmode)
	time.Sleep(5 * time.Second)
	bot := NewBot(token, &db)
	bot.Bot()
}
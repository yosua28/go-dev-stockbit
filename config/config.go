package config

// General
var LimitQuery uint64 = 100
var CORSAllowOrigin = []string{"*"}

// SSH
var SSHUser string = "yoga"
var SSHPassword string = "Welcome.0303"
var SSHHost string = "172.31.214.9"

// DB
var DBUser string = "root"
var DBPassword string = "Mncam777"
var DBName string = "mam_core"
var DBHost string = "localhost:3306"

// Email
// var EmailSMTPHost string = "smtp.gmail.com"
// var EmailSMTPPort uint64 = 587
// var EmailFrom string = "gameraja82@gmail.com"
// var EmailFromPassword string = "kmzway87aa"
var EmailSMTPHost string = "172.17.20.124"
var EmailSMTPPort uint64 = 25
var EmailFrom string = "cso.mam@mncgroup.com"
var EmailFromPassword string = "Welcome2MNCAM!"

// Session
var SessionExpired int64 = 10000
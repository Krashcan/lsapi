package main

import(
	"fmt"
	"log"
	"net/http"
	"database/sql"
	"os"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error
func main() {
	db,err = sql.Open("mysql","root:ATul1996@@@/fdata")
	if err!=nil{
		log.Fatal("Opening database: ",err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	http.HandleFunc("/",Home)
	http.HandleFunc("/login",Login)
	http.HandleFunc("/signup",SignUp)
	http.ListenAndServe(":8080",nil)
}

func Home(w http.ResponseWriter,r *http.Request){
	if r.Method =="GET"{
		t,err := template.ParseFiles(os.Getenv("GOPATH")+"/src/github.com/krashcan/lsapi/template/index.html")
		if err!=nil{
			log.Fatal("Parse File: ",err)
		}
		t.Execute(w,nil)	
	}else{
		if r.URL.Path=="/"{
			var dbPassword string
			err = db.QueryRow("SELECT password FROM userinfo WHERE username = ?",r.FormValue("username")).Scan(&dbPassword)
			if err!=nil{
				t,err := template.ParseFiles(os.Getenv("GOPATH")+"/src/github.com/krashcan/lsapi/template/login.html")
				if err!=nil{
					log.Fatal("Parse file: ",err)
				}
				t.Execute(w,"Wrong Username")
				return
			}
			err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(r.FormValue("password")))
			if err!=nil{
				t,err := template.ParseFiles(os.Getenv("GOPATH")+"/src/github.com/krashcan/lsapi/template/login.html")
				if err!=nil{
					log.Fatal("Parse file: ",err)
				}
				t.Execute(w,"Wrong Password")
				return
			}
		}else if r.URL.Path=="/add"{
			_,err := db.Exec("INSERT INTO user_players(username,player) VALUES(?,?)",r.FormValue("username"),r.FormValue("player"))
			if err!=nil{
				fmt.Fprintf(w,"Trouble connecting with database.")
				fmt.Println(err)
				return
			}
		}
		
		var players []string

		rows, err := db.Query(`SELECT player FROM user_players WHERE username=?`,r.FormValue("username"))
		defer rows.Close()
   		if err != nil {
   			log.Fatal(err) 
   		}
   		var player string
   		for i := 0; rows.Next(); i++ {
        	err := rows.Scan(&player)
        	if err != nil {
        		log.Fatal(err) 
        	}

        	players = append(players,player)
    	}

    	var userPlayers struct{
    		Name string
    		Players []string
    	}
    	userPlayers.Name = r.FormValue("username")
    	userPlayers.Players = players

    	t,err := template.ParseFiles(os.Getenv("GOPATH")+"/src/github.com/krashcan/lsapi/template/profile.html")
    	if err!=nil{
    		log.Fatal(err)
    	}
    	t.Execute(w,userPlayers)
	}	
}

func Login(w http.ResponseWriter,r *http.Request){
	t,err := template.ParseFiles(os.Getenv("GOPATH")+"/src/github.com/krashcan/lsapi/template/login.html")
	if err!=nil{
		log.Fatal("Parse file: ",err)
	}
	if r.Method =="GET"{
		t.Execute(w,nil)	
	}else{
		name := r.FormValue("username")
		pw := r.FormValue("password")
		var user string
		err = db.QueryRow("SELECT username FROM userinfo WHERE username =?",name).Scan(&user)

		switch{
			case err == sql.ErrNoRows:
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
				if err != nil {
					http.Error(w, "Server error, unable to create your account.(Encryption Failed)", 500)
					return
				}

				_, err = db.Exec("INSERT INTO userinfo(username, password) VALUES(?, ?)", name, hashedPassword)
				if err != nil {
					fmt.Println("query",err)
					http.Error(w, "Server error, unable to create your account.(Database Insertion Failed)", 500)
					return
				}
				t.Execute(w,"Signup done.Login to Continue.")
				return
			case err != nil:
				http.Error(w, "Server error, unable to create your account.(Database query failed)", 500)
				return
			default:
				t,err = template.ParseFiles(os.Getenv("GOPATH")+"/src/github.com/krashcan/lsapi/template/signup.html")
				if err!=nil{
					log.Fatal(err)
				}
				t.Execute(w,"User Exists. Please Try with a different username")
		}	
	}
	
}

func SignUp(w http.ResponseWriter,r *http.Request){

	t,err := template.ParseFiles(os.Getenv("GOPATH")+"/src/github.com/krashcan/lsapi/template/signUp.html")
	if err!=nil{
		log.Fatal("Parse file: ",err)
	}
	t.Execute(w,nil)	
}


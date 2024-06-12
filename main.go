package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type errors struct {
	UsernameError string
	PasswordError string
}
type home struct {
	Username2 string
}

var tmpl *template.Template
var errorV errors
var h home
var sessions = make(map[string]string)
var sessionID string
var sessionCookie http.Cookie
var c http.Cookie

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/logout", logoutHandler)

	fmt.Printf("Starting Server At port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store")

	_, err := r.Cookie("cookie")
	if err == nil { //means cookie here
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPost {
		w.Header().Set("Cache-Control", "no-cache,no-store")

		var UserName, PassWord = "sudhin.A", "sudhin123"
		fmt.Println("Success in login handler")
		if err := r.ParseForm(); err != nil {
			fmt.Println("Error here", err)
			http.Error(w, "Failed to Parse FormData", http.StatusInternalServerError)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")
		fmt.Println(username)
		fmt.Println(password)

		if UserName != username && PassWord == password {
			errorV.UsernameError = "Invalid Username"
			errorV.PasswordError = ""
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		if PassWord != password && UserName == username {
			errorV.PasswordError = "Invalid Password"
			errorV.UsernameError = ""
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		if UserName != username && PassWord != password {
			errorV.UsernameError = "Invalid Username"
			errorV.PasswordError = "Invalid Password"
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		if UserName == username && PassWord == password {
			fmt.Println("Success in LoginHandler")

			//Set the new cookie with the expiration time in the past
			sessionID = strconv.FormatInt(rand.Int63(), 16)

			sessionCookie = http.Cookie{Name: "cookie", Expires: time.Now().Add(time.Hour * 1), Value: sessionID}
			http.SetCookie(w, &sessionCookie)

			sessions[sessionID] = username

			fmt.Printf("Sessions created with session id %s and session data %v\n", sessionID, sessions[sessionID])

			http.Redirect(w, r, "/home", http.StatusSeeOther)
		}
	}

	tmpl, err = template.ParseFiles("templates/login.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.ExecuteTemplate(w, "login.html", errorV)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store")

	cookie, err := r.Cookie("cookie")
	if err != nil {
		fmt.Print("--------", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tmpl, err = template.ParseFiles("templates/home.html")
	if err != nil {
		log.Fatal(err)
	}
	sessionID = cookie.Value
	h.Username2 = sessions[sessionID]
	fmt.Println("____", h)
	tmpl.ExecuteTemplate(w, "home.html", h)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache,no-store")

	c = http.Cookie{Name: "cookie", Value: "", Expires: time.Now().AddDate(0, 0, -1), MaxAge: -1}

	http.SetCookie(w, &c)

	delete(sessions, sessionID)

	fmt.Println("in logoutHandler")
	fmt.Printf("%v\n", sessions)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

package handlers

import (
	"net/http"
	"html/template"
	"fmt"
)

var templates = template.Must(template.ParseFiles("view/templates/home.html", "view/templates/login.html", "view/templates/register.html", "view/templates/main.html"))
var uC = UserController{}

func homeHandler(w http.ResponseWriter, r *http.Request){
	templates.ExecuteTemplate(w, "home.html", nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodGet{
		templates.ExecuteTemplate(w, "register.html", nil)
		return
	}
	if r.Method == http.MethodPost{
		/*firstname := r.PostFormValue("firstname")
		lastname := r.PostFormValue("lastname")
		username := r.PostFormValue("username")
		email := r.PostFormValue("email")
		password := r.PostFormValue("password")
		isValid, message := userRegisterDataValidation([]string{firstname, lastname, username, email, password})
		if !isValid{
			templates.ExecuteTemplate(w, "register.html", message)
			return
		}
		uC.Add([]string{firstname, lastname, username, email, password})
		sess, err := globalSessions.SessionStart(w, r)
		if err != nil{
			http.Redirect(w, r, "/main/", http.StatusInternalServerError)
		}
		sess.Set("username", username)
		sess.Set("firstName", firstname)
		sess.Set("lastName", lastname)
		sess.Set("email", email)*/
		//TEST SESSIONS
		fmt.Println(globalSessions)
		sess, err := globalSessions.SessionStart(w, r)
		if err != nil{
			fmt.Println(err)
		}
		sess.Set("username", "username")
		sess.Set("email", "mail@mail.ru")
		http.Redirect(w, r, "/main/", http.StatusFound)
		return
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodGet{
		templates.ExecuteTemplate(w, "login.html", nil)
		return
	}
	if r.Method == http.MethodPost{
		username := r.PostFormValue("username")
		//password := r.PostFormValue("password")
		//TODO: Username and password validation
		sess, err := globalSessions.SessionStart(w, r)
		if err != nil{
			http.Redirect(w, r, "/main/", http.StatusInternalServerError)
		}
		sess.Set("username", username)
		http.Redirect(w, r, "/main/", http.StatusFound)
	}
}

func mainHandler (w http.ResponseWriter, r *http.Request){
	templates.ExecuteTemplate(w, "main.html", nil)
}
package handlers

import (
	"net/http"
	"html/template"
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
		firstname := r.PostFormValue("firstname")
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
		sess := globalSessions.SessionStart(w, r)
		sess.Set("username", username)
		sess.Set("firstName", firstname)
		sess.Set("lastName", lastname)
		sess.Set("email", email)
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
		sess := globalSessions.SessionStart(w, r)
		sess.Set("username", username)
		http.Redirect(w, r, "/main/", http.StatusFound)
	}
}

func mainHandler (w http.ResponseWriter, r *http.Request){
	templates.ExecuteTemplate(w, "main.html", nil)
}
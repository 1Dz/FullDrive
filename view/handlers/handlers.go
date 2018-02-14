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
		sess, err := globalSessions.SessionStart(w, r)
		if err != nil{
			http.Redirect(w, r, "/main/", http.StatusInternalServerError)
		}
		sess.Set("username", username)
		sess.Set("firstName", firstname)
		sess.Set("lastName", lastname)
		sess.Set("email", email)
		globalSessions.Driver.SessionUpdate(sess.SessionID())
		http.Redirect(w, r, "/", http.StatusFound)
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
			http.Redirect(w, r, "/home/", http.StatusInternalServerError)
		}
		if len(sess.Values()) == 0 {
			u := uC.GetByName(username)
			sess.Set("firstName", u.FirstName)
			sess.Set("lastName", u.LastName)
			sess.Set("email", u.Email)
			sess.Set("username", u.Username)
		}
		http.Redirect(w, r, "/main/", http.StatusFound)
	}
}

func mainHandler (w http.ResponseWriter, r *http.Request){
	s, e := globalSessions.SessionStart(w, r)
	if e != nil{
		http.Error(w, "Session Error", http.StatusInternalServerError)
	}
	username, ok := s.Values()["username"]
	if !ok {
		http.Error(w, "Wrong username", http.StatusNotFound)
	}
	u := uC.GetByName(username.(string))
	templates.ExecuteTemplate(w, "main.html", u)
}
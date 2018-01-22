package handlers

import (
	"net/http"
	"html/template"
	"encoding/json"
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
		//fmt.Fprint(w, firstname + " " + lastname + " " + username + " " + email + " " + password)
		isValid, message := userRegisterDataValidation([]string{firstname, lastname, username, email, password})
		if !isValid{
			templates.ExecuteTemplate(w, "register.html", message)
			return
		}
		uC.Add([]string{firstname, lastname, username, email, password})
		w.Header().Set("Content-Type", "application/json")
		js, _ := json.Marshal(uC.GetByName(username))
		templates.Execute(w, js)
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
		us := uC.GetByName(username)
		w.Header().Set("Content-Type", "application/json")
		js, _ := json.Marshal(us)
		templates.Execute(w, js)
		http.Redirect(w, r, "/main/", http.StatusFound)
	}
}

func mainHandler (w http.ResponseWriter, r *http.Request){
	templates.ExecuteTemplate(w, "main.html", nil)
}
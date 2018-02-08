package handlers

import (
	"net/http"
	"Conus/model"
	"Conus/persistence"
)

var globalSessions, _ = persistence.NewManager("pgm", 3600)
var pages = map[string] func(w http.ResponseWriter, r *http.Request){
	"/": homeHandler,
	"/register/": registerHandler,
	"/login/": loginHandler,
	"/main/": mainHandler,
}

type Controller interface{
	GetAll() *[]model.User
	GetByName(s string) *model.User
	GetById(f float64) *model.User
	Add(m []string) int
	Update(m []string)
	Delete(f float64)
}

type UserController struct{

}

type MyMux struct{

}

func Init(){
	persistence.Init()
	mux := &MyMux{}
	http.ListenAndServe(":8080", mux)
	go globalSessions.SessionGC()
}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request){
	if pages[r.URL.Path] != nil{
		pages[r.URL.Path](w, r)
		return
	}
	pages["/"](w, r)
}

func (c *UserController) GetAll() *[]model.User{
	r, e := persistence.GetAllUsers()
	checkError(e)
	return r
}

func (c *UserController) GetByName(s string) *model.User{
	r, e := persistence.GetUserByName(s)
	checkError(e)
	return &r
}

func (c *UserController) GetById(f float64) *model.User{
	r, e := persistence.GetUserById(f)
	checkError(e)
	return &r
}

func (c *UserController) Add(m []string) {
	e := persistence.AddUser(m)
	checkError(e)
}

func (c *UserController) Update(m []string) {
	e := persistence.UpdateUser(m)
	checkError(e)
}

func (c *UserController) Delete(f float64) {
	e := persistence.DeleteUser(f)
	checkError(e)
}

func checkError(e error) bool{
	if e != nil{
		panic(e.Error())
	}
	return true
}
package main

import (
	"Conus/view/handlers"
	"net/http"
	"Conus/persistence"
)

func main(){
	persistence.Init()
	mux := &handlers.MyMux{}
	http.ListenAndServe(":8080", mux)
}

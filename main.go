package main

import (
	"github.com/rmukubvu/login-service/handler"
	"log"
	"net/http"
)

func main() {
	r:= handler.InitRouter()
	log.Fatal(http.ListenAndServe(":8001",r))
}

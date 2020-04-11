package main

import (
	"github.com/gorilla/handlers"
	"github.com/rmukubvu/login-service/handler"
	"log"
	"net/http"
)

func main() {
	r := handler.InitRouter()
	log.Fatal(http.ListenAndServe(":8001", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)))
}

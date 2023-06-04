package main

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/mpgallage/xmcrud/database"
	"github.com/mpgallage/xmcrud/events"
	"github.com/mpgallage/xmcrud/handlers"
	"github.com/mpgallage/xmcrud/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	database.Init()
	defer database.Close()
	events.Init()
	defer events.Close()

	r := mux.NewRouter()
	r.HandleFunc("/company", middleware.ValidateMiddleware(handlers.CreateCompanyHandler)).Methods("POST")
	r.HandleFunc("/company/{id}", middleware.ValidateMiddleware(handlers.GetCompanyHandler)).Methods("GET")
	r.HandleFunc("/company/{id}", middleware.ValidateMiddleware(handlers.UpdateCompanyHandler)).Methods("PATCH")
	r.HandleFunc("/company/{id}", middleware.ValidateMiddleware(handlers.DeleteCompanyHandler)).Methods("DELETE")
	r.HandleFunc("/authenticate", handlers.CreateToken).Methods("POST")

	http.Handle("/", r)
	err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")), middleware.JsonContentTypeMiddleware(r))
	if err != nil {
		log.Fatal("Error starting server.", err)
		return
	}
}

package main

import (
	handler "finalexam/Handler"
	"net/http"

	"github.com/gorilla/mux"
)

var PORT = ":8088"

func main() {

	r := mux.NewRouter()

	userRoute := r.PathPrefix("/user").Subrouter()
	userRoute.HandleFunc("", handler.UsersHandler)
	userRoute.HandleFunc("/{id}", handler.UsersHandler)
	userRoute.Use(handler.MiddlewareAuth)

	socmedRoute := r.PathPrefix("/socialmedias").Subrouter()
	socmedRoute.HandleFunc("", handler.SocmedHandler)
	socmedRoute.HandleFunc("/{id}", handler.SocmedHandler)
	socmedRoute.Use(handler.MiddlewareAuth)

	photoRoute := r.PathPrefix("/photos").Subrouter()
	photoRoute.HandleFunc("", handler.PhotoHandler)
	photoRoute.HandleFunc("/{id}", handler.PhotoHandler)
	socmedRoute.Use(handler.MiddlewareAuth)

	commentRoute := r.PathPrefix("/comments").Subrouter()
	commentRoute.HandleFunc("", handler.CommentHandler)
	commentRoute.HandleFunc("/{id}", handler.CommentHandler)
	socmedRoute.Use(handler.MiddlewareAuth)

	r.HandleFunc("/login", handler.UserLogin).Methods("POST")
	r.HandleFunc("/register", handler.CreateUser).Methods("POST")
	http.Handle("/", r)
	http.ListenAndServe(PORT, nil)

}

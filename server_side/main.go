package main

import (
	"fmt"
	"github.com/Jla3eP/tetris/server_side/auth/handling"
	"github.com/gorilla/mux"
	"net/http"
)

const address = ":1234"

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/register", handling.CreateUser).Methods(http.MethodPost)

	server := &http.Server{
		Addr:    address,
		Handler: router,
	}
	err := server.ListenAndServeTLS("./client_side/localhost.crt", "./client_side/localhost.key")
	if err != nil {
		fmt.Println(err)
	}
}

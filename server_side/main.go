package main

//USE run_server.sh !!!

import (
	"fmt"
	"github.com/Jla3eP/tetris/server_side/auth/handling"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

const address = ":1234"

func main() {
	runServer()
}

func runServer() {
	router := mux.NewRouter()
	handleFuncs(router)

	server := &http.Server{
		Addr:    address,
		Handler: router,
	}
	fmt.Println("try to run server")
	stopCh := make(chan struct{})
	go func(stopCh chan struct{}) {
		time.Sleep(time.Second)
		select {
		case <-stopCh:
			fmt.Println("failed to run server")
		default:
			fmt.Println("server is running")
		}
	}(stopCh)
	err := server.ListenAndServeTLS("./keys/localhost.crt", "./keys/localhost.key")
	stopCh <- struct{}{}
	if err != nil {
		fmt.Println(err)
	}
}

func handleFuncs(router *mux.Router) {
	router.HandleFunc("/register", handling.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/logIn", handling.LogIn).Methods(http.MethodPost)
	router.HandleFunc("/updateSession", handling.UpdateSessionTimeout).Methods(http.MethodPost)
	router.HandleFunc("/findGame", handling.FindGame).Methods(http.MethodPost)
	router.HandleFunc("/getGameInfo", handling.GetGameInfo).Methods(http.MethodPost, http.MethodGet)
}

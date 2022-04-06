package main

import (
	"github.com/Jla3eP/tetris/server_side/auth/handling"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/register", handling.CreateUser)
	log.Fatalln(http.ListenAndServe("localhost:1234", mux))

	/*usr := User.User{NickName: "TestUser"}
	err := database.CreateUser(context.Background(), usr, "1234qwer")
	if err != nil {
		fmt.Println(err.Error())
	}

	if ok, err := database.VerifyPassword(context.Background(), usr, "1234qwer"); !ok {
		fmt.Println(err.Error())
	}*/
}

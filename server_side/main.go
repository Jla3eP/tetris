package main

import (
	"context"
	"fmt"
	"github.com/Jla3eP/tetris/server_side/auth/User"
	"github.com/Jla3eP/tetris/server_side/auth/database"
)

func main() {
	usr := User.User{NickName: "TestUser"}
	err := database.CreateUser(context.Background(), usr, "1234qwer")
	if err != nil {
		fmt.Println(err.Error())
	}

	if ok, err := database.VerifyPassword(context.Background(), usr, "1234qwer"); !ok {
		fmt.Println(err.Error())
	}
}

package main

import (
	"fmt"
	"github.com/Jla3eP/tetris/both_sides_code"
	"github.com/Jla3eP/tetris/client_side/requests"
	"os"
)

// 1 - login, 2 - password

func main() {
	args := os.Args
	fmt.Println(len(args))
	fmt.Println(args)

	if len(args) != 3 {
		fmt.Println("invalid input")
		return
	}
	info := &both_sides_code.AuthInfo{}
	info.Nickname = args[1]
	info.Password = args[2]

	err := requests.Register(info)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

}

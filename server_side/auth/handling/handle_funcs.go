package handling

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Jla3eP/tetris/server_side/auth/User"
	"github.com/Jla3eP/tetris/server_side/auth/database"
	"net/http"
	"unicode/utf8"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var buffer []byte

	_, err := r.Body.Read(buffer) //тут пусто, хмм

	fmt.Println(string(buffer))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userInfo := RegistationRequest{}
	json.Unmarshal(buffer, userInfo)

	username := userInfo.Nickname
	fmt.Println(username)
	if username == "" || utf8.RuneCountInString(username) < 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("nickname is too short\n"))
		return
	}

	password := userInfo.Password
	if password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if utf8.RuneCountInString(password) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("password is too short\n"))
		return
	}

	err = database.CreateUser(context.Background(), User.User{NickName: username}, password)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("%s, your account was created\n", username)))
}

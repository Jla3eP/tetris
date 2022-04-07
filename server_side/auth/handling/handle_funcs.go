package handling

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Jla3eP/tetris/server_side/auth/User"
	"github.com/Jla3eP/tetris/server_side/auth/database"
	"io/ioutil"
	"net/http"
	"unicode/utf8"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error() + "\n"))
		return
	}

	user := RegistationRequest{}
	err = json.Unmarshal(requestBody, &user)
	if err != nil {
		fmt.Println(err)
	}

	if user.Nickname == "" || utf8.RuneCountInString(user.Nickname) < 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("nickname is too short\n"))
		return
	}

	if user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if utf8.RuneCountInString(user.Password) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("password is too short\n"))
		return
	}

	err = database.CreateUser(context.Background(), User.User{NickName: user.Nickname}, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error() + "\n"))
		return
	}
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("%s, your account was created\n", user.Nickname)))
}

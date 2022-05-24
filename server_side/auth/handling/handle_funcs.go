package handling

import (
	"context"
	"fmt"
	"github.com/Jla3eP/tetris/server_side/auth/database"
	"github.com/Jla3eP/tetris/server_side/auth/user"
	"net/http"
	"sync"
	"time"
	"unicode/utf8"
)

const (
	sessionTimeout = 10 * time.Second
)

var (
	sessions        map[string]sessionValues
	sessionsCleaner = sync.Once{}
)

func FindGame(w http.ResponseWriter, r *http.Request) {

}

func UpdateSessionTimeout(w http.ResponseWriter, r *http.Request) {
	sessionsCleaner.Do(func() {
		go func() {
			ticker := time.NewTicker(1 * time.Second)
			for {
				for k, v := range sessions {
					<-ticker.C
					if time.Now().After(v.lastUpdate.Add(sessionTimeout)) {
						delete(sessions, k)
					}
				}
			}
		}()
	})

	info, err := getSessionInfoFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := database.GetIdByUsername(info.Nickname)
	if err != nil {
		return
	}

	session, ok := sessions[info.SessionKey]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if session.username != info.Nickname || session.userAgent != r.UserAgent() || session.id != id {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	session.lastUpdate = time.Now()
	sessions[info.SessionKey] = session
	w.WriteHeader(http.StatusOK)
}

func LogIn(w http.ResponseWriter, r *http.Request) {
	createdAt := time.Now()
	authInfo, err := getAuthInfoFromRequest(w, r)
	if err != nil {
		w.Write([]byte(err.Error() + "\n"))
	}
	usr := user.User{}
	usr.NickName = authInfo.Nickname

	ok, err := database.VerifyPassword(context.Background(), usr, authInfo.Password)
	if err != nil || !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	key, err := createSessionKey(r, usr.NickName, usr.ID, createdAt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(key))
	w.WriteHeader(http.StatusOK)
	sessions[key] = sessionValues{
		id:         usr.ID,
		username:   usr.NickName,
		userAgent:  r.UserAgent(),
		lastUpdate: time.Now(),
		createdAt:  createdAt,
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	usr, err := getAuthInfoFromRequest(w, r)

	if usr.Nickname == "" || utf8.RuneCountInString(usr.Nickname) < 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("nickname id too short\n"))
		return
	}

	if usr.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if utf8.RuneCountInString(usr.Password) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("password is too short\n"))
		return
	}

	if err = database.CreateUser(context.Background(), user.User{NickName: usr.Nickname}, usr.Password); err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error() + "\n"))
		return
	}
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("%s, your account was created\n", usr.Nickname)))
}

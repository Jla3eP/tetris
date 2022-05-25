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
	sessionsMu      = &sync.RWMutex{}
	sessions        = make(map[string]sessionValues)
	sessionsCleaner = sync.Once{}

	playersInQueue = make(map[string]struct{})
	queueMu        = &sync.RWMutex{}
)

func FindGame(w http.ResponseWriter, r *http.Request) {
	if info, err := auth(w, r); err == nil {
		queueMu.Lock()
		playersInQueue[info.Nickname] = struct{}{}
		queueMu.Unlock()
	}
}

func UpdateSessionTimeout(w http.ResponseWriter, r *http.Request) {
	sessionsCleaner.Do(func() {
		go func() {
			ticker := time.NewTicker(1 * time.Second)
			for {
				<-ticker.C //wait a second. who are you?
				sessionsMu.Lock()
				for k, v := range sessions {
					if time.Now().After(v.lastUpdate.Add(sessionTimeout)) {
						delete(sessions, k)

						queueMu.RLock()
						_, ok := playersInQueue[v.username]
						queueMu.RUnlock()

						if ok {
							queueMu.Lock()
							delete(playersInQueue, v.username)
							queueMu.Unlock()
						}
					}
				}
				sessionsMu.Unlock()
			}
		}()
	})

	if _, err := auth(w, r); err == nil {
		w.WriteHeader(http.StatusOK)
	}
}

func LogIn(w http.ResponseWriter, r *http.Request) {
	createdAt := time.Now()
	authInfo, err := getAuthInfoFromRequest(w, r)
	if err != nil {
		w.Write([]byte(err.Error() + "\n"))
	}
	usr := user.User{}
	usr.NickName = authInfo.Nickname
	usr.ID, err = database.GetIdByUsername(authInfo.Nickname)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(key))
	sessionsMu.Lock()
	sessions[key] = sessionValues{
		id:         usr.ID,
		username:   usr.NickName,
		userAgent:  r.UserAgent(),
		lastUpdate: time.Now(),
		createdAt:  createdAt,
	}
	sessionsMu.Unlock()
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

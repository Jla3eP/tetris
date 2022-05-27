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

	findGameCh    = make(chan struct{})
	gameFinder    = sync.Once{}
	currentGameID = int64(0)

	activeGames = make([]gameInfo, 1)
)

func FindGame(w http.ResponseWriter, r *http.Request) {
	if info, err := authUsingSessionKey(w, r); err == nil {
		queueMu.Lock()
		playersInQueue[info.Nickname] = struct{}{}
		queueMu.Unlock()
		findGameCh <- struct{}{}
	}
}

func UpdateSessionTimeout(w http.ResponseWriter, r *http.Request) {
	if _, err := authUsingSessionKey(w, r); err == nil {
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
		w.Write([]byte("nickname currentGameID too short\n"))
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

func init() {
	sessionsCleaner.Do(func() {
		cleanSessions()
	})
	gameFinder.Do(func() {
		tryToFindGame()
	})
}

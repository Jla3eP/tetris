package handling

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jla3eP/tetris/server_side/auth/database"
	"github.com/Jla3eP/tetris/server_side/auth/hash"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const infoToSaltFormat = "%s:%s:%v:%v"

func tryToFindGame() {
	go func() {
		for {
			<-findGameCh
			GameInfo := gameInfo{
				ID:      currentGameID,
				players: make([]string, 0, 2),
			}

			queueMu.Lock()
			for k := range playersInQueue {
				if len(GameInfo.players) < 2 {
					GameInfo.players = append(GameInfo.players, k)
				} else {
					activeGames = append(activeGames, GameInfo)
					for i := range GameInfo.players {
						delete(playersInQueue, GameInfo.players[i])
					}
					break
				}
			}
			queueMu.Unlock()
		}
	}()
}

func cleanSessions() {
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
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
}

func authUsingSessionKey(w http.ResponseWriter, r *http.Request) (*SessionUpdateRequest, error) {
	info, err := getSessionInfoFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, err
	}
	id, err := database.GetIdByUsername(info.Nickname)
	if err != nil {
		return nil, err
	}

	sessionsMu.RLock()
	session, ok := sessions[info.SessionKey]
	sessionsMu.RUnlock()

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, errors.New("can't find session by key")
	}

	if session.username != info.Nickname || session.userAgent != r.UserAgent() || session.id.String() != id.String() {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, errors.New("unauthorized: dismatch detected")
	}

	sessionsMu.Lock()
	session.lastUpdate = time.Now()
	sessions[info.SessionKey] = session
	sessionsMu.Unlock()

	return info, nil
}

func createSessionKey(r *http.Request, username string, ID primitive.ObjectID, createdAt time.Time) (string, error) {
	userAgent := r.UserAgent()
	sessionKey := saltToSessionKey(fmt.Sprintf(infoToSaltFormat, userAgent, username, ID, createdAt))

	key := ""
	for _, v := range sessionKey {
		key += strconv.Itoa(int(v))
	}
	return key, nil
}

func getSessionInfoFromRequest(r *http.Request) (*SessionUpdateRequest, error) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.New("invalid request body")
	}

	info := &SessionUpdateRequest{}
	if err = json.Unmarshal(requestBody, info); err != nil {
		return nil, err
	}

	return info, nil
}

func getAuthInfoFromRequest(w http.ResponseWriter, r *http.Request) (*AuthInfo, error) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error() + "\n"))
		return nil, errors.New("invalid request body")
	}

	usr := &AuthInfo{}

	if err = json.Unmarshal(requestBody, usr); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		log.Print("Invalid authUsingSessionKey json")
		return nil, errors.New("invalid authUsingSessionKey json")
	}

	return usr, nil
}

func saltToSessionKey(info string) []byte {
	return []byte(hash.CreateHash([]byte(info)))
}

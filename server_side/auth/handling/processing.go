package handling

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jla3eP/tetris/both_sides_code"
	"github.com/Jla3eP/tetris/server_side/auth/database"
	"github.com/Jla3eP/tetris/server_side/auth/hash"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	infoToSaltFormat  = "%s:%s:%v:%v"
	figureColorsCount = 5
)

var figuresCount int

func cleanActiveGame(gameID int64, playerIndex int) bool {
	if activeGames[gameID].status == gameStatusEnded {
		activeGames[gameID].playerWatching[playerIndex] = false
		for _, v := range activeGames[gameID].playerWatching {
			if v {
				return true
			}
		}
		delete(activeGames, gameID)
		return true
	}
	return false
}

func checkWatchers(gameID int64) {
	for _, v := range activeGames[gameID].playerWatching {
		if !v {
			return
		}
	}
	AG := activeGames[gameID]
	AG.status = gameStatusEnded
	activeGames[gameID] = AG
}

func appendFigures(gameID int64) {
	gamesMu.RLock()
	gi := generateFigures(activeGames[gameID])
	gamesMu.RUnlock()

	gamesMu.Lock()
	activeGames[gameID] = gi
	gamesMu.Unlock()
}

func needToGenerateNewFigures(gameID int64) bool {
	gamesMu.RLock()
	defer gamesMu.RUnlock()
	for _, v := range activeGames[gameID].playersFiguresIndexes {
		if v > len(activeGames[gameID].figures) {
			return true
		}
	}
	return false
}

func getUserFigureAndColor(gameID int64, figureIndex int) (int, int) {
	gamesMu.RLock()
	defer gamesMu.RUnlock()

	return activeGames[gameID].figures[figureIndex], activeGames[gameID].figuresColors[figureIndex]
}

func getPlayersFigureIndexAndIncrementIt(gameID int64, playerIndex int) int {
	gamesMu.Lock()
	defer gamesMu.Unlock()
	activeGames[gameID].playersFiguresIndexes[playerIndex]++
	return activeGames[gameID].playersFiguresIndexes[playerIndex] - 1
}

func calculateFiguresCount() {
	file, err := os.Open("../both_sides_code/figures_config.json")
	defer file.Close()
	if err != nil {
		log.Fatalln(err)
		return
	}

	buffer, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
		return
	}
	config := both_sides_code.FiguresConfig{}

	err = json.Unmarshal(buffer, &config)
	if err != nil {
		log.Fatalln(err.Error() + " module figure (init)")
		return
	}
	figuresCount = len(config.Figures)
}

func getFieldFromRequest(requestBody []byte) (*both_sides_code.FieldRequest, error) {
	fieldReq := both_sides_code.FieldRequest{}
	fmt.Println(string(requestBody))
	if err := json.Unmarshal(requestBody, &fieldReq); err != nil {
		return nil, err
	}

	return &fieldReq, nil
}

func getFieldInfoToPlayer(gameID int64, playerIndex int) (*both_sides_code.FieldResponse, error) {
	gamesMu.RLock()
	defer gamesMu.RUnlock()
	if _, ok := activeGames[gameID]; !ok {
		return nil, errors.New("can't find game")
	}
	for k, v := range activeGames[gameID].playersLastStatuses {
		if k != playerIndex && v != nil {
			fmt.Println(v)
			response := both_sides_code.FieldResponse{
				FieldRequest: both_sides_code.FieldRequest{
					History: make([]both_sides_code.EnemyFigure, 0, len(v)),
				},
			}
			for _, vv := range v {
				if len(vv.History) > 0 {
					response.History = append(response.History, vv.History[0])
				}
			}
			if playerIndex == 0 {
				activeGames[gameID].playersLastStatuses[1] = nil
			} else {
				activeGames[gameID].playersLastStatuses[0] = nil
			}
			return &response, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("can't find user with index \"%v\" in game with ID \"%v\"", playerIndex, gameID))
}

func clearSecretInfo(field *both_sides_code.FieldRequest) {
	if field != nil {
		field.Nickname = ""
		field.SessionKey = ""
	}
}

func setFieldInfo(gameID int64, playerIndex int, field *both_sides_code.FieldRequest) {
	if field != nil {
		gamesMu.Lock()
		defer gamesMu.Unlock()
		activeGames[gameID].playersLastStatuses[playerIndex] = append(activeGames[gameID].playersLastStatuses[playerIndex], field)
	}
}

func sendResponseWaiting(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	resp := both_sides_code.ResponseStatus{
		Comment: "status: waiting",
	}
	respJSON, _ := json.Marshal(resp)
	w.Write(respJSON)
}

func setGameStatus(gameID int64, status int) {
	gamesMu.Lock()
	defer gamesMu.Unlock()
	currentGame := activeGames[gameID]
	currentGame.status = status
	if status == gameStatusInProgress {
		currentGame = generateFigures(currentGame)
	}
	activeGames[gameID] = currentGame
}

func getGameStatus(gameID int64) int {
	gamesMu.RLock()
	defer gamesMu.RUnlock()
	return activeGames[gameID].status
}

func setPlayerActiveStatus(gameID int64, playerIndex int, status bool) {
	gamesMu.Lock()
	defer gamesMu.Unlock()
	activeGames[gameID].playerActive[playerIndex] = status
}

func arePlayersReady(gameID int64) bool {
	gamesMu.RLock()
	defer gamesMu.RUnlock()
	for _, v := range activeGames[gameID].playerActive {
		if !v {
			return false
		}
	}
	return true
}

func processUnexpectedFindGameResponse(info *both_sides_code.SessionUpdateRequest, w http.ResponseWriter) {
	if _, ok := playersInQueue[info.Nickname]; ok {
		sendResponseWaiting(w)
	} else {
		resp := both_sides_code.ResponseStatus{
			Comment: "try find game ;)",
		}
		respJSON, _ := json.Marshal(resp)

		w.WriteHeader(http.StatusPreconditionRequired)
		w.Write(respJSON)
	}
}

func findGameUsingUsername(username string) (int64, int, bool) {
	gamesMu.RLock()
	defer gamesMu.RUnlock()
	for k, v := range activeGames {
		for playerIndex, name := range v.players {
			if name == username {
				return k, playerIndex, true
			}
		}
	}

	return -1, -1, false
}

func tryToFindGame() {
	go func() {
		for {
			<-findGameCh
			GameInfo := getNewGameInfo()
			queueMu.Lock()

			if len(playersInQueue) >= 2 {
				for k := range playersInQueue {
					if len(GameInfo.players) < 2 {
						GameInfo.players = append(GameInfo.players, k)
					} else {
						break
					}
				}
				gamesMu.Lock()
				activeGames[GameInfo.ID] = GameInfo
				gamesMu.Unlock()

				for i := range GameInfo.players {
					delete(playersInQueue, GameInfo.players[i])
				}
			}

			queueMu.Unlock()
		}
	}()
}

func cleanSessions() {
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
}

func authUsingSessionKey(w http.ResponseWriter, r *http.Request, reqBody []byte) (*both_sides_code.SessionUpdateRequest, error) {
	info, err := getSessionInfoFromRequest(reqBody)
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

func getSessionInfoFromRequest(reqBody []byte) (*both_sides_code.SessionUpdateRequest, error) {
	info := &both_sides_code.SessionUpdateRequest{}
	if err := json.Unmarshal(reqBody, info); err != nil {
		return nil, err
	}

	return info, nil
}

func getAuthInfoFromRequest(w http.ResponseWriter, r *http.Request) (*both_sides_code.AuthInfo, error) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error() + "\n"))
		return nil, errors.New("invalid request body")
	}

	usr := &both_sides_code.AuthInfo{}

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

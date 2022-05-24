package handling

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jla3eP/tetris/server_side/auth/hash"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const infoToSaltFormat = "%s:%s:%v:%v"

func createSessionKey(r *http.Request, username string, ID primitive.ObjectID, createdAt time.Time) (string, error) {
	userAgent := r.UserAgent()
	sessionKey := saltToSessionKey(fmt.Sprintf(infoToSaltFormat, userAgent, username, ID, createdAt))
	return sessionKey, nil
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

		log.Print("Invalid auth json")
		return nil, errors.New("invalid auth json")
	}

	return usr, nil
}

func saltToSessionKey(info string) string {
	return hash.CreateHash([]byte(info))
}

package requests

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jla3eP/tetris/both_sides_code"
	"github.com/Jla3eP/tetris/client_side/field"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client        = http.DefaultClient
	sessionKey    []byte
	StopPingingCh = make(chan struct{})
	pinger        = sync.Once{}
	nickName      string
)

func StopPinging() {
	StopPingingCh <- struct{}{}
}

func LogIn() {
	info, err := logInInfo()
	if err != nil {
		log.Fatalln(err)
		return
	}
	req, _ := http.NewRequest(http.MethodPost, "https://localhost:1234/logIn", info)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		return
	}

	sessionKey, _ = ioutil.ReadAll(resp.Body)
	//pingServer()
}

func Register(info *both_sides_code.AuthInfo) error {
	infoJSON, err := json.Marshal(info)
	buffer := bytes.Buffer{}
	buffer.Write(infoJSON)
	if err != nil {
		return err
	}

	file, err := os.Create("../logInfo.json")
	defer file.Close()
	if err != nil {
		return err
	}

	req, _ := http.NewRequest(http.MethodPost, "https://localhost:1234/register", &buffer)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("unexpected status code %v", resp.StatusCode))
	}

	file.Write(infoJSON)

	return nil
}

func FindGameRequest() error {
	request := bytes.Buffer{}
	request.Write([]byte(fmt.Sprintf(`{"nickname": "%s", "session_key": "%s"}`, nickName, sessionKey)))

	req, _ := http.NewRequest(http.MethodPost, "https://localhost:1234/findGame", &request)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("unexpected code %v", resp.StatusCode))
	}
	return nil
}

func GetGameInfoAndSendMyInfo(reqInfo *both_sides_code.FieldRequest) (*field.Figure, *field.Figure, error) {
	if sessionKey == nil {
		LogIn()
	}

	reqBody := bytes.Buffer{}
	if reqInfo == nil {
		reqBody.WriteString(fmt.Sprintf(`{"nickname": "%s", "session_key": "%s"}`, nickName, sessionKey))
	} else {
		reqInfo.Nickname = nickName
		reqInfo.SessionKey = string(sessionKey)
		JSON, err := json.Marshal(reqInfo)
		if err != nil {
			return nil, nil, err
		}
		reqBody.Write(JSON)
	}

	req, _ := http.NewRequest(http.MethodPost, "https://localhost:1234/getGameInfo", &reqBody)
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBody))

	if isWaitingResponseType(respBody) {
		return nil, nil, errors.New("waiting")
	}

	fieldResp := &both_sides_code.FieldResponse{}
	json.Unmarshal(respBody, fieldResp)

	yourFigure := field.GetFigureUsingIndex(fieldResp.FigureID)
	yourFigure.Color = fieldResp.FigureColor

	if !fieldResp.EnemyFigureSent {
		return yourFigure, nil, nil
	}

	enemyFigure := field.GetFigureUsingIndex(fieldResp.EnemyFigureID)
	enemyFigure.Color = fieldResp.EnemyFigureColor
	enemyFigure.CurrentCoords = fieldResp.EnemyFigureCoords
	enemyFigure.CurrentRotateIndex = fieldResp.EnemyFigureRotateIndex

	return yourFigure, enemyFigure, nil
}

func isWaitingResponseType(respBody []byte) bool {
	rs := &both_sides_code.ResponseStatus{
		Comment: "",
	}
	_ = json.Unmarshal(respBody, rs)
	return rs.Comment == "status: waiting"
}

func pingServer() {
	pinger.Do(func() {
		go func() {
			ticker := time.NewTicker(200 * time.Millisecond)
			requestBody := bytes.Buffer{}
			requestBody.Write([]byte(fmt.Sprintf(`{"nickname": "%s", "session_key": "%s"}`, nickName, sessionKey)))

			req, _ := http.NewRequest(http.MethodPost, "https://localhost:1234/updateSession", &requestBody)

		Loop:
			for {
				<-ticker.C
				select {
				case <-StopPingingCh:
					break Loop
				default:
					break
				}
				resp, err := client.Do(req)
				if err != nil {
					log.Fatalln(err)
					return
				}
				if resp.StatusCode != http.StatusOK {
					log.Fatalln(errors.New(fmt.Sprintf("unexpected status code %v", resp.StatusCode)))
					return
				}
			}
		}()
	})
}

func logInInfo() (*bytes.Buffer, error) {
	file, err := os.Open("./logInfo.json")
	defer file.Close()
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	buffer, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	authInfo := both_sides_code.AuthInfo{}
	err = json.Unmarshal(buffer, &authInfo)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	nickName = authInfo.Nickname

	return bytes.NewBuffer(buffer), nil
}

func init() {
	client.Transport = tr
	sessionKey = nil
}

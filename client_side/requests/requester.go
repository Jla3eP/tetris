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
)

var (
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client     = http.DefaultClient
	sessionKey []byte
	nickName   string
)

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

func GetGameInfoAndSendMyInfo(reqInfo *both_sides_code.FieldRequest) (*field.Figure, []*field.Figure, error) {
	if sessionKey == nil {
		LogIn()
	}

	reqBody := bytes.Buffer{}
	if reqInfo == nil {
		reqBody.WriteString(fmt.Sprintf(`{"nickname": "%s", "session_key": "%s"}`, nickName, sessionKey))
	} else {
		reqInfo.Nickname = nickName
		reqInfo.SessionKey = string(sessionKey)
		reqInfo.History[0].EnemyFigureSent = true
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

	if string(respBody) == "end" {
		return nil, nil, errors.New("end")
	}

	if isWaitingResponseType(respBody) {
		return nil, nil, errors.New("waiting")
	}

	fieldResp := &both_sides_code.FieldResponse{}
	json.Unmarshal(respBody, fieldResp)

	yourFigure := field.GetFigureUsingIndex(fieldResp.FigureID)
	yourFigure.Color = fieldResp.FigureColor

	if len(fieldResp.History) == 0 {
		return yourFigure, nil, nil
	}

	enemyFigures := make([]*field.Figure, 0, 1)
	for _, v := range fieldResp.History {
		enemyFigure := field.GetFigureUsingIndex(v.EnemyFigureID)
		enemyFigure.Color = v.EnemyFigureColor
		enemyFigure.CurrentCoords = v.EnemyFigureCoords
		enemyFigure.CurrentRotateIndex = v.EnemyFigureRotateIndex
		enemyFigures = append(enemyFigures, enemyFigure)
	}

	return yourFigure, enemyFigures, nil
}

func isWaitingResponseType(respBody []byte) bool {
	rs := &both_sides_code.ResponseStatus{
		Comment: "",
	}
	_ = json.Unmarshal(respBody, rs)
	return rs.Comment == "status: waiting"
}

func logInInfo() (*bytes.Buffer, error) {
	authInfo := both_sides_code.AuthInfo{
		Nickname: os.Args[1],
		Password: os.Args[2],
	}
	buffer, err := json.Marshal(authInfo)
	if err != nil {
		return nil, err
	}
	nickName = authInfo.Nickname

	return bytes.NewBuffer(buffer), nil
}

func init() {
	client.Transport = tr
	sessionKey = nil
}

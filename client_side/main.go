package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
)

//func main() {
//	mg := mainGame.NewGame()
//	et.SetWindowSize(render.GetX(), render.GetY())
//	et.SetWindowTitle("TetriX")
//	if err := et.RunGame(mg); err != nil {
//		log.Fatal(err)
//	}
//}

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.DefaultClient
	client.Transport = tr

	request := bytes.Buffer{}
	request.Write([]byte(`{"nickname": "Jla3eP", "password": "1234qwer"}`))

	req, err := http.NewRequest(http.MethodPost, "https://localhost:1234/logIn", &request)
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("err != nil", err)
	} else {
		respBody, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("%v %s\n", resp.StatusCode, string(respBody))

		request := bytes.Buffer{}
		body := fmt.Sprintf(`{"nickname": "Jla3eP", "session_key": "%s"}`, respBody)
		fmt.Println(body)
		request.Write([]byte(body))

		req, err = http.NewRequest(http.MethodPost, "https://localhost:1234/updateSession", &request)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("err != nil", err)
		} else {
			respBody, _ := ioutil.ReadAll(resp.Body)
			fmt.Printf("%v %s\n", resp.StatusCode, string(respBody))
		}
	}
}

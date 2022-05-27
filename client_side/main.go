package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
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

	req, _ := http.NewRequest(http.MethodPost, "https://localhost:1234/logIn", &request)
	resp, _ := client.Do(req)

	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("%v %s\n", resp.StatusCode, string(respBody))

	{
		time.Sleep(11 * time.Second)

		request := bytes.Buffer{}
		request.Write([]byte(fmt.Sprintf(`{"nickname": "Jla3eP", "session_key": "%s"}`, respBody)))

		req, err := http.NewRequest(http.MethodPost, "https://localhost:1234/findGame", &request)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("err != nil", err)
		} else {
			respBody, _ := ioutil.ReadAll(resp.Body)
			fmt.Printf("%v %s\n", resp.StatusCode, string(respBody))
		}
		resp, err = client.Do(req)
		if err != nil {
			fmt.Println("err != nil", err)
		} else {
			respBody, _ := ioutil.ReadAll(resp.Body)
			fmt.Printf("%v %s\n", resp.StatusCode, string(respBody))
		}
	}
}

package main

import (
	"github.com/Jla3eP/tetris/client_side/mainGame"
	"github.com/Jla3eP/tetris/client_side/render"
	et "github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	mg := mainGame.NewGame()
	et.SetWindowSize(render.GetX(), render.GetY())
	et.SetWindowTitle("TetriX")
	if err := et.RunGame(mg); err != nil {
		log.Fatal(err)
	}
}

//func main() {
//	mainGame.StartGame()
//	/*tr := &http.Transport{
//		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
//	}
//	client := http.DefaultClient
//	client.Transport = tr
//
//	request := bytes.Buffer{}
//	request.Write([]byte(`{"nickname": "Jla3eP", "password": "1234qwer"}`))
//
//	req, err := http.NewRequest("POST", "https://localhost:1234/register", &request)
//	resp, err := client.Do(req)
//
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		respBody, _ := ioutil.ReadAll(resp.Body)
//		fmt.Printf("%v %s", resp.StatusCode, string(respBody))
//	}*/
//}

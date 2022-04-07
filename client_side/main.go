package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.DefaultClient
	client.Transport = tr

	request := bytes.Buffer{}
	request.Write([]byte(`{"nickname": "Jla3eP", "password": "1234qwer"}`))

	req, err := http.NewRequest("POST", "https://localhost:1234/register", &request)
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	} else {
		respBody, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("%v %s", resp.StatusCode, string(respBody))
	}
}

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func getHostname(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return string(body), err
}

func main() {
	url := "https://poc-app-2-cloudhub-eu-west-1.dev.hydra.sophos.com/hostname"

	resp, err := getHostname(url)

	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("Response:", resp)
}

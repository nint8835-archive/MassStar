package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// BotToken represents an individual token to be used
type BotToken struct {
	Name  string `json:"name,omitempty"`
	Token string `json:"token,omitempty"`
}

func main() {
	var tokens []BotToken

	tokenFile, _ := os.Open("tokens.json")
	defer tokenFile.Close()
	tokenBytes, _ := ioutil.ReadAll(tokenFile)

	json.Unmarshal(tokenBytes, &tokens)

	fmt.Println(tokens)
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
)

// Config is struct holding token and gist ID which are used to sync the config files from the system.
type Config struct {
	GithubToken   string
	GistID        string
	EncryptionKey string
}

var config = &Config{
	GistID:        "",
	GithubToken:   "",
	EncryptionKey: "",
}

func generateNewConfig() {
	fmt.Print("Enter Gist ID: ")
	fmt.Scanln(&config.GistID)
	fmt.Print("Enter GitHub token: ")
	fmt.Scanln(&config.GithubToken)
	fmt.Print("Enter your encryption key (needs to have 32 characters): ")
	fmt.Scanln(&config.EncryptionKey)
	fmt.Println("You can change the settings later in " + path.Join(getHomeDir(), CONFIGFILE))

	syncGist := getSync()

	fmt.Println(".. reading content of the gist")
	syncGist.ReadRemote()

	fmt.Println(".. saving list of files into the config along the ID and the token")
	content, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path.Join(getHomeDir(), CONFIGFILE), content, 0600)
	if err != nil {
		panic(err)
	}
}

func loadConfiguration() {
	configFilepath := path.Join(getHomeDir(), CONFIGFILE)

	content, err := ioutil.ReadFile(configFilepath)
	if err != nil {
		generateNewConfig()
		return
	}

	err = json.Unmarshal(content, &config)
	if err != nil {
		generateNewConfig()
		return
	}

	fmt.Println(config)
}

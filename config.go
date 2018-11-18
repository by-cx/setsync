package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
)

// Config is struct holding token and gist ID which are used to sync the config files from the system.
type Config struct {
	Mode          string // possible choices: gist, s3
	GithubToken   string // GitHub token, only for gist mode
	GistID        string // Gist ID, only for gist mode
	EncryptionKey string // Encryption key to protect the data
	S3AccessKey   string // only for s3 mode
	S3SecretKey   string // only for s3 mode
	S3Bucket      string // only for s3 mode
	S3Endpoint    string // only for s3 mode
}

var config = &Config{
	GistID:        "",
	GithubToken:   "",
	EncryptionKey: "",
}

func generateNewConfig() {
	fmt.Print("Enter the mode (s3, gist): ")
	fmt.Scanln(&config.Mode)

	if config.Mode == "s3" {
		fmt.Print("Enter access key: ")
		fmt.Scanln(&config.S3AccessKey)
		fmt.Print("Enter secret key: ")
		fmt.Scanln(&config.S3SecretKey)
		fmt.Print("Enter bucket: ")
		fmt.Scanln(&config.S3Bucket)
		fmt.Print("Enter endpoint: ")
		fmt.Scanln(&config.S3Endpoint)
	}

	if config.Mode == "gist" {
		fmt.Print("Enter Gist ID: ")
		fmt.Scanln(&config.GistID)
		fmt.Print("Enter GitHub token: ")
		fmt.Scanln(&config.GithubToken)
		fmt.Print("Enter your encryption key (needs to have 32 characters): ")
		fmt.Scanln(&config.EncryptionKey)
		fmt.Println("You can change the settings later in " + path.Join(getHomeDir(), CONFIGFILE))
	}

	syncGist := getSync()

	fmt.Println(".. reading content of the " + config.Mode)
	syncGist.ReadRemote()

	fmt.Println(".. saving list of files into the config file")
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
}

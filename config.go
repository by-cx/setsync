package main

// Config is struct holding token and gist ID which are used to sync the config files from the system.
type Config struct {
	GithubToken string
	GistID      string
}

var config = &Config{
	GistID:      "",
	GithubToken: "",
}

package main

// Config is struct holding token and gist ID which are used to sync the config files from the system.
type Config struct {
	GithubToken string
	GistID      string
}

var config = &Config{
	GistID:      "",
	GithubToken: "",
	///home/cx/.config/terminator/config /home/cx/.ssh/config /home/cx/.config/fish/config.fish /home/cx/.config/fish/functions/elxvpn.fish /home/cx/.vimrc /home/cx/.gitconfig /home/cx/.config/Slack/storage/slack-teams
}

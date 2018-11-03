package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/cosiner/flag"
)

// CONFIGFILE contains path to the config file
const CONFIGFILE = ".setsync"

// Cmd line arguments
type Cmd struct {
	Upload struct {
		Enable bool
	} `usage:"uploads config files into the gist"`
	Download struct {
		Enable bool
	} `usage:"downloads config files from the gist"`
	Add struct {
		Enable bool
		Files  []string `args:"true"`
	} `usage:"adds file into the gist"`
	Remove struct {
		Enable bool
		Files  []string `args:"true"`
	} `usage:"removes file from the gist"`
	List struct {
		Enable bool
	} `usage:"list files saved in the gist"`
	Help struct {
		Enable bool
	} `usage:"show help"`
}

func generateNewConfig() {
	fmt.Print("Enter Gist ID: ")
	fmt.Scanln(&config.GistID)
	fmt.Print("Enter GitHub token: ")
	fmt.Scanln(&config.GithubToken)

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	syncGist := SyncGist{
		GistID: config.GistID,
	}

	fmt.Println(".. reading content of the gist")
	syncGist.ReadRemote()

	fmt.Println(".. saving list of files into the config along the ID and the token")
	content, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path.Join(usr.HomeDir, CONFIGFILE), content, 0600)
	if err != nil {
		panic(err)
	}
}

func loadConfiguration() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	configFilepath := path.Join(usr.HomeDir, CONFIGFILE)

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

func main() {
	loadConfiguration()

	var cmd Cmd

	set := flag.NewFlagSet(flag.Flag{})
	set.StructFlags(&cmd)
	set.Parse()
	//set.Help(false)
	//fmt.Println()

	syncGist := SyncGist{
		GistID: config.GistID,
		Files:  make(map[string]*SyncFile),
	}

	if cmd.Upload.Enable {
		fmt.Println("Start uploading the config files")
		fmt.Println(".. reading the gist")
		syncGist.ReadRemote()
		fmt.Println(".. reading local content of the files")
		syncGist.ReadLocal()
		fmt.Println(".. writing the content to the gist")
		syncGist.WriteRemote()
	} else if cmd.Download.Enable {
		fmt.Println("Start downloading the config files")
		fmt.Println(".. reading the gist")
		syncGist.ReadRemote()
		fmt.Println(".. writing content of the config files")
		syncGist.WriteLocal()
	} else if cmd.Add.Enable {
		if len(cmd.Add.Files) == 0 {
			fmt.Println("No files")
		}

		fmt.Println("New files")

		usr, err := user.Current()
		if err != nil {
			panic(err)
		}

		for _, file := range cmd.Add.Files {
			fullPath, err := filepath.Abs(file)
			if err != nil {
				panic(err)
			}

			if strings.HasPrefix(fullPath, usr.HomeDir) {
				finalPath := strings.Trim(fullPath[len(usr.HomeDir):], "/")

				fmt.Println(".. adding " + finalPath)
				syncGist.Files[finalPath] = &SyncFile{
					Filename: &finalPath,
				}
			} else {
				fmt.Println(".." + file + " is not located in your home directory")
			}
		}

		fmt.Println(".. reading local content of files")
		syncGist.ReadLocal()
		fmt.Println(".. writing the content to the gist")
		syncGist.WriteRemote()
	} else if cmd.Remove.Enable {
		fmt.Println("Removing files from the gist")
		fmt.Println(".. reading the content of the gist")
		syncGist.ReadRemote()

		for _, file := range cmd.Remove.Files {
			if _, ok := syncGist.Files[file]; ok {
				fmt.Println(".. removing " + file)
				syncGist.Files[file] = &SyncFile{
					Filename: nil,
				}
			} else {
				fmt.Println(".. file " + file + " not found")
			}
		}
		fmt.Println(".. writing the content to the gist")
		syncGist.WriteRemote()

	} else if cmd.Help.Enable {
		set.Help(true)
	} else {
		fmt.Println(".. reading the gist")
		syncGist.ReadRemote()
		fmt.Println("")

		for _, file := range syncGist.Files {
			fmt.Println("* " + *file.Filename)
		}
	}
}

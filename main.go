package main

import (
	"errors"
	"fmt"
	"os/user"
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

func getHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	return usr.HomeDir
}

func getSync() *Sync {
	var backend BackendInterface

	if config.Mode == "gist" {
		backend = &GistBackend{
			GistID:      config.GistID,
			GitHubToken: config.GithubToken,
		}
	} else if config.Mode == "s3" {
		backend = &S3Backend{
			Endpoint:  config.S3Endpoint,
			AccessKey: config.S3AccessKey,
			SecretKey: config.S3SecretKey,
			Bucket:    config.S3Bucket,
		}
	} else {
		panic(errors.New("Error: unknown mode"))
	}

	sync := &Sync{
		Backend: backend,
		Files:   make(map[string]*SyncFile),
	}

	return sync
}

func main() {
	loadConfiguration()

	var cmd Cmd

	set := flag.NewFlagSet(flag.Flag{})
	set.StructFlags(&cmd)
	set.Parse()

	syncGist := getSync()

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

		for _, file := range cmd.Add.Files {
			fullPath, err := filepath.Abs(file)
			if err != nil {
				panic(err)
			}

			if strings.HasPrefix(fullPath, getHomeDir()) {
				finalPath := strings.Trim(fullPath[len(getHomeDir()):], "/")

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

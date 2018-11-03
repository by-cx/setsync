package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/levigross/grequests"
)

const GITHUB_ENDPOINT = "https://api.github.com"

type GistResponse struct {
	URL   string `json:"url"`
	Files map[string]struct {
		Filename string `json:"filename"`
		Content  string `json:"content"`
		Size     int    `json:"size"`
		Type     string `json:"type"`
		Language string `json:"language"`
	} `json:"files"`
}

type GistWriteRequestFile struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

type GistWriteRequest struct {
	Description string                          `json:"description"`
	Files       map[string]GistWriteRequestFile `json:"files"`
}

type SyncFile struct {
	Filename *string
	Content  string
}

// SyncGist represents all files we synchronize between this computer and Gist
type SyncGist struct {
	GistID string
	Files  map[string]*SyncFile
}

func (s *SyncGist) getPrefix() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	return usr.HomeDir
}

// UpToDate checks if local version is up to date. True means loaded version is in sync. False means it's not.
func (s *SyncGist) UpToDate() bool {
	return true
}

// ReadRemote reads content of the files from the gist.
func (s *SyncGist) ReadRemote() {
	s.Files = make(map[string]*SyncFile)

	resp, err := grequests.Get(
		GITHUB_ENDPOINT+"/gists/"+s.GistID, &grequests.RequestOptions{
			Headers: map[string]string{
				"Content-type":  "application/json",
				"Authorization": "token " + config.GithubToken,
			},
		})
	if err != nil {
		panic(err)
	}

	gist := GistResponse{}

	err = json.Unmarshal(resp.Bytes(), &gist)
	if err != nil {
		panic(err)
	}

	s.Files = make(map[string]*SyncFile)

	for _, file := range gist.Files {
		filename := strings.Replace(file.Filename, "__", "/", -1)
		s.Files[file.Filename] = &SyncFile{
			Filename: &filename,
			Content:  file.Content,
		}
	}
}

// WriteLocal writes content of s.Files into the gist.
func (s *SyncGist) WriteLocal() {
	for _, file := range s.Files {
		// Ignore error here
		os.MkdirAll(path.Join(s.getPrefix(), path.Dir(*file.Filename)), 0755)

		err := ioutil.WriteFile(path.Join(s.getPrefix(), *file.Filename), []byte(file.Content), 0644)
		if err != nil {
			panic(err)
		}
	}
}

// WriteRemote writes content of s.Files into the gist.
func (s *SyncGist) WriteRemote() {
	var requestData = &GistWriteRequest{
		Description: "test",
		Files:       make(map[string]GistWriteRequestFile),
	}
	for key, file := range s.Files {
		requestData.Files[strings.Replace(key, "/", "__", -1)] = GistWriteRequestFile{
			Filename: strings.Replace(key, "/", "__", -1),
			Content:  file.Content,
		}
	}

	_, err := grequests.Patch(GITHUB_ENDPOINT+"/gists/"+s.GistID, &grequests.RequestOptions{
		Headers: map[string]string{
			"Content-type":  "application/json",
			"Authorization": "token " + config.GithubToken,
		},
		JSON: requestData,
	})

	if err != nil {
		panic(err)
	}
}

// ReadLocal loads content of the files from the local files
func (s *SyncGist) ReadLocal() {
	for _, file := range s.Files {
		content, err := ioutil.ReadFile(path.Join(s.getPrefix(), *file.Filename))
		if err != nil {
			panic(err)
		}

		s.Files[*file.Filename] = &SyncFile{
			Filename: file.Filename,
			Content:  string(content),
		}
	}
}

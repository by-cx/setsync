package main

import (
	"encoding/json"
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

type GistBackend struct {
	GistID      string
	GitHubToken string
}

// Read content of the files from the gist
func (g *GistBackend) Read() *map[string]*SyncFile {
	resp, err := grequests.Get(
		GITHUB_ENDPOINT+"/gists/"+g.GistID, &grequests.RequestOptions{
			Headers: map[string]string{
				"Content-type":  "application/json",
				"Authorization": "token " + g.GitHubToken,
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

	files := make(map[string]*SyncFile)

	for _, file := range gist.Files {
		filename := strings.Replace(file.Filename, "__", "/", -1)
		files[file.Filename] = &SyncFile{
			Filename: &filename,
			Content:  file.Content,
		}
	}

	return &files
}

// Write content of the files from the gist
func (g *GistBackend) Write(files *map[string]*SyncFile) {
	var requestData = &GistWriteRequest{
		Description: "test",
		Files:       make(map[string]GistWriteRequestFile),
	}
	for key, file := range *files {
		requestData.Files[strings.Replace(key, "/", "__", -1)] = GistWriteRequestFile{
			Filename: strings.Replace(key, "/", "__", -1),
			Content:  file.Content,
		}
	}

	_, err := grequests.Patch(GITHUB_ENDPOINT+"/gists/"+g.GistID, &grequests.RequestOptions{
		Headers: map[string]string{
			"Content-type":  "application/json",
			"Authorization": "token " + g.GitHubToken,
		},
		JSON: requestData,
	})

	if err != nil {
		panic(err)
	}
}

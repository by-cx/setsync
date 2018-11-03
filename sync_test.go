package main

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

var lastFilename string
var lastContent string

const testContent = "test"
const testFilename = ".test_config"

type DummyBackend struct {
}

func (d *DummyBackend) Read() *map[string]*SyncFile {
	files := make(map[string]*SyncFile)
	filename := testFilename
	files[testFilename] = &SyncFile{
		Filename: &filename,
		Content:  lastContent,
	}

	return &files
}

func (d *DummyBackend) Write(files *map[string]*SyncFile) {
	for key, file := range *files {
		lastFilename = key
		lastContent = file.Content
	}
}

func TestSync(t *testing.T) {
	sync := Sync{
		Backend: &DummyBackend{},
	}
	sync.ReadRemote()
	assert.Equal(t, testFilename, *sync.Files[testFilename].Filename)

	sync.WriteLocal()

	content, err := ioutil.ReadFile(path.Join(getHomeDir(), testFilename))

	assert.Nil(t, err)
	assert.Equal(t, lastContent, string(content))

	sync.WriteRemote()

	assert.Equal(t, testFilename, lastFilename)
	assert.Equal(t, lastContent, lastContent)
}

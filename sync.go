package main

import (
	"io/ioutil"
	"os"
	"path"
)

// BackendInterface is interface for all backend structs
type BackendInterface interface {
	Write(*map[string]*SyncFile)
	Read() *map[string]*SyncFile
}

// SyncFile represents file in the Sync structure
type SyncFile struct {
	Filename *string
	Content  string
}

// Sync represents all files we synchronize between this computer and Gist
type Sync struct {
	Files   map[string]*SyncFile
	Backend BackendInterface
}

func (s *Sync) getPrefix() string {
	return getHomeDir()
}

// UpToDate checks if local version is up to date. True means loaded version is in sync. False means it's not.
func (s *Sync) UpToDate() bool {
	return true
}

// ReadRemote reads content of the files from the gist.
func (s *Sync) ReadRemote() {
	s.Files = *s.Backend.Read()

	for key, file := range s.Files {
		if len(file.Content) <= 12 {
			s.Files[key].Content = ""
		} else {
			decryptedContent, err := decrypt(file.Content)
			if err != nil {
				panic(err)
			}

			s.Files[key].Content = string(decryptedContent)
		}
	}
}

// WriteLocal writes content of s.Files into the gist.
func (s *Sync) WriteLocal() {
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
func (s *Sync) WriteRemote() {
	for key, file := range s.Files {
		encryptedContent, err := encrypt([]byte(file.Content))
		if err != nil {
			panic(err)
		}
		s.Files[key].Content = encryptedContent
	}

	s.Backend.Write(&s.Files)
}

// ReadLocal loads content of the files from the local files
func (s *Sync) ReadLocal() {
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

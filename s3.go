package main

import (
	"bytes"
	"strings"

	minio "github.com/minio/minio-go"
)

// S3Backend is backend for S3 storage of the data
type S3Backend struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Endpoint  string
}

// getS3Connection returns configured S3 client
func (s *S3Backend) getS3Connection() *minio.Client {
	s3Client, err := minio.New(s.Endpoint, s.AccessKey, s.SecretKey, true)
	if err != nil {
		panic(err)
	}

	return s3Client
}

// Read content of the files from the gist
func (s *S3Backend) Read() *map[string]*SyncFile {
	files := make(map[string]*SyncFile)

	client := s.getS3Connection()

	doneCh := make(chan struct{})
	defer close(doneCh)

	for object := range client.ListObjectsV2(s.Bucket, "", false, doneCh) {
		if object.Err != nil {
			panic(object.Err)
		}
		key := strings.Replace(object.Key, "__", "/", -1)

		reader, err := client.GetObject(s.Bucket, object.Key, minio.GetObjectOptions{})
		if err != nil {
			panic(err)
		}

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(reader)
		if err != nil {
			panic(err)
		}

		files[key] = &SyncFile{
			Filename: &key,
			Content:  buf.String(),
		}
	}

	return &files
}

// Write content of the files from the gist
func (s *S3Backend) Write(files *map[string]*SyncFile) {
	client := s.getS3Connection()

	var key string
	var reader *strings.Reader

	for filename, file := range *files {
		key = strings.Replace(filename, "/", "__", -1)
		reader = strings.NewReader(file.Content)

		_, err := client.PutObject(
			s.Bucket,
			key,
			reader,
			int64(reader.Len()),
			minio.PutObjectOptions{},
		)
		if err != nil {
			panic(err)
		}
	}
}

package client

import (
	"io"
	"net/http"

	"google.golang.org/api/drive/v2"
)

type RemoteFileClient interface {
	ConfigureRemoteFileDestination() error
	List() ([]*File, error)
	Stream(string, int64, int64) (io.ReadCloser, error)
}

type DriveClient struct {
	client       *http.Client
	driveService *drive.Service
}

type File struct {
	Name string
	Id   string
	Size int64
}

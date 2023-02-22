package models

import (
	"time"
)

type File struct {
	FileId     string    `json:"fileid"`
	Name       string    `json:"name"`
	Size       int64     `json:"size"`
	Created_At time.Time `json:"created_at"`
}

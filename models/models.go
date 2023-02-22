package models

import (
	"time"
)

type File struct {
	FileId     string    `json:"fileid" gorm:"string:not null:default:null"`
	Name       string    `json:"name" gorm:"string:not null:default:null"`
	Size       int64     `json:"size" gorm:"int:not null:default:null"`
	Created_At time.Time `json:"created_at"`
}

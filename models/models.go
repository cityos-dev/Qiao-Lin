package models

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	FileId     uuid.UUID `json:"fileid" gorm:"string:not null:default:null"`
	Name       string    `json:"name" gorm:"string:not null:default:null"`
	Size       int64     `json:"size" gorm:"int:not null:default:null"`
	Created_At time.Time `json:"created_at"`
}

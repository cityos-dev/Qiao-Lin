package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	FileId uuid.UUID `json:"fileid" gorm:"string:not null:default:null"`
	Name   string    `json:"name" gorm:"string:not null:default:null"`
	Size   int64     `json:"size" gorm:"int:not null:default:null"`
	// Content []byte    `json:"content" gorm:"byte:not null:default:null"`
}

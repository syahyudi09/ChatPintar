package model

import (
	"time"
)

type FileModel struct {
	FileID	string     
	FileName  string   
	FilePath  string   
	Uploader  string    
	CreatedAt time.Time
}
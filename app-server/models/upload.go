package models

import "time"

type File struct {
	ID         int64     `json:"id"`
	FileName   string    `json:"file_name"`
	FileType   string    `json:"file_type"`
	FileSize   int64     `json:"file_size"`
	FileHash   string    `json:"file_hash"`
	CreateTime time.Time `json:"create_time"`
}

package models

type File struct {
	ID         int64  `json:"id"`
	FileName   string `json:"file_name"`
	FileType   string `json:"file_type"`
	FileSize   int64  `json:"file_size"`
	FileHash   string `json:"file_hash"`
	CreateTime string `json:"create_time"`
}

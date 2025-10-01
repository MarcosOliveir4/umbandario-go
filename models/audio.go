package models

import "time"

type AudioFile struct {
	ID        string    `json:"id"`
	Filename  string    `json:"filename"`
	Filetype  string    `json:"filetype"`
	Path      string    `json:"path"`
	LineID    string    `json:"line_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AudioFileList struct {
	Dados []AudioFile `json:"dados"`
}

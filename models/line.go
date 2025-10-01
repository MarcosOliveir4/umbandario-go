package models

import "time"

type Line struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LineList struct {
	Dados []Line `json:"dados"`
}

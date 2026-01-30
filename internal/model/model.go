package model

import "time"

type NewNewsRecord struct {
	Author    string    `json:"author,omitempty"`
	Title     string    `json:"title,omitempty"`
	Summary   string    `json:"summary,omitempty"`
	Content   string    `json:"content,omitempty"`
	Source    string    `json:"source,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	Tags      []string  `json:"tags,omitempty"`
}

package model

import (
	"github.com/logeshwarann-dev/news-api-rest/internal/news"
)

type NewsRecord struct {
	Author    string   `json:"author,omitempty"`
	Title     string   `json:"title,omitempty"`
	Summary   string   `json:"summary,omitempty"`
	Content   string   `json:"content,omitempty"`
	Source    string   `json:"source,omitempty"`
	CreatedAt string   `json:"created_at,omitempty"`
	Tags      []string `json:"tags,omitempty"`
}

type AllNewsRecords struct {
	NewsRecords []news.Record `json:"news,omitempty"`
}

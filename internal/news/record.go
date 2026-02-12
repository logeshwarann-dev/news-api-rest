package news

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Record struct {
	bun.BaseModel `bun:"table:news"`
	Id            uuid.UUID `json:"id" bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	Author        string    `json:"author" bun:"author,nullzero,notnull"`
	Title         string    `json:"title" bun:"title,nullzero,notnull"`
	Summary       string    `json:"summary" bun:"summary,nullzero,notnull"`
	Content       string    `json:"content" bun:"content,nullzero,notnull"`
	Source        string    `json:"source" bun:"source,nullzero,notnull"`
	Tags          []string  `json:"tags" bun:"tags,nullzero,notnull,array"`
	CreatedAt     time.Time `json:"created_at" bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `json:"updated_at" bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	DeletedAt     time.Time `json:"deleted_at" bun:"deleted_at,nullzero,soft_delete"`
}

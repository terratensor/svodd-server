package answer

import (
	"context"
	"time"
)

type Entry struct {
	ID         *int64     `json:"id"`
	Username   string     `json:"username"`
	Text       string     `json:"text"`
	AvatarFile string     `json:"avatar_file"`
	Url        string     `json:"url"`
	Role       string     `json:"role"`
	Datetime   *time.Time `json:"datetime"`
	DataID     int64      `json:"data_id,omitempty"`
	ParentID   int64      `json:"parent_id,omitempty"`
	Type       int        `json:"type"`
	Position   int        `json:"position"`
}

type StorageInterface interface {
	// FindByUrl(ctx context.Context, url string) (*Entry, error)
	FindAllByUrl(ctx context.Context, url string) (*[]Entry, error)
	Insert(ctx context.Context, entry *Entry) (*int64, error)
	// Update(ctx context.Context, entry *Entry) error
}

type Entries struct {
	Storage StorageInterface
}

func NewAnswerStorage(store StorageInterface) *Entries {
	return &Entries{
		Storage: store,
	}
}

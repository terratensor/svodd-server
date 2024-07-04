package answer

import (
	"time"
)

type Entry struct {
	ID         *int64     `json:"id"`
	Url        string     `json:"url"`
	Username   string     `json:"username"`
	Text       string     `json:"text"`
	AvatarFile string     `json:"avatar_file"`
	Role       string     `json:"role"`
	Datetime   *time.Time `json:"datetime"`
	DataID     string     `json:"data_id,omitempty"`
	ParentID   string     `json:"parent_id,omitempty"`
	Type       string     `json:"type"`
	Position   int        `json:"position"`
	Chunk      int        `json:"chunk"`
}

type StorageInterface interface {
	// FindByUrl(ctx context.Context, url string) (*Entry, error)
	// Insert(ctx context.Context, entry *Entry) (*int64, error)
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

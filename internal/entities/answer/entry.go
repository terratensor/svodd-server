package answer

import (
	"context"
	"log/slog"
	"time"

	"github.com/terratensor/svodd-server/internal/lib/logger/sl"
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
	FindAllByUrl(ctx context.Context, url string) (*[]Entry, error)
	Insert(ctx context.Context, entry *Entry) (*int64, error)
	Update(ctx context.Context, entry *Entry) error
}

type Entries struct {
	Storage StorageInterface
}

func NewAnswerStorage(store StorageInterface) *Entries {
	return &Entries{
		Storage: store,
	}
}

func (es *Entries) Insert(e *Entry, logger *slog.Logger) error {
	id, err := es.Storage.Insert(context.Background(), e)
	if err != nil {
		logger.Error(
			"failed insert entry",
			slog.String("url", e.Url),
			sl.Err(err),
		)
		return err
	}
	logger.Info(
		"entry successful inserted",
		slog.Int64("id", *id),
		slog.String("url", e.Url),
	)
	return nil
}

func (es *Entries) Update(e *Entry, logger *slog.Logger) error {
	err := es.Storage.Update(context.Background(), e)
	if err != nil {
		logger.Error(
			"failed update entry",
			slog.String("url", e.Url),
			sl.Err(err),
		)
		return err
	}
	logger.Info(
		"entry successful updated",
		slog.Int64("id", *e.ID),
		slog.String("url", e.Url),
	)
	return nil
}

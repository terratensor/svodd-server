package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/terratensor/svodd-server/internal/config"
	"github.com/terratensor/svodd-server/internal/entities/answer"
	"github.com/terratensor/svodd-server/internal/qaparser/questionanswer"
	"github.com/terratensor/svodd-server/internal/storage/manticore"
)

type App struct {
	manticoreStorages *[]answer.Entries
}

func NewApp(cfg *config.Config) *App {

	var manticoreStorages []answer.Entries
	// Создаем срез клиетнов мантикоры по количеству индексов в конфиге
	for _, index := range cfg.ManticoreIndex {
		manticoreStorages = append(manticoreStorages, *NewEntriesStorage(index.Name))
	}

	return &App{
		manticoreStorages: &manticoreStorages,
	}
}
func NewEntriesStorage(index string) *answer.Entries {
	var storage answer.StorageInterface

	manticoreClient, err := manticore.New(index)
	if err != nil {
		log.Printf("failed to initialize manticore client for index %v, %v", index, err)
		os.Exit(1)
	}

	storage = manticoreClient

	return answer.NewAnswerStorage(storage)
}

func (app *App) Process(entry *questionanswer.Entry) error {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	answerEntries := *makeAnswerEntries(entry)

	for _, storage := range *app.manticoreStorages {

		store := storage

		dbe, err := store.Storage.FindAllByUrl(context.Background(), entry.Url.String())
		if err != nil {
			log.Fatalf("failed to find by url: %v", err)
		}

		if dbe == nil || len(*dbe) == 0 {

			for _, e := range answerEntries {
				err = store.Insert(&e, logger)
				if err != nil {
					log.Fatalf("failed to insert new entry: %v", err)
				}
			}
		} else {
			log.Printf("entry already exists")
			for n, e := range answerEntries {
				if n < len(*dbe) {
					e.ID = (*dbe)[n].ID
					err = store.Update(&e, logger)
					if err != nil {
						return err
					}
				} else {
					err = store.Insert(&e, logger)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

const TypeAQTeaser = 4
const TypeAQFragmnt = 5
const TypeAQComment = 3

// MakeAnswerEntries generates a slice of answer.Entry based on the given questionanswer.Entry.
//
// Parameters:
// - entry: a pointer to a questionanswer.Entry struct representing the entry to generate answer entries from.
//
// Return:
// - a pointer to a slice of answer.Entry representing the generated answer entries.
func makeAnswerEntries(entry *questionanswer.Entry) *[]answer.Entry {
	var entries []answer.Entry
	position := 1

	text := fmt.Sprintf("<h4>%v</h4> <p><span class=\"link\">%v</span></p>", entry.Title, entry.Video.String())
	if len(entry.Fragments) == 0 {
		text += entry.Html
	}

	answerEntry := answer.Entry{
		Username: entry.Title,
		Text:     text,
		Url:      entry.Url.String(),
		Datetime: entry.Datetime,
		Type:     TypeAQTeaser,
		Position: position,
	}
	position++
	entries = append(entries, answerEntry)

	for _, fragm := range entry.Fragments {
		answerEntry := answer.Entry{
			Username: entry.Title,
			Text:     fragm.QuestionAnswer,
			Url:      entry.Url.String(),
			Datetime: entry.Datetime,
			Type:     TypeAQFragmnt,
			Position: position,
		}
		position++
		entries = append(entries, answerEntry)
	}

	for _, comment := range entry.Comments {
		DataID, _ := strconv.ParseInt(comment.DataID, 10, 64)
		answerEntry := answer.Entry{
			Username:   comment.Username,
			Text:       comment.Text,
			Url:        entry.Url.String(),
			AvatarFile: comment.AvatarFile.String(),
			Role:       comment.Role,
			Datetime:   comment.Datetime,
			DataID:     DataID,
			Type:       TypeAQComment,
			Position:   position,
		}
		position++
		entries = append(entries, answerEntry)
	}

	return &entries
}

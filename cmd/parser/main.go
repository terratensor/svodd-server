package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/terratensor/svodd-server/internal/app"
	"github.com/terratensor/svodd-server/internal/config"
	"github.com/terratensor/svodd-server/internal/entities/answer"
	"github.com/terratensor/svodd-server/internal/lib/httpclient"
	"github.com/terratensor/svodd-server/internal/qaparser/qavideo"
	"github.com/terratensor/svodd-server/internal/qaparser/questionanswer"
	"github.com/terratensor/svodd-server/internal/workerpool"
)

func main() {
	cfg := config.MustLoad()

	ch := make(chan *url.URL, cfg.EntryChanBuffer)

	wg := &sync.WaitGroup{}
	for _, parserCfg := range cfg.Parsers {

		wg.Add(1)
		parser := qavideo.NewParser(parserCfg, *cfg.Delay, *cfg.RandomDelay)
		go parser.Run(ch, wg)
	}

	var allTask []*workerpool.Task

	// Создаем срез клиетнов мантикоры по количеству индексов в конфиге
	var manticoreStorages []answer.Entries
	for _, index := range cfg.ManticoreIndex {
		manticoreStorages = append(manticoreStorages, *app.NewEntriesStorage(index.Name))
	}

	for page := range ch {
		log.Printf("page: %v", page)
		task := workerpool.NewTask(func(data interface{}) error {
			taskID := data.(*url.URL)
			time.Sleep(100 * time.Millisecond)
			log.Printf("Task %v processed\n", taskID.String())

			entry := questionanswer.NewEntry(taskID)

			client := httpclient.New(nil)
			err := entry.FetchData(client)
			if err != nil {
				return err
			}

			err = SavingEntry(entry, &manticoreStorages)

			return err

		}, page, &manticoreStorages)

		allTask = append(allTask, task)
	}

	pool := workerpool.NewPool(allTask, cfg.Workers)
	pool.Run()

	wg.Wait()
	log.Println("finished, all tasks successfully processed.")
}

func SavingEntry(entry *questionanswer.Entry, manticoreStorages *[]answer.Entries) error {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	answerEntries := *MakeAnswerEntries(entry)

	for _, storage := range *manticoreStorages {

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
func MakeAnswerEntries(entry *questionanswer.Entry) *[]answer.Entry {
	var entries []answer.Entry
	position := 1
	answerEntry := answer.Entry{
		Username: entry.Title,
		Text:     fmt.Sprintf("<h4>%v</h4> <p><span class=\"link\">%v</span></p>", entry.Title, entry.Video.String()),
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

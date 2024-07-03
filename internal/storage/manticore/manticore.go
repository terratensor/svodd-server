package manticore

import (
	openapiclient "github.com/manticoresoftware/manticoresearch-go"

	"github.com/terratensor/svodd-server/internal/entities/answer"
)

var _ answer.StorageInterface = &Client{}

type Response struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Hits     struct {
		Total         int    `json:"total"`
		TotalRelation string `json:"total_relation"`
		Hits          []struct {
			Id     string `json:"_id"`
			Score  int    `json:"_score"`
			Source struct {
				Title      string `json:"title"`
				Summary    string `json:"summary"`
				Content    string `json:"content"`
				ResourceID int    `json:"resource_id"`
				Chunk      int    `json:"chunk"`
				Published  int64  `json:"published"`
				Updated    int64  `json:"updated"`
				Created    int64  `json:"created"`
				UpdatedAt  int64  `json:"updated_at"`
				Language   string `json:"language"`
				Url        string `json:"url"`
				Author     string `json:"author"`
				Number     string `json:"number"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type Client struct {
	apiClient *openapiclient.APIClient
	Index     string
}
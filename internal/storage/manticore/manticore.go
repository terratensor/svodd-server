package manticore

import (
	"context"
	"fmt"
	"log"

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

func New(tbl string) (*Client, error) {
	// Initialize ApiClient
	configuration := openapiclient.NewConfiguration()
	configuration.Servers = openapiclient.ServerConfigurations{
		{
			// URL: "http://manticore:9308", // Здесь должна быть переменная окружения manticore host:port
			URL:         "http://localhost:9308",
			Description: "Default Manticore Search HTTP",
		},
	}
	//configuration.ServerURL(1, map[string]string{"URL": "http://manticore:9308"})
	apiClient := openapiclient.NewAPIClient(configuration)

	query := fmt.Sprintf(`show tables like '%v'`, tbl)

	// Проверяем существует ли таблица tbl, если нет, то создаем
	resp, _, err := apiClient.UtilsAPI.Sql(context.Background()).Body(query).Execute()
	if err != nil {
		return nil, err
	}
	data := resp[0]["data"].([]interface{})

	if len(data) > 0 {
		myMap := data[0].(map[string]interface{})
		indexValue := myMap["Index"]

		if indexValue != tbl {
			err := createTable(apiClient, tbl)
			if err != nil {
				return nil, err
			}
		}
	} else {
		err := createTable(apiClient, tbl)
		if err != nil {
			return nil, err
		}
	}

	return &Client{apiClient: apiClient, Index: tbl}, nil
}

func createTable(apiClient *openapiclient.APIClient, tbl string) error {

	log.Println("creating table", tbl)
	query := fmt.Sprintf("create table %v(username text, `text` text, avatar_file text, url string, role string, datetime timestamp, data_id int, parent_id int, type int, position int, chunk int) min_infix_len='3' index_exact_words='1' morphology='stem_en, stem_ru' index_sp='1'", tbl)

	sqlRequest := apiClient.UtilsAPI.Sql(context.Background()).Body(query)
	_, _, err := apiClient.UtilsAPI.SqlExecute(sqlRequest)
	if err != nil {
		return err
	}

	return nil
}


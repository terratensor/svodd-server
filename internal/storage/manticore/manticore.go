package manticore

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

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
			Id     int64 `json:"_id"`
			Score  int   `json:"_score"`
			Source struct {
				Username   string `json:"username"`
				Text       string `json:"text"`
				AvatarFile string `json:"avatar_file"`
				Url        string `json:"url"`
				Role       string `json:"role"`
				Datetime   int64  `json:"datetime"`
				DataID     int    `json:"data_id"`
				ParentID   int    `json:"parent_id"`
				Type       int    `json:"type"`
				Position   int    `json:"position"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type DBEntry struct {
	Username   string `json:"username"`
	Text       string `json:"text"`
	Url        string `json:"url"`
	AvatarFile string `json:"avatar_file"`
	Role       string `json:"role"`
	Datetime   int64  `json:"datetime"`
	DataID     int64  `json:"data_id"`
	ParentID   int64  `json:"parent_id"`
	Type       int    `json:"type"`
	Position   int    `json:"position"`
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

func castTime(value *time.Time) int64 {
	if value == nil || time.Time.IsZero(*value) {
		return 0
	}
	return value.Unix()
}

func createTable(apiClient *openapiclient.APIClient, tbl string) error {

	log.Println("creating table", tbl)
	query := fmt.Sprintf("create table %v(username text, `text` text, avatar_file text, url string, role string, datetime timestamp, data_id int, parent_id int, type int, position int) min_infix_len='3' index_exact_words='1' morphology='stem_en, stem_ru' index_sp='1'", tbl)

	sqlRequest := apiClient.UtilsAPI.Sql(context.Background()).Body(query)
	_, _, err := apiClient.UtilsAPI.SqlExecute(sqlRequest)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Insert(ctx context.Context, entry *answer.Entry) (*int64, error) {

	dbe := &DBEntry{
		Username:   entry.Username,
		Text:       entry.Text,
		Url:        entry.Url,
		AvatarFile: entry.AvatarFile,
		Role:       entry.Role,
		Datetime:   castTime(entry.Datetime),
		DataID:     entry.DataID,
		ParentID:   entry.ParentID,
		Type:       entry.Type,
		Position:   entry.Position,
	}

	//marshal into JSON buffer
	buffer, err := json.Marshal(dbe)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	var doc map[string]interface{}
	err = json.Unmarshal(buffer, &doc)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling buffer: %v", err)
	}

	idr := openapiclient.InsertDocumentRequest{
		Index: c.Index,
		Doc:   doc,
	}

	resp, r, err := c.apiClient.IndexAPI.Insert(ctx).InsertDocumentRequest(idr).Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v", r)
		return nil, fmt.Errorf("error when calling `IndexAPI.Insert``: %v", err)
	}

	return resp.Id, nil
}

func (c *Client) FindAllByUrl(ctx context.Context, url string) (*[]answer.Entry, error) {
	// response from `Search`: SearchRequest
	searchRequest := *openapiclient.NewSearchRequest(c.Index)

	filter := map[string]interface{}{"url": url}
	query := map[string]interface{}{"equals": filter}
	limit := 10000
	sort := []map[string]interface{}{{"position": "asc"}}

	searchRequest.SetQuery(query)
	searchRequest.SetLimit(int32(limit))
	searchRequest.SetSort(sort)

	resp, r, err := c.apiClient.SearchAPI.Search(ctx).SearchRequest(searchRequest).Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		return nil, fmt.Errorf("error when calling `SearchAPI.Search.Equals``: %v", err)
	}

	res := &Response{}
	respBody, _ := io.ReadAll(r.Body)
	err = json.Unmarshal(respBody, res)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", resp)
	}

	hits := res.Hits.Hits

	var entries []answer.Entry
	for _, hit := range hits {

		id := hit.Id
		sr := hit.Source
		jsonData, err := json.Marshal(sr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %v", resp)
		}

		var dbe DBEntry
		err = json.Unmarshal(jsonData, &dbe)
		if err != nil {
			log.Fatal(err)
		}

		datetime := time.Unix(dbe.Datetime, 0)

		ent := &answer.Entry{
			ID:         &id,
			Username:   dbe.Username,
			Text:       dbe.Text,
			Url:        dbe.Url,
			AvatarFile: dbe.AvatarFile,
			Role:       dbe.Role,
			Datetime:   &datetime,
			DataID:     dbe.DataID,
			ParentID:   dbe.ParentID,
			Type:       dbe.Type,
			Position:   dbe.Position,
		}

		entries = append(entries, *ent)
	}

	return &entries, nil
}

func getEntryID(response *openapiclient.SearchResponse) (*int64, error) {
	if len(response.Hits.Hits) == 0 {
		return nil, nil
	}

	hit := response.Hits.Hits[0]

	id, err := strconv.ParseInt(hit["_id"].(string), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ID to int64: %w", err)
	}

	return &id, nil
}

func makeDBEntry(resp *openapiclient.SearchResponse) *DBEntry {
	var hits []map[string]interface{}
	hits = resp.Hits.Hits

	// Если слайс Hits пустой (0) значит нет совпадений
	if len(hits) == 0 {
		return nil
	}

	hit := hits[0]

	sr := hit["_source"]
	jsonData, err := json.Marshal(sr)

	var dbe DBEntry
	err = json.Unmarshal(jsonData, &dbe)
	if err != nil {
		log.Fatal(err)
	}

	return &dbe
}

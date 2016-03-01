package groveclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func Int(value int) *int {
	return &value
}

func String(value string) *string {
	return &value
}

func Bool(value bool) *bool {
	return &value
}

type Config struct {
	Host string
}

type Client struct {
	Config
	httpClient *http.Client
}

type Post struct {
	Uid              string           `json:"uid"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	Document         *json.RawMessage `json:"document,omitempty"`
	ExternalDocument *json.RawMessage `json:"external_document,omitempty"`
	Sensitive        *json.RawMessage `json:"sensitive,omitempty"`
	Protected        *json.RawMessage `json:"protected,omitempty"`
	Tags             []string         `json:"tags"`
	Deleted          bool             `json:"deleted"`
	Published        bool             `json:"deleted"`
	ExternalId       string           `json:"external_id"`
}

type DocumentAttributes map[string]interface{}

func (post *Post) GetDocument() (DocumentAttributes, error) {
	return post.extractDocumentAttributes(post.Document)
}

func (post *Post) GetExternalDocument() (DocumentAttributes, error) {
	return post.extractDocumentAttributes(post.ExternalDocument)
}

func (post *Post) GetSensitive() (DocumentAttributes, error) {
	return post.extractDocumentAttributes(post.Sensitive)
}

func (post *Post) GetProtected() (DocumentAttributes, error) {
	return post.extractDocumentAttributes(post.Protected)
}

func (post *Post) extractDocumentAttributes(raw *json.RawMessage) (DocumentAttributes, error) {
	var attrs DocumentAttributes
	if raw != nil {
		if err := json.Unmarshal(*raw, &attrs); err != nil {
			return nil, err
		}
	}
	return attrs, nil
}

type GetOptions struct {
	Limit *int
	Raw   *bool
}

type GetResultItem struct {
	Post *Post `json:"post"`
}

type GetOutput struct {
	Posts []GetResultItem `json:"posts"`
}

func NewClient(config Config) *Client {
	return &Client{config, &http.Client{
		Timeout: 60 * time.Second,
	}}
}

func (client *Client) GetMany(uids []string, options GetOptions) (*GetOutput, error) {
	queryString := make(url.Values)
	if options.Raw != nil && *options.Raw == false {
		queryString["raw"] = []string{"false"}
	}
	if options.Limit != nil {
		queryString["limit"] = []string{strconv.Itoa(*options.Limit)}
	}

	uidList := strings.Join(uids, ",")
	if len(uids) > 1 {
		uidList = uidList + ","
	}

	url := fmt.Sprintf("http://%s/api/grove/v1/posts/%s?%s",
		client.Host, uidList, queryString.Encode())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.httpClient.Do(req)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Request [GET %s] failed with status code %d: %s",
			url, resp.StatusCode, resp.Status)
	}

	var output GetOutput
	reader := json.NewDecoder(resp.Body)
	if err := reader.Decode(&output); err != nil {
		return nil, err
	}

	return &output, nil
}

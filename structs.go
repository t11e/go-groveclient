package groveclient

import (
	"encoding/json"
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

type DocumentAttributes map[string]interface{}

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

type PostItem struct {
	Post *Post `json:"post"`
}

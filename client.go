package groveclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	pc "github.com/t11e/go-pebbleclient"
)

//go:generate go run vendor/github.com/vektra/mockery/cmd/mockery/mockery.go -name=Client -case=underscore

type Client interface {
	Get(uid string, options GetOptions) (*PostItem, error)
	GetMany(uids []string, options GetManyOptions) (*GetManyOutput, error)
	Update(uid string, pu PostUpdate, options UpdateOptions) (*PostItem, error)
}

type PostUpdate struct {
	CreatedAt        *time.Time      `json:"created_at,omitempty"`
	UpdatedAt        *time.Time      `json:"updated_at,omitempty"`
	Document         json.RawMessage `json:"document,omitempty"`
	ExternalDocument json.RawMessage `json:"external_document,omitempty"`
	Sensitive        json.RawMessage `json:"sensitive,omitempty"`
	Protected        json.RawMessage `json:"protected,omitempty"`
	Tags             *[]string       `json:"tags,omitempty",`
	Deleted          *bool           `json:"deleted,omitempty"`
	Published        *bool           `json:"published,omitempty"`
	ExternalId       string          `json:"external_id,omitempty"`
	Version          int             `json:"version,omitempty"`
}

type client struct {
	c pc.Client
}

type GetOptions struct {
	Raw *bool
}

type GetManyOptions struct {
	Limit *int
	Raw   *bool
}

type UpdateOptions struct {
	Merge      *bool   `json:"merge,omitempty"`
	ExternalID *string `json:"external_id,omitempty"`
}

type GetManyOutput struct {
	Posts []PostItem `json:"posts"`
}

// Register registers us in a connector.
func Register(connector *pc.Connector) {
	connector.Register((*Client)(nil), func(client pc.Client) (pc.Service, error) {
		return New(client)
	})
}

func New(pebbleClient pc.Client) (Client, error) {
	return &client{pebbleClient.WithOptions(pc.Options{
		ServiceName: "grove",
		APIVersion:  1,
	})}, nil
}

func (c *client) Get(uid string, options GetOptions) (*PostItem, error) {
	params := pc.Params{
		"raw": options.Raw != nil && *options.Raw,
		"uid": uid,
	}

	var out PostItem
	err := c.c.Get("/posts/:uid", &pc.RequestOptions{
		Params: params,
	}, &out)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *client) GetMany(uids []string, options GetManyOptions) (*GetManyOutput, error) {
	uidList := strings.Join(uids, ",")
	if len(uids) > 1 {
		uidList = uidList + ","
	}

	params := pc.Params{
		"raw":  options.Raw != nil && *options.Raw,
		"uids": uidList,
	}
	if options.Limit != nil {
		params["limit"] = *options.Limit
	}

	var out GetManyOutput
	err := c.c.Get("/posts/:uids", &pc.RequestOptions{
		Params: params,
	}, &out)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *client) Update(uid string, pu PostUpdate, options UpdateOptions) (*PostItem, error) {
	params := pc.Params{
		"uid": uid,
	}
	if options.Merge != nil {
		params["merge"] = *options.Merge
	}
	if options.ExternalID != nil {
		params["external_id"] = *options.ExternalID
	}
	payload, err := json.Marshal(&struct {
		Post PostUpdate `json:"post"`
	}{pu})
	if err != nil {
		return nil, err
	}
	result := PostItem{}
	err = c.c.Put("/posts/:uid", &pc.RequestOptions{
		Params: params,
	}, bytes.NewReader(payload), &result)
	if err == nil {
		return &result, nil
	}
	if reqErr, ok := err.(*pc.RequestError); ok {
		switch reqErr.Resp.StatusCode {
		case 404:
			return &result, NoSuchPostError{uid}
		case 409:
			return &result, ConflictError{uid}
		}
	}
	return &result, err
}

type NoSuchPostError struct {
	UID string
}

func (e NoSuchPostError) Error() string {
	return fmt.Sprintf("grove post not found: %s", e.UID)
}

type ConflictError struct {
	UID string
}

func (e ConflictError) Error() string {
	return fmt.Sprintf("grove post failed to update due to version conflict: %s", e.UID)
}

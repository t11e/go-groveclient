package groveclient

import (
	"bytes"
	"encoding/json"
	"strings"

	pc "github.com/t11e/go-pebbleclient"
)

//go:generate go run vendor/github.com/vektra/mockery/cmd/mockery/mockery.go -name=Client -case=underscore

type Client interface {
	Get(uid string, options GetOptions) (*PostItem, error)
	GetMany(uids []string, options GetManyOptions) (*GetManyOutput, error)
	Update(postItem *PostItem) (*PostItem, error)
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
		ApiVersion:  1,
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

func (c *client) Update(postItem *PostItem) (*PostItem, error) {
	params := pc.Params{
		"uid": postItem.Post.Uid,
	}
	payload, err := json.Marshal(postItem)
	if err != nil {
		return nil, err
	}
	var out PostItem
	err = c.c.Put("/posts/:uid", &pc.RequestOptions{
		Params: params,
	}, bytes.NewReader(payload), &out)
	if err != nil {
		return nil, err
	}
	return &out, err
}

package groveclient

import (
	"strings"

	pc "github.com/t11e/go-pebbleclient"
)

type Client struct {
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

func New(client pc.Client) (*Client, error) {
	return &Client{client.Options(pc.Options{
		ServiceName: "grove",
		ApiVersion:  1,
	})}, nil
}

func (client *Client) Get(uid string, options GetOptions) (*PostItem, error) {
	params := pc.Params{
		"raw": options.Raw != nil && *options.Raw,
		"uid": uid,
	}

	var out PostItem
	err := client.c.Get("/posts/:uid", &pc.RequestOptions{
		Params: params,
	}, &out)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (client *Client) GetMany(uids []string, options GetManyOptions) (*GetManyOutput, error) {
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
	err := client.c.Get("/posts/:uids", &pc.RequestOptions{
		Params: params,
	}, &out)
	if err != nil {
		return nil, err
	}
	return &out, err
}

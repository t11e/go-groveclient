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
	return &Client{client}, nil
}

func (client *Client) Get(uid string, options GetOptions) (*PostItem, error) {
	params := pc.Params{
		"raw": options.Raw != nil && *options.Raw,
	}

	var out PostItem
	err := client.c.Get(pc.FormatPath("/posts/:uid", pc.Params{"uid": uid}), &pc.RequestOptions{
		Params: params,
	}, &out)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (client *Client) GetMany(uids []string, options GetManyOptions) (*GetManyOutput, error) {
	params := pc.Params{
		"raw": options.Raw != nil && *options.Raw,
	}
	if options.Limit != nil {
		params["limit"] = *options.Limit
	}

	uidList := strings.Join(uids, ",")
	if len(uids) > 1 {
		uidList = uidList + ","
	}

	var out GetManyOutput
	err := client.c.Get(pc.FormatPath("/posts/:uids", pc.Params{"uids": uidList}), &pc.RequestOptions{
		Params: params,
	}, &out)
	if err != nil {
		return nil, err
	}
	return &out, err
}

package groveclient

import (
	"net/http"
	"strconv"
	"strings"

	pc "github.com/t11e/go-pebbleclient"
)

type Client struct {
	client *pc.Client
}

func New(client *pc.Client) (*Client, error) {
	return &Client{client}, nil
}

// NewFromRequest constructs a new client from an HTTP request.
func NewFromRequest(options pc.ClientOptions, req *http.Request) (*Client, error) {
	if options.AppName == "" {
		options.AppName = "grove"
	}
	client, err := pc.NewFromRequest(options, req)
	if err != nil {
		return nil, err
	}
	return &Client{client}, nil
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

func (client *Client) GetMany(uids []string, options GetOptions) (*GetOutput, error) {
	params := pc.Params{
		"raw": options.Raw != nil && *options.Raw,
	}
	if options.Limit != nil {
		params["limit"] = []string{strconv.Itoa(*options.Limit)}
	}

	uidList := strings.Join(uids, ",")
	if len(uids) > 1 {
		uidList = uidList + ","
	}

	var out GetOutput
	err := client.client.Get(uidList, &params, &out)
	if err != nil {
		return nil, err
	}
	return &out, err
}

package groveclient

import (
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (c *MockClient) Get(uid string, options GetOptions) (*PostItem, error) {
	args := c.Called(uid, options)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PostItem), nil
}

func (c *MockClient) GetMany(uids []string, options GetManyOptions) (*GetManyOutput, error) {
	args := c.Called(uids, options)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetManyOutput), nil
}

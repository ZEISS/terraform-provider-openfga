package sdk

import (
	"context"

	openfga "github.com/openfga/go-sdk/client"
	"github.com/zeiss/pkg/cast"
)

// CreateStore ...
func (c *Client) CreateStore(name string) (*Store, error) {
	resp, err := c.fga.CreateStore(context.Background()).Body(openfga.ClientCreateStoreRequest{Name: name}).Execute()
	if err != nil {
		return nil, err
	}

	store := Store{
		ID:   resp.GetId(),
		Name: resp.GetName(),
	}

	return cast.Ptr(store), nil
}

// GetStore ...
func (c *Client) GetStore(id string) (*Store, error) {
	resp, err := c.fga.GetStore(context.Background()).Options(openfga.ClientGetStoreOptions{StoreId: cast.Ptr(id)}).Execute()
	if err != nil {
		return nil, err
	}

	store := Store{
		ID:   resp.GetId(),
		Name: resp.GetName(),
	}

	return cast.Ptr(store), nil
}

// DeleteStore ...
func (c *Client) DeleteStore(id string) error {
	_, err := c.fga.DeleteStore(context.Background()).Options(openfga.ClientDeleteStoreOptions{StoreId: cast.Ptr(id)}).Execute()
	if err != nil {
		return err
	}

	return nil
}

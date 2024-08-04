package sdk

import (
	"context"

	openfga "github.com/openfga/go-sdk/client"
	"github.com/zeiss/pkg/cast"
)

// CreateStore ...
func (c *Client) CreateStore(store Store) (*Store, error) {
	resp, err := c.fga.CreateStore(context.Background()).Body(openfga.ClientCreateStoreRequest{Name: store.Name}).Execute()
	if err != nil {
		return nil, err
	}

	store = Store{
		ID:   resp.GetId(),
		Name: resp.GetName(),
	}

	return cast.Ptr(store), nil
}

// DeleteStore ...
func (c *Client) DeleteStore(store Store) error {
	_, err := c.fga.DeleteStore(context.Background()).Options(openfga.ClientDeleteStoreOptions{StoreId: cast.Ptr(store.ID)}).Execute()
	if err != nil {
		return err
	}

	return nil
}

package sdk

import (
	"context"

	openfga "github.com/openfga/go-sdk/client"
	"github.com/zeiss/pkg/cast"
)

// TupleModel ...
type TupleModel struct {
	ID string `json:"id,omitempty"`
}

// CreateTuple ...
func (c *Client) CreateTuple(ctx context.Context, store, model, user, relation, object string) (*TupleModel, error) {
	body := openfga.ClientWriteRequest{
		Writes: []openfga.ClientTupleKey{
			{
				User:     user,
				Relation: relation,
				Object:   object,
			},
		},
	}

	_, err := c.fga.Write(ctx).Options(openfga.ClientWriteOptions{StoreId: cast.Ptr(store), AuthorizationModelId: cast.Ptr(model)}).Body(body).Execute()
	if err != nil {
		return nil, err
	}

	tuple := TupleModel{}

	return cast.Ptr(tuple), nil
}

// GetTuple ...
func (c *Client) GetTuple(ctx context.Context, store, model, user, relation, object string) (*TupleModel, error) {
	_, err := c.fga.Read(ctx).Options(openfga.ClientReadOptions{StoreId: cast.Ptr(store)}).Body(openfga.ClientReadRequest{}).Execute()
	if err != nil {
		return nil, err
	}

	tuple := TupleModel{}

	return cast.Ptr(tuple), nil
}

// DeleteTuple ...
func (c *Client) DeleteTuple(ctx context.Context, store, model, user, relation, object string) error {
	body := openfga.ClientWriteRequest{
		Deletes: []openfga.ClientTupleKeyWithoutCondition{
			{
				User:     user,
				Relation: relation,
				Object:   object,
			},
		},
	}

	_, err := c.fga.Write(ctx).Options(openfga.ClientWriteOptions{StoreId: cast.Ptr(store), AuthorizationModelId: cast.Ptr(model)}).Body(body).Execute()
	if err != nil {
		return err
	}

	return nil
}

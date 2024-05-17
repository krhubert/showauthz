package client

import (
	"context"

	"rift/authz/testauthz"
)

func StartTestServer(ctx context.Context) (*Client, error) {
	authclient, err := testauthz.StartMemServer(ctx)
	if err != nil {
		return nil, err
	}

	client := &Client{c: authclient}
	if err := client.MigrateSchema(ctx); err != nil {
		return nil, err
	}
	return client, nil
}

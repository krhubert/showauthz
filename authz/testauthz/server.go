package testauthz

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/spicedb/pkg/cmd/datastore"
	"github.com/authzed/spicedb/pkg/cmd/server"
	"github.com/authzed/spicedb/pkg/cmd/util"
)

func StartMemServer(ctx context.Context) (*authzed.ClientWithExperimental, error) {
	ds, err := datastore.NewDatastore(
		ctx,
		datastore.DefaultDatastoreConfig().ToOption(),
	)
	if err != nil {
		return nil, err
	}
	srv, err := server.NewConfigWithOptions(
		server.WithDatastore(ds),
		server.WithDispatchMaxDepth(50),
		server.WithMaximumPreconditionCount(1000),
		server.WithMaximumUpdatesPerWrite(1000),
		server.WithStreamingAPITimeout(30*time.Second),
		server.WithMaxCaveatContextSize(4096),
		server.WithMaxRelationshipContextSize(25000),
		server.WithGRPCServer(util.GRPCServerConfig{
			Network: util.BufferedNetwork,
			Enabled: true,
		}),
		server.WithGRPCAuthFunc(func(ctx context.Context) (context.Context, error) {
			return ctx, nil
		}),
		server.WithHTTPGateway(util.HTTPServerConfig{HTTPEnabled: false}),
		server.WithMetricsAPI(util.HTTPServerConfig{HTTPEnabled: false}),
		server.WithDispatchServer(util.GRPCServerConfig{Enabled: false}),
	).Complete(ctx)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := srv.Run(ctx); err != nil {
			panic(fmt.Sprintf("error running server: %v", err))
		}
	}()

	conn, err := srv.GRPCDialContext(ctx)
	if err != nil {
		return nil, err
	}

	return &authzed.ClientWithExperimental{
		Client: authzed.Client{
			SchemaServiceClient:      v1.NewSchemaServiceClient(conn),
			PermissionsServiceClient: v1.NewPermissionsServiceClient(conn),
			WatchServiceClient:       v1.NewWatchServiceClient(conn),
		},
		ExperimentalServiceClient: v1.NewExperimentalServiceClient(conn),
	}, nil
}
